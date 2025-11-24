package cacher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"core/services/config"

	"github.com/vn-go/bx"
)

func NewCacheService(cfgSvc *config.ConfigService) (CacheService, error) {
	var cache bx.Cache
	var err error
	cfg := cfgSvc.Get()
	if err != nil {
		return nil, err
	}

	if cfg.Cache.Type == "badger" {
		os.MkdirAll(cfg.Cache.Badger.Directory, 0755)
		//Cannot create lock file
		i := 0
		ok := false
		for !ok {
			bagerDir := filepath.Join(cfg.Cache.Badger.Directory, fmt.Sprintf("app[%d]", i))
			cache, err = bx.Cacher.NewBadgerCache(bagerDir)
			if err != nil {
				if strings.Contains(err.Error(), "Cannot create lock file") {
					i++
					continue
				} else {
					log.Fatalln(err)
					return nil, err
				}
			} else {
				ok = true
			}
		}

	} else if cfg.Cache.Type == "redis" {
		cache = bx.Cacher.NewRedisCache(
			cfg.Cache.Redis.Address,
			cfg.Cache.Redis.Password,
			cfg.Cache.Redis.Db,
		)
	} else if cfg.Cache.Type == "memcache" {
		cache = bx.Cacher.NewMemCache(strings.Split(cfg.Cache.Memcache.Servers, ",")...)
	} else {
		cache = bx.Cacher.NewInMemoryCache()
	}

	return &CacheServiceImpl{
		cacher: cache,
	}, nil
}
