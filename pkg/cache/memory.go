package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type item struct {
	value      []byte
	expiration int64
}

type MemoryCache struct {
	items     sync.Map
	closeChan chan struct{}
}

func NewMemoryCache(cfg *MemoryConfig) (*MemoryCache, error) {
	c := &MemoryCache{
		closeChan: make(chan struct{}),
	}
	go c.cleanupLoop()
	return c, nil
}

func (c *MemoryCache) cleanupLoop() {
	if c.closeChan == nil {
		return
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		case <-c.closeChan:
			return
		}
	}
}

func (c *MemoryCache) cleanupExpired() {
	now := time.Now().UnixNano()
	c.items.Range(func(key, value any) bool {
		it := value.(item)
		if it.expiration > 0 && now > it.expiration {
			c.items.Delete(key)
		}
		return true
	})
}

func (c *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	val, ok := c.items.Load(key)
	if !ok {
		return "", ErrKeyNotFound
	}

	it := val.(item)
	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		c.items.Delete(key)
		return "", ErrKeyNotFound
	}

	return string(it.value), nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var expiration int64
	if exp > 0 {
		expiration = time.Now().Add(exp).UnixNano()
	}

	c.items.Store(key, item{
		value:      data,
		expiration: expiration,
	})
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.items.Delete(key)
	return nil
}

func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	val, ok := c.items.Load(key)
	if !ok {
		return false, nil
	}

	it := val.(item)
	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		c.items.Delete(key)
		return false, nil
	}

	return true, nil
}

func (c *MemoryCache) GetObject(ctx context.Context, key string, dest interface{}) error {
	data, err := c.GetBytes(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *MemoryCache) SetObject(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return c.Set(ctx, key, value, exp)
}

func (c *MemoryCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	val, ok := c.items.Load(key)
	if !ok {
		return nil, ErrKeyNotFound
	}

	it := val.(item)
	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		c.items.Delete(key)
		return nil, ErrKeyNotFound
	}

	return it.value, nil
}

func (c *MemoryCache) Close() error {
	if c.closeChan != nil {
		close(c.closeChan)
		c.closeChan = nil
	}
	c.items.Range(func(key, _ any) bool {
		c.items.Delete(key)
		return true
	})
	return nil
}
