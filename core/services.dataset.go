/*
defined datasset service
*/
package core

import (
	"context"
	"strings"
)

type datasetService interface {
	Execute(ctx context.Context, user *UserClaims, datasetName, fields, filters, sort string, limit, offset uint64) (interface{}, error)
}

type datasetServiceSql struct {
	tenant tenantService
}

func NewDatasetServiceSql(tenant tenantService) datasetService {
	return &datasetServiceSql{
		tenant: tenant,
	}
}
func (ds *datasetServiceSql) Execute(ctx context.Context, user *UserClaims, datasetName, fields, filters, sort string, limit, offset uint64) (interface{}, error) {
	db, err := ds.tenant.GetTenant(user.Tenant)
	if err != nil {
		return nil, err
	}

	dataset := db.DatasetWithContext(ctx)
	dataset.From(datasetName)
	if fields != "" {
		dataset.Select(fields)
	}
	if filters != "" {
		dataset.Where(filters)
	}
	if sort != "" {
		dataset.Sort(strings.Split(sort, ",")...)
	}
	if limit > 0 {
		dataset.Limit(limit)
	}
	if offset > 0 {
		dataset.Offset(offset)
	}
	return dataset.ToArray()
}
