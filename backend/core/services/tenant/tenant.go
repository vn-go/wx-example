package tenant

import (
	"core/models"
	"core/services/config"
	"sync"

	"github.com/vn-go/dx"
)

type TenantService struct {
	db      *dx.DB
	isMulti bool
}
type initGetSecret struct {
	val  string
	err  error
	once sync.Once
}

var initGetSecretMap sync.Map

func (tenantSvc *TenantService) GetSecret(tenant string) (string, error) {
	a, _ := initGetSecretMap.LoadOrStore(tenant, &initGetSecret{})
	i := a.(*initGetSecret)
	i.once.Do(func() {
		db, err := tenantSvc.GetDb(tenant)
		if err != nil {
			i.err = err
		}
		appItem := &models.SysApp{}
		err = db.First(appItem)
		if err != nil {
			i.err = err
		}
		i.val = appItem.SecretKey
	})
	return i.val, i.err
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
