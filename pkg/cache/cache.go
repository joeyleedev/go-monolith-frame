package cache

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("cache: key not found")

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// 扩展方法
	GetObject(ctx context.Context, key string, dest interface{}) error
	SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// 资源释放
	Close() error
}

// 缓存类型
type Type string

const (
	TypeRedis  Type = "redis"
	TypeMemory Type = "memory"
)

// New 工厂函数，根据配置创建对应实现
func New(cfg *Config) (Cache, error) {
	switch cfg.Type {
	case TypeRedis:
		return NewRedisCache(cfg.Redis)
	case TypeMemory:
		return NewMemoryCache(cfg.Memory)
	default:
		return NewMemoryCache(nil) // 默认内存缓存
	}
}
