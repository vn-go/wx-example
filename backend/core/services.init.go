package core

import (
	"core/services/config"
	"sync"

	"github.com/vn-go/bx"
)

var cfgService *config.ConfigService

func initConfig(configYamlPath string) *config.ConfigService {
	cfgService = config.NewConfigService()
	cfgService.SetConfigFilePath(configYamlPath)
	return cfgService
}
func newConfigService() *config.ConfigService {
	return cfgService
}

var Services *Service
var once sync.Once

func Start(configYamlPath string) {
	once.Do(func() {
		initConfig(configYamlPath)
		svcRet := (&bx.Wire[Service]{}).WireThenGet(
			servicesInjectors...,
		)
		Services = svcRet
	})

}
func (svc Service) Close() {
	svc.Db.Close()
}
