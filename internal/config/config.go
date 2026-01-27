package config

import (
	"fmt"
	"os"
	"time"

	"go-monolith-frame/pkg/cache"
	"go-monolith-frame/pkg/logger"
	"go-monolith-frame/pkg/mysql"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig  `yaml:"server"`
	Database mysql.Config  `yaml:"database"`
	Cache    cache.Config  `yaml:"cache"`
	Log      logger.Config `yaml:"log"`
	JWT      JWTConfig     `yaml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr         string        `yaml:"addr"`          // 服务地址，如 :8080
	Mode         string        `yaml:"mode"`          // gin运行模式: debug, release, test
	ReadTimeout  time.Duration `yaml:"read_timeout"`  // 读取超时
	WriteTimeout time.Duration `yaml:"write_timeout"` // 写入超时
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `yaml:"secret"`      // JWT密钥
	ExpireTime time.Duration `yaml:"expire_time"` // 过期时间
}

var globalConfig *Config

// Load 加载配置文件
func Load(env string) (*Config, error) {
	configPath := fmt.Sprintf("configs/config.%s.yaml", env)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}
