package acc

import (
	"context"
	"core/models"
	"time"
)

func (acc *AccService) LoginAndGetJWT(ctx context.Context, tenant, username, password string) (string, error) {
	user := &models.SysUsers{}
	err := acc.db.First(user, "username=?", username)
	if err != nil {
		return "", err
	}
	ok, err := acc.pwdSvc.ComparePassword(ctx, tenant, username, password, user.HashPassword)
	if err != nil {
		return "", err
	}
	if ok {
		secret, err := acc.tenantSvc.GetSecret(tenant)
		if err != nil {
			return "", err
		}
		return acc.jwtSvc.NewJWTWithSecret(secret, "app", user, 2*time.Hour)
	} else {
		return "", nil
	}
}
