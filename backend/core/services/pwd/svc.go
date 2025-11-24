package pwd

import (
	"context"
	"core/services/cacher"
	"core/services/config"
	"fmt"
)

type PwdService interface {
	HashPassword(username, password string) (string, error)
	ComparePassword(ctx context.Context, tenant string, username, password, hashPass string) (bool, error)
}

func NewPwdService2(
	cfgSvc *config.ConfigService,
	//cache cacher.CacheService
) (PwdService, error) {
	panic("OK")
}
func NewPwdService(
	cfgSvc *config.ConfigService,
	cache cacher.CacheService) (PwdService, error) {
	cfg := cfgSvc.Get()
	if cfg.Jwt.HashPasswordType == "bcrypt" {
		return &bcryptPasswordService{
			cache: cache,
		}, nil
	}
	if cfg.Jwt.HashPasswordType == "argon2" {
		return &argon2PasswordService{
			cache: cache,
		}, nil
	}
	return nil, fmt.Errorf("%s was not support", cfg.Jwt.HashPasswordType)

}
