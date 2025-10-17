// ----------------------------------------------------------------------------------
// Tenant Service
// ----------------------------------------------------------------------------------
// This file implements the tenant management service.
// Responsibilities include:
//   - Creating and managing tenants
//   - Maintaining tenant-specific databases
//   - Generating and caching JWT secret keys
//   - Creating default admin accounts
//
// The service uses the `dx` package for database access and the `bx` package
// for caching function results. Synchronization primitives like `sync.Once`
// are used to ensure idempotent operations.
// ----------------------------------------------------------------------------------

package core

import (
	"context"
	"core/models"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

/*
TenantService defines methods for managing tenants.
*/
type tenantService interface {
	/*
		Retrieve a tenant database from the master database.
		If the tenant does not exist, return a dbError with ErrorType = NOT_FOUND.
	*/
	GetTenant(tenant string) (*dx.DB, error)

	/*
		Each call to GetTenant creates a new *dx.DB (once per tenant).
		All tenant databases are kept alive during the lifetime of the app.
		Call CloseAllTenants() when the app is shutting down.
	*/
	CloseAllTenants() error

	/*
		Create a tenant. This method will insert a new tenant with:
		  - Database name
		  - Tenant name
		  - Shared secret key (used to generate JWT tokens)
	*/
	CreateTenant(ctx context.Context, dbName, name string) error

	/*
		Get the shared secret key for a tenant.
		The secret key is cached to avoid repeated database lookups.
	*/
	GetSecretKey(ctx context.Context, dbName string) (string, error)

	/*
		Create the default admin tenant account.
		This special account is allowed to create, update, and delete tenants.
	*/
	CreateDefaultAdminTenantAccount(ctx context.Context) error
	Login(ctx context.Context, username, password string) (*OAuthResponse, error)
	GetAppInfo(ctx context.Context) (app *models.App, err error)
	CreateDefaultRootUser(tenant string, ctx context.Context) error
}

type tenantServiceSql struct {
	db          *dx.DB
	mapDbtenant map[string]*dx.DB
	jwtSvc      jwtService
	PwdSvc      passwordService
	usrRepo     userRepo
	cfg         *configInfo
}

/*
Retrieve a tenant's shared secret key.
Uses bx.OnceCall to cache the result and avoid repeated database queries.
*/
func (s *tenantServiceSql) GetSecretKey(ctx context.Context, dbName string) (string, error) {
	return bx.OnceCall[tenantServiceSql, string]("GetSecretKey/"+dbName, func() (string, error) {
		if s.cfg.Tenant.IsMulti {
			// if multi tenant get share secret from tanent db
			tenantData := &models.Tenant{}
			err := s.db.First(tenantData, "dbName=?", dbName)
			if err != nil {
				return "", err
			}
			return tenantData.ShareSecret, nil
		} else {
			// if not multi tenant get share secret from current db
			appData := &models.App{}
			err := s.db.First(appData, "name=?", s.cfg.AppName)
			if err != nil {
				return "", err
			}
			return appData.ShareSecret, nil
		}

	})
}

var (
	// Ensures ShareSecret is generated if missing for an existing tenant.
	updateTenantShareSecretIfEmptyOnce sync.Once
)

/*
CreateTenant inserts a new tenant into the master database.
If the tenant already exists, it ensures a ShareSecret is present.
It also creates a default root user for the tenant database.
*/
func (s *tenantServiceSql) CreateTenant(ctx context.Context, dbName, name string) error {
	if dbName == "" {
		return errors.New("dbName can not be  empty")
	}
	shareSecret, err := s.jwtSvc.GenerateSecret()
	if err != nil {
		return err
	}
	tenant, err := dx.NewDTO[models.Tenant]()
	if err != nil {
		return err
	}

	// TODO: Add validation for tenant name and dbName
	tenant.Name = name
	tenant.DbName = dbName
	tenant.ShareSecret = shareSecret

	err = s.db.InsertWithContext(ctx, tenant)
	if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
		if dbErr.ErrorType == dx.Errors.DUPLICATE {
			var makeShareSecret error
			updateTenantShareSecretIfEmptyOnce.Do(func() {
				// Ensure ShareSecret is set if missing
				r := s.db.Model(&models.Tenant{}).
					Where("dbName=? and (shareSecret is null or ShareSecret=?)", tenant.DbName, "").
					Update(map[string]interface{}{"shareSecret": tenant.ShareSecret})
				if r.Error != nil {
					makeShareSecret = r.Error
				}
			})
			if makeShareSecret != nil {
				return makeShareSecret
			}
		}
	}

	dbTenant, err := s.GetTenant(dbName)
	if err != nil {
		return err
	}
	hashPass, err := s.PwdSvc.HashPassword(s.cfg.DefaultAuth.Username, s.cfg.DefaultAuth.Password)
	if err != nil {
		return err
	}
	err = s.usrRepo.CreateDefaultUser(dbTenant, ctx, s.cfg.DefaultAuth.Username, hashPass)
	if err != nil {
		return err
	}
	return nil
}

/*
GetTenant opens (once) a new tenant database by its name.
The database is cached in mapDbtenant for reuse.
*/
func (s *tenantServiceSql) GetTenant(tenant string) (*dx.DB, error) {
	if !s.cfg.Tenant.IsMulti {
		if err := s.CreateDefaultAdminTenantAccount(context.Background()); err != nil {
			return nil, err
		}
		return s.db, nil
	}
	return bx.OnceCall[tenantServiceSql, *dx.DB](tenant, func() (*dx.DB, error) {
		tenantData := &models.Tenant{}
		err := s.db.First(tenantData, "dbName=?", tenant)
		if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
			if dbErr.ErrorType != dx.Errors.NOTFOUND {
				return nil, dbErr
			}
		}
		if tenant == s.db.DbName {
			return s.db, nil
		}
		db, err := s.db.NewDB(tenant)
		if err != nil {
			return nil, err
		}
		if s.mapDbtenant == nil {
			s.mapDbtenant = make(map[string]*dx.DB)
		}
		s.mapDbtenant[tenant] = db
		return db, nil
	})
}

