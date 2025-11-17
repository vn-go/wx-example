package core

import (
	"context"
	"core/models"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

type serviceOAuthUser struct {
	HashPassword string
	UserId       string
	Username     string
	Email        *string
	RoleCode     *string
}
type OAuthResponse struct {
	AccessToken  string `json:"access_token"`            // Bắt buộc
	TokenType    string `json:"token_type"`              // Bắt buộc, thường là "Bearer"
	ExpiresIn    int    `json:"expires_in"`              // Bắt buộc (giây)
	RefreshToken string `json:"refresh_token,omitempty"` // Tùy chọn
	IDToken      string `json:"id_token,omitempty"`      // Tùy chọn (nếu dùng OIDC)
	Scope        string `json:"scope,omitempty"`         // Tùy chọn
}

type serviceAuth interface {
	LoadUserCache(tenant string) error
	Login(tenant string, ctx context.Context, username, password string) (*OAuthResponse, error)
	//verify authorization header
	//return
	//user,tenant,error
	Verify(ctx context.Context, authorization string) (*models.User, string, error)
	DeleteUserCache(ctx context.Context, tenant, username string) error
	//check can user can use api in viewPath at tenant
	Authorize(context context.Context, tenant string, user *models.User, viewPath, apiPath string) (bool, error)
	ChangePassword(ctx context.Context, user *UserClaims, newPass string) error
}

type serviceOAuth struct {
	user        userRepo
	cache       cacheService
	passwordSvc passwordService
	tenant      tenantService
	jwtSvc      jwtService
	rabcSvc     rabcService
	secretSvc   *secretService
	db          *dx.DB
}

func (s *serviceOAuth) ChangePassword(ctx context.Context, user *UserClaims, newPass string) error {
	db, err := s.tenant.GetTenant(user.Tenant)
	if err != nil {
		return err
	}
	newHashPassword, err := s.passwordSvc.HashPassword(user.Username, newPass)
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Model(&models.User{}).Where("username=?", user.Username).Update(map[string]interface{}{
		"hashPassword": newHashPassword,
	}).Error
	if err != nil {
		return err
	}
	s.DeleteUserCache(ctx, user.Tenant, user.Username)
	return nil
}

func (s *serviceOAuth) Authorize(context context.Context, tenant string, user *models.User, viewPath string, apiPath string) (bool, error) {
	//check if user is admin
	if user.IsSysAdmin {
		err := s.rabcSvc.ResgisterView(context, tenant, viewPath, apiPath, user.Username)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	panic("unimplemented")
}
func (s *serviceOAuth) generateAccessToken(ctx context.Context, tenant string, user *serviceOAuthUser) (string, error) {
	key, err := s.tenant.GetSecretKey(ctx, tenant)
	if err != nil {
		return "", err
	}

	return s.jwtSvc.NewJWTWithSecret(key, user.UserId, tenant, bx.IsNull(user.RoleCode, ""), bx.IsNull(user.Email, ""), time.Hour*4)
}
func (s *serviceOAuth) generateRefreshToken(ctx context.Context, tenant string, user *serviceOAuthUser) (string, error) {
	db, err := s.tenant.GetTenant(tenant)
	if err != nil {
		return "", err
	}
	ret := uuid.NewString()
	err = db.WithContext(ctx).Insert(&models.RefreshToken{
		Token:  ret,
		UserId: user.UserId,
	})
	if err != nil {
		if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
			if dbErr.ErrorType == dx.Errors.DUPLICATE {
				err := db.WithContext(ctx).Model(&models.RefreshToken{}).Where("userId=?", user.UserId).Update(map[string]interface{}{
					"token": ret,
				}).Error
				if err != nil {
					return "", err
				}
				return ret, nil
			}
			return "", dbErr
		}
		return "", err
	}
	return ret, nil

}

type OAuthResponseCacheItem struct {
	Tanent   string
	Username string
	Oauth    OAuthResponse
}

func (s *serviceOAuth) getUserByUsername(ctx context.Context, tenant string, username string) (*serviceOAuthUser, error) {
	ret := &serviceOAuthUser{}
	err := s.cache.GetObject(ctx, tenant, strings.ToLower(username), ret)
	if err == nil {
		return ret, nil
	}
	db, err := s.tenant.GetTenant(tenant) // get tenant db by tenant name
	if err != nil {
		return nil, err
	}
	ret, err = dx.QueryItem[serviceOAuthUser](db, `user(HashPassword,username,UserId,email),
												role(code RoleCode),
												from(user,role,left(user.roleId=role.id)),
												where(username=?)`, username)
	if err != nil {
		return nil, err
	}
	s.cache.AddObject(ctx, tenant, strings.ToLower(username), ret, 4)
	return ret, nil
}
func (s *serviceOAuth) LoadUserCache(tenant string) error {
	db, err := s.tenant.GetTenant(tenant) // get tenant db by tenant name
	if err != nil {
		return err
	}
	users, err := dx.QueryItems[serviceOAuthUser](db, `user(HashPassword,username,UserId,email),
												role(code RoleCode),
												from(user,role,left(user.roleId=role.id)),sort(UserId)`)
	if err != nil {
		return err
	}
	for _, user := range users {

		s.cache.AddObject(context.Background(), tenant, strings.ToLower(user.Username), user, 4)
	}
	return nil

}

var loadUserCacheOne sync.Once

func (s *serviceOAuth) Login(tenant string, ctx context.Context, username, password string) (*OAuthResponse, error) {
	loadUserCacheOne.Do(func() {
		s.LoadUserCache(tenant)
	})
	//key := fmt.Sprintf("%s:%s@%s(Login)", username, password, tenant)
	ret := &OAuthResponse{}
	cacheItem := &OAuthResponseCacheItem{}
	if err := s.cache.GetObject(ctx, tenant, strings.ToLower(fmt.Sprintf("%s/%s", username, password)), cacheItem); err == nil {
		return &cacheItem.Oauth, nil
	}

	user, err := s.getUserByUsername(ctx, tenant, username)
	// user, err := s.user.GetUserByName(db, ctx, username)
	if err != nil {
		return nil, err
	}
	ok, err := s.passwordSvc.ComparePassword(ctx, tenant, username, password, user.HashPassword)
	if err != nil {
		return nil, err
	}
	if ok {

		accessToken, err := s.generateAccessToken(ctx, tenant, user)
		if err != nil {
			return nil, err
		}
		refreshToken, err := s.generateRefreshToken(ctx, tenant, user)
		if err != nil {
			return nil, err
		}
		ret = &OAuthResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    30,
			RefreshToken: refreshToken,
		}
		s.cache.AddObject(ctx, tenant, strings.ToLower(fmt.Sprintf("%s/%s", username, password)), OAuthResponseCacheItem{
			Tanent:   tenant,
			Username: strings.ToLower(username),
			Oauth:    *ret,
		}, 4)

		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		return nil, nil
	}

}

