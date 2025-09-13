package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Driver string `mapstructure:"driver"`
		DSN    string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	Cache struct {
		Type string `mapstructure:"type"`
	} `mapstructure:"cache"`
}
type initLoadConfig struct {
	val  *Config
	err  error
	once sync.Once
}

var cacheLoadConfig sync.Map

func LoadConfig(path string) (*Config, error) {
	actually, _ := cacheLoadConfig.LoadOrStore("LoadConfig", &initLoadConfig{})
	init := actually.(*initLoadConfig)
	init.once.Do(func() {
		init.val, init.err = loadConfig(path)
	})
	return init.val, init.err
}
func loadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Enable ENV support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
