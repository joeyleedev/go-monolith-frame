package mysql

import (
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	DSN                       string        `yaml:"dsn"`                           // 数据库连接字符串
	MaxOpenConns              int           `yaml:"max_open_conns"`                // 最大打开连接数
	MaxIdleConns              int           `yaml:"max_idle_conns"`                // 最大空闲连接数
	ConnMaxLifetime           time.Duration `yaml:"conn_max_lifetime"`             // 连接最大生命周期
	ConnMaxIdleTime           time.Duration `yaml:"conn_max_idle_time"`            // 连接最大空闲时间
	LogLevelStr               string        `yaml:"log_level"`                     // 日志级别（字符串）
	SlowThreshold             time.Duration `yaml:"slow_threshold"`                // 慢查询时间阈值，单位为毫秒
	IgnoreRecordNotFoundError bool          `yaml:"ignore_record_not_found_error"` // 是否忽略记录未找到错误
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