/*
CloseAllTenants closes all tenant database connections.
Call this when the application is shutting down.
*/
func (s *tenantServiceSql) CloseAllTenants() error {
	for k, v := range s.mapDbtenant {
		fmt.Printf("close %s\n", k)
		err := v.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	// Ensures default admin tenant account is created only once.
	createDefaultAdminTenantAccountOnce sync.Once
)

/*
CreateDefaultAdminTenantAccount creates a default admin user in the master database.
This operation is performed only once.
*/
func (s *tenantServiceSql) CreateDefaultAdminTenantAccount(ctx context.Context) error {
	var err error
	createDefaultAdminTenantAccountOnce.Do(func() {
		app, err := dx.NewDTO[models.App]()
		if err != nil {
			return
		}
		app.Name = s.cfg.AppName
		app.ShareSecret, err = s.jwtSvc.GenerateSecret()
		if err != nil {
			return
		}
		err = s.db.Insert(app)
		if err != nil {
			if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
				if dbErr.ErrorType == dx.Errors.DUPLICATE {
					// If already exists, ignore error
					err = nil
				} else {
					return
				}
			}
		}
		adminTenantUser, err := dx.NewDTO[models.User]()
		if err != nil {
			return
		}
		adminTenantUser.Username = s.cfg.Tenant.DefaultAdminUser
		adminTenantUser.IsTenantAdmin = dx.Ptr(true)
		hashPasswords, err := s.PwdSvc.HashPassword(
			s.cfg.Tenant.DefaultAdminUser,
			s.cfg.Tenant.DefaultAdminPassword,
		)
		if err != nil {
			return
		}
		adminTenantUser.HashPassword = hashPasswords
		err = s.db.Insert(adminTenantUser)
		if err != nil {
			if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
				if dbErr.ErrorType == dx.Errors.DUPLICATE {
					// If already exists, ignore error
					err = nil
					return
				}
			}
		}
	})
	return err
}
func (s *tenantServiceSql) GetAppInfo(ctx context.Context) (app *models.App, err error) {
	if err := s.db.WithContext(ctx).First(&app, "name=?", s.cfg.AppName); err != nil {
		return nil, err
	}
	return app, nil
}
func (s *tenantServiceSql) Login(ctx context.Context, username, password string) (*OAuthResponse, error) {
	err := s.CreateDefaultAdminTenantAccount(ctx)
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	if err := s.db.First(user, "username=? and isActive=? and IsTenantAdmin=?", username, true, true); err != nil {
		if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
			if dbErr.ErrorType == dx.Errors.NOTFOUND {
				return nil, nil
			}
		}
	}
	ok, err := s.PwdSvc.ComparePassword(ctx, "", username, password, user.HashPassword)
	if err != nil {
		return nil, err
	}
	//get sharesecret form db
	app := &models.App{}
	if err := s.db.First(app, "name=?", s.cfg.AppName); err != nil {
		return nil, err
	}

	if ok {
		token, err := s.jwtSvc.NewJWTWithSecret(app.ShareSecret, user.UserId, "admin", "", "", time.Hour*2)
		if err != nil {
			return nil, err
		}
		return &OAuthResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   int((time.Hour * 2).Seconds()),
		}, nil
	} else {
		return nil, nil
	}
}

type initCreateDefaultRootUser struct {
	err  error
	once sync.Once
}

var initCreateDefaultRootUserCache sync.Map

func (s *tenantServiceSql) CreateDefaultRootUser(tenant string, ctx context.Context) error {
	a, _ := initCreateDefaultRootUserCache.LoadOrStore(tenant, &initCreateDefaultRootUser{})
	i := a.(*initCreateDefaultRootUser)
	i.once.Do(func() {
		hashPwd, err := s.PwdSvc.HashPassword(s.cfg.DefaultAuth.Username, s.cfg.DefaultAuth.Password)
		if err != nil {
			i.err = err
			return
		}
		if s.cfg.Tenant.IsMulti {
			db, err := s.GetTenant(tenant)
			if err != nil {
				i.err = err
				return
			}
			err = s.usrRepo.CreateDefaultUser(db, ctx, s.cfg.DefaultAuth.Username, hashPwd)
			if err != nil {
				if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
					if dbErr.ErrorType == dx.Errors.DUPLICATE {
						return
					}
				}
				i.err = err
				return
			}
		} else {
			err = s.usrRepo.CreateDefaultUser(s.db, ctx, s.cfg.DefaultAuth.Username, hashPwd)
			if err != nil {
				if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
					if dbErr.ErrorType == dx.Errors.DUPLICATE {
						return
					}
				}
				i.err = err
				return
			}
		}
		return
	})
	return i.err

}

/*
newTenantService is a constructor for tenantServiceSql.
*/
func newTenantService(
	db *dx.DB,
	jwtSvc jwtService,
	PwdSvc passwordService,
	usrRepo userRepo,
	cfg *configInfo,
) (tenantService, error) {
	ret := &tenantServiceSql{
		db:      db,
		jwtSvc:  jwtSvc,
		PwdSvc:  PwdSvc,
		usrRepo: usrRepo,
		cfg:     cfg,
	}
	if cfg.Tenant.IsMulti {
		if err := ret.CreateDefaultAdminTenantAccount(context.Background()); err != nil {
			return nil, err
		}
	} else {
		return ret, ret.CreateDefaultRootUser("", context.Background())
	}

	return ret, nil
}
