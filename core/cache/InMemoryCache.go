package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type InMemoryCache struct {
	mu    sync.RWMutex
	store map[string][]byte
	ttl   map[string]time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[string][]byte),
		ttl:   make(map[string]time.Time),
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string, target any) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	exp, ok := c.ttl[key]
	if ok && time.Now().After(exp) {
		return errors.New("cache expired")
	}

	val, ok := c.store[key]
	if !ok {
		return errors.New("not found")
	}
	return json.Unmarshal(val, target)
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.store[key] = data
	if ttl > 0 {
		c.ttl[key] = time.Now().Add(ttl)
	}
	return nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	delete(c.ttl, key)
	return nil
}

func (c *InMemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	exp, ok := c.ttl[key]
	if ok && time.Now().After(exp) {
		return false, nil
	}
	_, exists := c.store[key]
	return exists, nil
}

func (c *InMemoryCache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = nil
	c.ttl = nil
	return nil
}
