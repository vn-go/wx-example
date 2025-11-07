package core

import "context"

type viewService struct {
	tenantDb tenantService
}

func (v *viewService) GetListOfViews(context context.Context, data *UserClaims) (any, error) {
	db, err := v.tenantDb.GetTenant(data.Tenant)
	if err != nil {
		return nil, err
	}
	return db.DslToArray("uiView()")
}

func newViewService(tenantDb tenantService) *viewService {
	return &viewService{tenantDb: tenantDb}
}
