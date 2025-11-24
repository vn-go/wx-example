package cacher

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/vn-go/bx"
)

type CacheService interface {
	// Set(service any, ctx context.Context, key string, data any) error
	// Get(service any, ctx context.Context, key string, target any) error
	// Delete(service any, ctx context.Context, key string) error
	AddObject(ctx context.Context, tenant, key string, target any, timeInHour int) error
	GetObject(ctx context.Context, tenant, key string, target any) error
	DeleteObject(ctx context.Context, tenant, key string, target any) error
}

type CacheServiceImpl struct {
	cacher    bx.Cache
	prefixKey string
}

func (c *CacheServiceImpl) AddObject(ctx context.Context, tenant, key string, target any, timeInHour int) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	keyStr := c.sha265(fmt.Sprintf("%s,%s://%s@%s/%s", c.prefixKey, tenant, key, typ.PkgPath(), typ.String()))
	return c.cacher.Set(ctx, keyStr, target, time.Hour*time.Duration(timeInHour))
}
func (c *CacheServiceImpl) GetObject(ctx context.Context, tenant, key string, target any) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	keyStr := c.sha265(fmt.Sprintf("%s,%s://%s@%s/%s", c.prefixKey, tenant, key, typ.PkgPath(), typ.String()))
	return c.cacher.Get(ctx, keyStr, target)
}
func (c *CacheServiceImpl) DeleteObject(ctx context.Context, tenant, key string, target any) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	keyStr := c.sha265(fmt.Sprintf("%s,%s://%s@%s/%s", c.prefixKey, tenant, key, typ.PkgPath(), typ.String()))
	return c.cacher.Delete(ctx, keyStr)
}
