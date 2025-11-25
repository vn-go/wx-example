package acc

import (
	"context"
	"core/models"
	"core/services/aes"
	"core/services/base"
	"core/services/jwt"
	"core/services/pwd"
	"core/services/tenant"

	"github.com/vn-go/dx"
)

type AccService struct {
	base.BaseService
	db        *dx.DB
	jwtSvc    *jwt.JwtService
	pwdSvc    pwd.PwdService
	tenantSvc *tenant.TenantService

	aes.AesService
}

func (acc *AccService) ListUsers(context context.Context, data *jwt.Indentifier) (any, error) {
	db, err := acc.tenantSvc.GetDb(data.Tenant)
	if err != nil {
		return nil, err
	}
	return db.DslToArray("SysUsers()")
}

func (acc *AccService) ValidateToken(context context.Context, authorization string) (*jwt.Indentifier, error) {
	token, err := acc.jwtSvc.ExtractToken(authorization)
	if err != nil {
		return nil, err
	}
	if token.TokenType != "Bearer" {
		return nil, acc.ErrService.ForbiddenErr()
	}
	tokenClaims, err := acc.jwtSvc.DecodeJwtTokenWithouSignature(token.Token)
	if err != nil {
		return nil, acc.ErrService.Unauthenticate()
	}
	tenant := tokenClaims.Data.Tenant
	secret, err := acc.tenantSvc.GetSecret(tenant)
	if err != nil {
		return nil, err
	}
	ok, err := acc.jwtSvc.VerifyToken(token.Token, secret)
	if err != nil || !ok {
		return nil, acc.ErrService.Unauthenticate()
	}

	return &tokenClaims.Data, nil
}

func (acc *AccService) GetAllUsers(context context.Context) ([]models.SysUsers, error) {
	users := []models.SysUsers{}
	if err := acc.db.DslQuery(&users, "SysUsers()"); err != nil {
		return nil, err
	}
	return users, nil
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
