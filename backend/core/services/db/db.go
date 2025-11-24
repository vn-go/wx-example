package db

import (
	"core/services/config"
	"sync"

	"github.com/vn-go/dx"
)

type DbService struct {
	db     *dx.DB
	cfgSvc *config.ConfigService
}

var dbServiceGetOnce sync.Once

func (dbSvc *DbService) Get() *dx.DB {
	var err error
	dbServiceGetOnce.Do(func() {
		cgf := dbSvc.cfgSvc.Get()
		dbSvc.db, err = dx.Open(cgf.Database.Driver, cgf.Database.DSN)
	})
	if err != nil {
		panic(err)
	}
	return dbSvc.db
}
func (dbSvc *DbService) Close() {
	if dbSvc.db != nil {
		dbSvc.db.Close()
	}
}
func NewDbService(cfgSvc *config.ConfigService) *DbService {
	return &DbService{
		cfgSvc: cfgSvc,
	}
}
