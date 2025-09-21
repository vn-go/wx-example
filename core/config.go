package core

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

type configInfo struct {
	Database struct {
		Driver string `mapstructure:"driver"`
		DSN    string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	Cache struct {
		Type     string `mapstructure:"type"`
		Memcache struct {
			Servers string `mapstructure:"server`
		} `mapstructure:"memcache"`
		Redis struct {
			Address  string `mapstructure:address`
			Password string `mapstructure.password`
			Db       int    `mapstructure.db`
		} `mapstructure:"redis"`
		Badger struct {
			Directory string `mapstructure:directory`
		} `mapstructure:"bager"`
	} `mapstructure:"cache"`
	Jwt struct {
		Sercret          string `mapstructure:"default-sercret-key"`
		HashPasswordType string `mapstructure:"hash-type"`
		SecretLen        int    `mapstructure:"share-secret-len"`
	} `mapstructure:"JWT"`
	Log struct {
		Path     string `mapstructure:"path"`
		Size     int    `mapstructure:"size"`
		Age      int    `mapstructure:"age"`
		Backup   int    `mapstructure:"backup"`
		Compress bool   `mapstructure:"compress"`
	} `mapstructure:"log"`
	Cryptor struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"cryptor"`
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

func newConfig() (*configInfo, error) {
	log.Printf("loading config %s", "./config.yaml")
	t := time.Now().UTC()
	defer func() {
		n := time.Since(t).Milliseconds()
		log.Printf("loading config %s,in %d", "./config.yaml", n)
	}()
	return loadConfig("./config.yaml")
}

var (
	newDBOnce    sync.Once
	dbConnection *dx.DB
)

func newDB2(cfg *configInfo) (*dx.DB, error) {
	var err error
	newDBOnce.Do(func() {
		dbConnection, err = dx.Open(cfg.Database.Driver, cfg.Database.DSN)
	})
	return dbConnection, err
}
func newDB(cfg *configInfo) (*dx.DB, error) {
	fmt.Println("new db")
	db, err := dx.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)
	return db, nil
}
