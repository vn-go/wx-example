package acc

import (
	"core/services/jwt"
	"core/services/pwd"
	"core/services/tenant"

	"github.com/vn-go/dx"
)

type AccService struct {
	db        *dx.DB
	jwtSvc    *jwt.JwtService
	pwdSvc    pwd.PwdService
	tenantSvc *tenant.TenantService
}

func NewAccService(
	db *dx.DB,
	jwtSvc *jwt.JwtService,
	pwdSvc pwd.PwdService,
	tenantSvc *tenant.TenantService) *AccService {
	return &AccService{
		db:        db,
		jwtSvc:    jwtSvc,
		pwdSvc:    pwdSvc,
		tenantSvc: tenantSvc,
	}
}
