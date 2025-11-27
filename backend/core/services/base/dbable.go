package base

import (
	"core/services/data"
)

type dbableService struct {
}
type BaseService struct {
	dataSvc data.DataSignService
}
