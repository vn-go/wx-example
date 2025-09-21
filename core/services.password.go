package core

import (
	"context"
	"fmt"
)

type passwordService interface {
	HashPassword(username, password string) (string, error)
	ComparePassword(ctx context.Context, tenant string, username, password, hashPass string) (bool, error)
	DeleteCacheByUsername(ctx context.Context, tenant, username string) error
}
type passwordServiceType struct {
	Bcrypt passwordService
	Argon2 passwordService
}
type passwordServiceImpl struct {
}

func newpasswordService(cfg *configInfo, cache cacheService) (passwordService, error) {
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

// var PasswordService = &passwordServiceType{
// 	Bcrypt: &bcryptPasswordService{},
// 	Argon2: &argon2PasswordService{},
// }
