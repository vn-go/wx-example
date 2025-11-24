package core

import (
	"core/services/acc"
	"core/services/config"
	"core/services/pwd"
	"core/services/tenant"

	"core/services/app"
	"core/services/cacher"
	"core/services/jwt"

	"github.com/vn-go/dx"
)

type services struct {
	ConfigSvc *config.ConfigService
	CacheSvc  cacher.CacheService
	AppSvc    *app.AppService
	//DbSvc     *db.DbService
	JwtSvc    *jwt.JwtService
	PwdSvc    pwd.PwdService
	AccSvc    *acc.AccService
	Db        *dx.DB
	TenantSvc *tenant.TenantService
}

var servicesInjectors = []any{
	newConfigService,
	cacher.NewCacheService,
	app.NewAppService,
	newDb,
	jwt.NewJwtService,
	pwd.NewPwdService,
	//db.NewDbService,
	acc.NewAccService,
	tenant.NewTenantService,
}
