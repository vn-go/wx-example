package base

import (
	"core/services/data"
	"core/services/errs"
)

type dbableService struct {
}
type BaseService struct {
	errs.ErrService
	dataSvc data.DataSignService
}
