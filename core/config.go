package core

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/vn-go/bx"
)

type configInfo struct {
	Database struct {
		Driver string `mapstructure:"driver"`
		DSN    string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	Cache struct {
		Type string `mapstructure:"type"`
	} `mapstructure:"cache"`
}
type configType struct {
}

func loadConfig(path string) (*configInfo, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Enable ENV support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg configInfo
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
func (c *configType) Load(path string) (*configInfo, error) {
	return bx.OnceCall[configType]("Load/"+path, func() (*configInfo, error) {
		return loadConfig(path)
	})
}

var Config = &configType{}
