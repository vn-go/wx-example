package tenant

import (
	"core/models"
	"core/services/config"

	"github.com/vn-go/dx"
)

type TenantService struct {
	db      *dx.DB
	isMulti bool
}

func (tenantSvc *TenantService) GetSecret(tenant string) (string, error) {
	db, err := tenantSvc.GetDb(tenant)
	if err != nil {
		return "", err
	}
	appItem := &models.SysApp{}
	err = db.First(appItem)
	if err != nil {
		return "", err
	}
	return appItem.SecretKey, nil
}

func (tenantSvc *TenantService) GetDb(tenant string) (*dx.DB, error) {
	if tenantSvc.isMulti {
		return tenantSvc.db.NewDB(tenant)
	} else {
		return tenantSvc.db, nil
	}
}

func NewTenantService(
	db *dx.DB,
	cfgSvc *config.ConfigService,
) *TenantService {

	return &TenantService{
		db:      db,
		isMulti: cfgSvc.Get().Tenant.IsMulti,
	}
}
