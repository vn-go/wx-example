package app

import (
	"core/services/config"
	"core/services/jwt"
	"core/services/pwd"

	"github.com/vn-go/dx"
)

/*
This service will start app setup everything
Such as:

	New some table in Database with default data
*/
type AppService struct {
	cfgSvc *config.ConfigService
	jwtSvc *jwt.JwtService
	pwdSvc pwd.PwdService

	db *dx.DB
}

func NewAppService(
	cfgSvc *config.ConfigService,
	jwtSvc *jwt.JwtService,
	db *dx.DB,
	pwdSvc pwd.PwdService,
) *AppService {
	return &AppService{
		cfgSvc: cfgSvc,
		jwtSvc: jwtSvc,
		db:     db,
		pwdSvc: pwdSvc,
	}
}
