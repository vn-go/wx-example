package core

import (
	"context"
	"core/models"
	"fmt"

	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

type tenantService interface {
	GetTenant(tenant string) (*dx.DB, error)
	CloseAllTenants() error
	CreateTenant(ctx context.Context, dbName, name string) error
	GetSecretKey(ctx context.Context, dbName string) (string, error)
}
type tenantServiceSql struct {
	db          *dx.DB
	mapDbtenant map[string]*dx.DB
	jwtSvc      jwtService
}

func (s *tenantServiceSql) GetSecretKey(ctx context.Context, dbName string) (string, error) {
	return bx.OnceCall[tenantServiceSql, string]("GetSecretKey/"+dbName, func() (string, error) {
		tenantData := &models.Tenant{}
		err := s.db.First(tenantData, "dbName=?", dbName)
		if err != nil {
			return "", err
		}
		return tenantData.ShareSecret, nil
	})
}
func (s *tenantServiceSql) CreateTenant(ctx context.Context, dbName, name string) error {
	shareSecret, err := s.jwtSvc.GenerateSecret()
	if err != nil {
		return err
	}
	tenant, err := dx.NewDTO[models.Tenant]()
	if err != nil {
		return err
	}

	tenant.Name = name
	tenant.DbName = dbName
	tenant.ShareSecret = shareSecret
	err = s.db.InsertWithContext(ctx, tenant)
	if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
		if dbErr.ErrorType == dx.Errors.DUPLICATE {

			r := s.db.Model(&models.Tenant{}).Where("dbName=?", tenant.DbName).Update(map[string]interface{}{
				"shareSecret": tenant.ShareSecret,
			})
			if r.Error != nil {
				return r.Error
			}
		}
	}
	return nil

}
func (s *tenantServiceSql) GetTenant(tenant string) (*dx.DB, error) {
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
func newTenantService(db *dx.DB, jwtSvc jwtService) tenantService {

	return &tenantServiceSql{
		db:     db,
		jwtSvc: jwtSvc,
	}
}
