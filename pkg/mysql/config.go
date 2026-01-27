package mysql

import "time"

// Config 数据库配置
type Config struct {
	DSN             string        `yaml:"dsn"`                // 数据库连接字符串
	MaxOpenConns    int           `yaml:"max_open_conns"`     // 最大打开连接数
	MaxIdleConns    int           `yaml:"max_idle_conns"`     // 最大空闲连接数
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`  // 连接最大生命周期
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"` // 连接最大空闲时间
}
