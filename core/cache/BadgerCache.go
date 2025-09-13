package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type BadgerCache struct {
	db *badger.DB
}

func NewBadgerCache(path string) (*BadgerCache, error) {
	opts := badger.DefaultOptions(path).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerCache{db: db}, nil
}

func (c *BadgerCache) Get(ctx context.Context, key string, target any) error {
	return c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return json.Unmarshal(val, target)
	})
}

func (c *BadgerCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), data)
		if ttl > 0 {
			e = e.WithTTL(ttl)
		}
		return txn.SetEntry(e)
	})
}

func (c *BadgerCache) Delete(ctx context.Context, key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (c *BadgerCache) Exists(ctx context.Context, key string) (bool, error) {
	err := c.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (c *BadgerCache) Close() error {
	return c.db.Close()
}