func (s *serviceOAuth) DeleteUserCache(ctx context.Context, tenant, username string) error {
	key, err := s.tenant.GetSecretKey(ctx, tenant)
	if err != nil {
		return err
	}
	var userVerify OAuthResponseCacheItem
	if err := s.cache.GetObject(ctx, tenant, username, &userVerify); err == nil {
		if err := s.cache.DeleteObject(ctx, tenant, strings.ToLower(username), &userVerify); err != nil {
			return err
		}
		err = s.jwtSvc.DeleteVerifyJWTWithSecretCache(key, userVerify.Oauth.AccessToken)
		if err != nil {
			return err
		}

	}
	err = s.cache.DeleteObject(ctx, tenant, username, &ComparePasswordCacheItem{})
	if err != nil {
		return err
	}
	err = s.cache.DeleteObject(ctx, tenant, strings.ToLower(username), &serviceOAuthUser{})
	if err != nil {
		return err
	}
	return nil

}

type cacheUserVerify struct {
	User               models.User
	Secret, AuthHeader string
}

func (s *serviceOAuth) Verify(ctx context.Context, authorization string) (user *models.User, tenant string, err error) {

	palyload, err := s.jwtSvc.DecodeJWTNoVerify(authorization)

	if err != nil {
		return nil, "", err
	}
	if palyload.Issuer == "admin" {
		app, err := s.tenant.GetAppInfo(ctx)
		if err != nil {
			return nil, "", err
		}
		_, err = s.jwtSvc.VerifyJWTWithSecret(app.ShareSecret, authorization)
		if err != nil {
			return nil, "", err

		}
		var userVerify cacheUserVerify
		if err := s.cache.GetObject(ctx, palyload.Issuer, palyload.Subject, &userVerify); err != nil {
			userFind, err := s.user.GetUserByUserId(s.db, ctx, palyload.Subject)
			if err != nil {
				return nil, "", err
			}
			userVerify = cacheUserVerify{
				User:       *userFind,
				Secret:     app.ShareSecret,
				AuthHeader: authorization,
			}

			hoursLeft := int(time.Until(palyload.ExpiresAt.Time).Hours())
			err = s.cache.AddObject(ctx, palyload.Issuer, palyload.Subject, userVerify, hoursLeft)
			if err != nil {
				return nil, "", err
			}

		}

		return &userVerify.User, palyload.Issuer, err
	} else {
		shareSecret, err := s.tenant.GetSecretKey(ctx, palyload.Issuer)
		if err != nil {
			return nil, "", err
		}
		_, err = s.jwtSvc.VerifyJWTWithSecret(shareSecret, authorization)
		if err != nil {
			return nil, "", err
		}
		var userVerify cacheUserVerify
		if err := s.cache.GetObject(ctx, palyload.Issuer, palyload.Subject, &userVerify); err != nil {
			userFind, err := s.user.GetUserByUserId(s.db, ctx, palyload.Subject)
			if err != nil {
				return nil, "", err
			}
			userVerify = cacheUserVerify{
				User:       *userFind,
				Secret:     shareSecret,
				AuthHeader: authorization,
			}

			hoursLeft := int(time.Until(palyload.ExpiresAt.Time).Hours())
			err = s.cache.AddObject(ctx, palyload.Issuer, palyload.Subject, userVerify, hoursLeft)
			if err != nil {
				return nil, "", err
			}

		}

		return &userVerify.User, palyload.Issuer, err

	}

}
func newServiceAuth(
	user userRepo,
	cache cacheService,
	passwordSvc passwordService,
	tenant tenantService,
	jwtSvc jwtService,
	db *dx.DB,
	rabc rabcService,
	secretSvc *secretService,

) serviceAuth {
	return &serviceOAuth{
		user:        user,
		cache:       cache,
		passwordSvc: passwordSvc,
		tenant:      tenant,
		jwtSvc:      jwtSvc,
		db:          db,
		rabcSvc:     rabc,
		secretSvc:   secretSvc,
	}
}
