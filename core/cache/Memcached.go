package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemCache struct {
	client *memcache.Client
}

func NewMemCache(servers ...string) *MemCache {
	return &MemCache{
		client: memcache.New(servers...),
	}
}

func (c *MemCache) Get(ctx context.Context, key string, target any) error {
	item, err := c.client.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(item.Value, target)
}

func (c *MemCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(ttl.Seconds()),
	})
}

func (c *MemCache) Delete(ctx context.Context, key string) error {
	return c.client.Delete(key)
}

func (c *MemCache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return false, nil
	}
	return err == nil, err
}

func (c *MemCache) Close() error {
	return nil // gomemcache không có Close
}
