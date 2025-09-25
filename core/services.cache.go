package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/vn-go/bx"
)

type cacheService interface {
	// Set(service any, ctx context.Context, key string, data any) error
	// Get(service any, ctx context.Context, key string, target any) error
	// Delete(service any, ctx context.Context, key string) error
	AddObject(ctx context.Context, tenant, key string, target any, timeInHour int) error
	GetObject(ctx context.Context, tenant, key string, target any) error
	DeleteObject(ctx context.Context, tenant, key string, target any) error
}
type cacheServiceImpl struct {
	cacher    bx.Cache
	prefixKey string
}
type initCacheServiceImplSha265 struct {
	val  string
	once sync.Once
}

var cacheInitCacheServiceImplSha265 sync.Map

func (c *cacheServiceImpl) sha265(content string) string {
	a, _ := cacheInitCacheServiceImplSha265.LoadOrStore(content, &initCacheServiceImplSha265{})
	i := a.(*initCacheServiceImplSha265)
	i.once.Do(func() {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s://%s", c.prefixKey, content)))
		i.val = hex.EncodeToString(hash[:])
	})
	return i.val
}

//	func (c *cacheServiceImpl) Set(service any, ctx context.Context, key string, data any) error {
//		return c.cacher.Set(ctx, c.sha265(reflect.TypeOf(service).String()+"/"+key), data, time.Hour*4)
//	}
//
//	func (c *cacheServiceImpl) Get(service any, ctx context.Context, key string, target any) error {
//		return c.cacher.Get(ctx, c.sha265(reflect.TypeOf(service).String()+"/"+key), target)
//	}
//
//	func (c *cacheServiceImpl) Delete(service any, ctx context.Context, key string) error {
//		return c.cacher.Delete(ctx, c.sha265(reflect.TypeOf(service).String()+"/"+key))
//	}
func (c *cacheServiceImpl) AddObject(ctx context.Context, tenant, key string, target any, timeInHour int) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	keyStr := c.sha265(fmt.Sprintf("%s,%s://%s@%s/%s", c.prefixKey, tenant, key, typ.PkgPath(), typ.String()))
	return c.cacher.Set(ctx, keyStr, target, time.Hour*time.Duration(timeInHour))
}
func (c *cacheServiceImpl) GetObject(ctx context.Context, tenant, key string, target any) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	keyStr := c.sha265(fmt.Sprintf("%s,%s://%s@%s/%s", c.prefixKey, tenant, key, typ.PkgPath(), typ.String()))
	return c.cacher.Get(ctx, keyStr, target)
}
func (c *cacheServiceImpl) DeleteObject(ctx context.Context, tenant, key string, target any) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	keyStr := c.sha265(fmt.Sprintf("%s,%s://%s@%s/%s", c.prefixKey, tenant, key, typ.PkgPath(), typ.String()))
	return c.cacher.Delete(ctx, keyStr)
}

var cac bx.Cache
var newCacherOne sync.Once

func newCacher() (bx.Cache, error) {
	var err error
	newCacherOne.Do(func() {
		cac = bx.Cacher.NewInMemoryCache()
	})
	return cac, err

}
func newCacheServiceImpl(cfg *configInfo) cacheService {
	if cfg.Cache.Type == "app" {
		return &cacheServiceImpl{
			cacher: bx.Cacher.NewInMemoryCache(),
		}
	} else {
		panic(fmt.Sprintf("Not implement ccahe type %s", cfg.Cache.Type))
	}
}
func newCacheService(cfg *configInfo) (cacheService, error) {
	var cache bx.Cache
	var err error
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

	return &cacheServiceImpl{
		cacher: cache,
	}, nil
}
