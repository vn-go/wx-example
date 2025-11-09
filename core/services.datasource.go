package core

import (
	"context"

	"github.com/vn-go/dx"
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
	qr := dx.NewDynamicQuery(selector)
	qr.Join(dsName)
	if filter != "" {
		qr.Filter(filter)
	}

	return qr.ToArrayWithContext(ctx, db)
}
func newDataService(tenantSvc tenantService) dataService {
	return &dataServiceImpl{
		tenantSvc: tenantSvc,
	}
}
