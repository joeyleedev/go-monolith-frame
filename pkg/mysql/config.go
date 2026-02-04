package mysql

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	DSN                       string        `mapstructure:"dsn"`                           // 数据库连接字符串
	MaxOpenConns              int           `mapstructure:"max_open_conns"`                // 最大打开连接数
	MaxIdleConns              int           `mapstructure:"max_idle_conns"`                // 最大空闲连接数
	ConnMaxLifetime           time.Duration `mapstructure:"conn_max_lifetime"`             // 连接最大生命周期
	ConnMaxIdleTime           time.Duration `mapstructure:"conn_max_idle_time"`            // 连接最大空闲时间
	LogLevelStr               string        `mapstructure:"log_level"`                     // 日志级别（字符串）
	SlowThreshold             time.Duration `mapstructure:"slow_threshold"`                // 慢查询时间阈值，单位为毫秒
	IgnoreRecordNotFoundError bool          `mapstructure:"ignore_record_not_found_error"` // 是否忽略记录未找到错误
}

// GetLogLevel 将字符串转换为 gorm logger.LogLevel
func (c *Config) GetLogLevel() logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(c.LogLevelStr)) {
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

// Validate 验证配置并设置合理的默认值
func (c *Config) Validate() error {
	// 验证 DSN 不为空
	if c.DSN == "" {
		return fmt.Errorf("database DSN is required")
	}

	// 设置连接池参数的默认值
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = 100
	}
	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnMaxLifetime <= 0 {
		c.ConnMaxLifetime = time.Hour
	}
	if c.ConnMaxIdleTime <= 0 {
		c.ConnMaxIdleTime = 10 * time.Minute
	}

	// 验证慢查询阈值为正数
	if c.SlowThreshold < 0 {
		return fmt.Errorf("slow threshold must be non-negative, got %v", c.SlowThreshold)
	}

	return nil
}
