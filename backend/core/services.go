package core

import (
	"core/services/acc"
	"core/services/aes"
	"core/services/config"
	"core/services/data"
	"core/services/errs"
	"core/services/pwd"
	"core/services/tenant"

	"core/services/app"
	"core/services/cacher"
	"core/services/jwt"

	"github.com/vn-go/dx"
)

type Service struct {
	ConfigSvc *config.ConfigService
	CacheSvc  cacher.CacheService
	AppSvc    *app.AppService
	//DbSvc     *db.DbService
	JwtSvc    *jwt.JwtService
	PwdSvc    pwd.PwdService
	AccSvc    *acc.AccService
	Db        *dx.DB
	TenantSvc *tenant.TenantService
	ErrSvc    *errs.ErrService
	AesSvc    *aes.AesService
	DataSvc   *data.DataSignService
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
	errs.NewErrService,
	aes.NewAesService,
	data.NewDataSignService,
}

type User jwt.Indentifier
