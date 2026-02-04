package cache

import "time"

type Config struct {
	Type   Type          `mapstructure:"type"` // redis / memory
	Redis  *RedisConfig  `mapstructure:"redis"`
	Memory *MemoryConfig `mapstructure:"memory"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MemoryConfig struct {
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}
