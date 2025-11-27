package acc

import (
	"context"
	"core/models"
	"time"

	"github.com/vn-go/dx"
)

func (acc *AccService) LoginAndGetJWT(ctx context.Context, tenant, username, password string) (*OAuthResponse, error) {
	user := &models.SysUsers{}
	err := acc.db.First(user, "username=?", username)
	if err != nil {
		if dx.Errors.IsRecordNotFound(err) {
			return nil, acc.errSvc.Unauthenticate()
		}
		return nil, err
	}
	ok, err := acc.pwdSvc.ComparePassword(ctx, tenant, username, password, user.HashPassword)
	if err != nil {
		return nil, err
	}
	if ok {
		secret, err := acc.tenantSvc.GetSecret(tenant)
		if err != nil {
			return nil, err
		}
		tk, err := acc.jwtSvc.NewJWTWithSecret(secret, tenant, user, 2*time.Hour)
		if err != nil {
			return nil, err
		}
		refresToken, err := acc.CreateRefreshToken(ctx, tenant, username, tk)
		if err != nil {
			return nil, err
		}
		return &OAuthResponse{
			AccessToken:  tk,
			TokenType:    "Bearer",
			ExpiresIn:    int((2 * time.Hour).Seconds()),
			RefreshToken: refresToken,
		}, nil
		// return acc.jwtSvc.NewJWTWithSecret(secret, "app", user, 2*time.Hour)
	} else {
		return nil, acc.errSvc.Unauthenticate()
	}
}
