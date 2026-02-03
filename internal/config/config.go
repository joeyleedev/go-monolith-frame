package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go-monolith-frame/pkg/cache"
	"go-monolith-frame/pkg/logger"
	"go-monolith-frame/pkg/mysql"
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
	v := viper.New()

	// 1. 设置配置文件路径和名称
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")

	// 2. 启用自动环境变量读取
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 3. 读取配置文件（作为默认值）
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 4. 反序列化到 Config 结构
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}
