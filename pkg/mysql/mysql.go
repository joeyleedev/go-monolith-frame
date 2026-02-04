package mysql

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// New 初始化数据库连接，返回 GORM 实例
func New(cfg *Config) (*gorm.DB, error) {
	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: NewZapGormLogger(
			zap.L(),
			cfg.GetLogLevel(),
			cfg.SlowThreshold,
			cfg.IgnoreRecordNotFoundError,
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.AutoMigrate(models...)
}

// Close 关闭数据库连接
func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping 检测数据库连接状态
func Ping(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Ping()
}

// Stats 返回连接池状态统计信息
func Stats(db *gorm.DB) map[string]interface{} {
	if db == nil {
		return map[string]interface{}{
			"error": "database not initialized",
		}
	}
	sqlDB, err := db.DB()
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("failed to get sql.DB: %v", err),
		}
	}
	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
