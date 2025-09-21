package core

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/vn-go/bx"
)

type cacheService interface {
	Set(service any, ctx context.Context, key string, data any) error
	Get(service any, ctx context.Context, key string, target any) error
	Delete(service any, ctx context.Context, key string) error
}
type cacheServiceImpl struct {
	cacher bx.Cache
}

func (c *cacheServiceImpl) Set(service any, ctx context.Context, key string, data any) error {
	return c.cacher.Set(ctx, reflect.TypeOf(service).String()+"/"+key, data, time.Hour*4)
}
func (c *cacheServiceImpl) Get(service any, ctx context.Context, key string, target any) error {
	return c.cacher.Get(ctx, reflect.TypeOf(service).String()+"/"+key, target)
}
func (c *cacheServiceImpl) Delete(service any, ctx context.Context, key string) error {
	return c.cacher.Delete(ctx, reflect.TypeOf(service).String()+"/"+key)
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
func newBadgerCache() (cacheService, error) {
	cache, err := bx.Cacher.NewBadgerCache("./Badger")
	if err != nil {
		return nil, err
	}
	return &cacheServiceImpl{
		cacher: cache,
	}, nil
}
