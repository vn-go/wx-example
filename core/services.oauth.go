package core

import (
	"context"
	"core/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

type OAuthResponse struct {
	AccessToken  string `json:"access_token"`            // Bắt buộc
	TokenType    string `json:"token_type"`              // Bắt buộc, thường là "Bearer"
	ExpiresIn    int    `json:"expires_in"`              // Bắt buộc (giây)
	RefreshToken string `json:"refresh_token,omitempty"` // Tùy chọn
	IDToken      string `json:"id_token,omitempty"`      // Tùy chọn (nếu dùng OIDC)
	Scope        string `json:"scope,omitempty"`         // Tùy chọn
}

type serviceAuth interface {
	Login(tenant string, ctx context.Context, username, password string) (*OAuthResponse, error)
	Verify(ctx context.Context, authorization string) (*models.User, error)
}
type serviceOAuth struct {
	user        userRepo
	cache       cacheService
	passwordSvc passwordService
	tenant      tenantService
	jwtSvc      jwtService
	db          *dx.DB
}

func (s *serviceOAuth) generateAccessToken(ctx context.Context, tenant string, user *models.User) (string, error) {
	key, err := s.tenant.GetSecretKey(ctx, tenant)
	if err != nil {
		return "", err
	}
	roldeCode := bx.IsNull(user.RoleCode, "")
	return s.jwtSvc.NewJWTWithSecret(key, user.UserId, tenant, roldeCode, bx.IsNull(user.Email, ""), time.Hour*4)
}
func (s *serviceOAuth) generateRefreshToken(ctx context.Context, tenant string, user *models.User) (string, error) {
	db, err := s.tenant.GetTenant(tenant)
	if err != nil {
		return "", err
	}
	ret := uuid.NewString()
	err = db.WithContext(ctx).Insert(&models.RefreshToken{
		Token:  ret,
		UserId: user.Id,
	})
	if err != nil {
		if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
			if dbErr.ErrorType == dx.Errors.DUPLICATE {
				err := db.WithContext(ctx).Model(&models.RefreshToken{}).Where("userId=?", user.Id).Update(map[string]interface{}{
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

func (s *serviceOAuth) Login(tenant string, ctx context.Context, username, password string) (*OAuthResponse, error) {
	//key := fmt.Sprintf("%s:%s@%s(Login)", username, password, tenant)
	ret := &OAuthResponse{}
	cacheItem := &OAuthResponseCacheItem{}
	if err := s.cache.GetObject(ctx, tenant, strings.ToLower(username), cacheItem); err == nil {
		return &cacheItem.Oauth, nil
	}
	db, err := s.tenant.GetTenant(tenant) // get tenant db by tenant name
	if err != nil {
		return nil, err
	}
	user, err := s.user.GetUserByName(db, ctx, username)
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
		s.cache.AddObject(ctx, tenant, strings.ToLower(username), OAuthResponseCacheItem{
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
func (s *serviceOAuth) Verify(ctx context.Context, authorization string) (user *models.User, err error) {
	palyload, err := s.jwtSvc.DecodeJWTNoVerify(authorization)
	if err != nil {
		return nil, err
	}
	if palyload.Issuer == "admin" {
		app, err := s.tenant.GetAppInfo(ctx)
		if err != nil {
			return nil, err
		}
		_, err = s.jwtSvc.VerifyJWTWithSecret(app.ShareSecret, authorization)
		if err != nil {
			return nil, err
		}
		return s.user.GetUserByUserId(s.db, ctx, palyload.Subject)
	} else {
		shareSecret, err := s.tenant.GetSecretKey(ctx, palyload.Issuer)
		if err != nil {
			return nil, err
		}
		_, err = s.jwtSvc.VerifyJWTWithSecret(shareSecret, authorization)
		if err == nil {
			return nil, err
		}
		return s.user.GetUserByUserId(s.db, ctx, palyload.Subject)
	}

}
func newServiceAuth(
	user userRepo,
	cache cacheService,
	passwordSvc passwordService,
	tenant tenantService,
	jwtSvc jwtService,
	db *dx.DB,

) serviceAuth {
	return &serviceOAuth{
		user:        user,
		cache:       cache,
		passwordSvc: passwordSvc,
		tenant:      tenant,
		jwtSvc:      jwtSvc,
		db:          db,
	}
}
