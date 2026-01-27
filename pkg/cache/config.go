package cache

import "time"

type Config struct {
	Type   Type          `yaml:"type"` // redis / memory
	Redis  *RedisConfig  `yaml:"redis"`
	Memory *MemoryConfig `yaml:"memory"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MemoryConfig struct {
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
}
