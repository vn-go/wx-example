package core

import (
	"context"
	"core/models"
)

// type DataSource[TResult any] struct {
// }
type dataService interface {
	GetList(ctx context.Context, user *UserClaims, dsName string, filter string, selector string) (any, error)
	GetSQL(ctx context.Context, user *UserClaims, dsName string, filter string, selector string) (string, error)
}
type dataServiceImpl struct {
	tenantSvc tenantService
}

func (ds *dataServiceImpl) GetSQL(ctx context.Context, user *UserClaims, dsName string, filter string, selector string) (string, error) {
	db, err := ds.tenantSvc.GetTenant(user.Tenant)
	if err != nil {
		return "", err
	}
	source := db.ModelDatasource(dsName)

	if selector != "" {
		source.Select(selector)
	}
	if filter != "" {
		source.Where(filter)
	}
	data, err := source.ToSql()
	return data.Sql, err
}

func (ds *dataServiceImpl) GetList(ctx context.Context, user *UserClaims, dsName, filter, selector string) (any, error) {

	db, err := ds.tenantSvc.GetTenant(user.Tenant)

	if err != nil {
		return nil, err
	}
	err = db.Insert(&models.TrackFilter{
		DsName:   dsName,
		Filter:   filter,
		Selector: selector,
	})
	if err != nil {
		return nil, err
	}
	source := db.ModelDatasource(dsName)
	source.ToSql()
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
