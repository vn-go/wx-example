package core

import (
	"core/services/config"
	"sync"

	"github.com/vn-go/dx"
)

var globalDb *dx.DB
var newDbOnce sync.Once

func newDb(cfgSvc *config.ConfigService) *dx.DB {
	db, err := dx.Open(cfgSvc.Get().Database.Driver, cfgSvc.Get().Database.DSN)
	if err != nil {
		panic(err)
	}
	return db
}
