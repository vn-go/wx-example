package core

import (
	"context"
)

// type DataSource[TResult any] struct {
// }
type dataService interface {
	GetList(ctx context.Context, user *UserClaims, dsName string, filter string, selector string) (any, error)
}
type dataServiceImpl struct {
	tenantSvc tenantService
}

func (ds *dataServiceImpl) GetList(ctx context.Context, user *UserClaims, dsName, filter, selector string) (any, error) {
	db, err := ds.tenantSvc.GetTenant(user.Tenant)
	if err != nil {
		return nil, err
	}
	source := db.ModelDatasource(dsName)

	if selector != "" {
		source.Select(selector)
	}
	if filter != "" {
		source.Where(filter)
	}
	data, err := source.ToDict()
	return data, err
}
func newDataService(tenantSvc tenantService) dataService {
	return &dataServiceImpl{
		tenantSvc: tenantSvc,
	}
}
