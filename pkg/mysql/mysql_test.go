package mysql

import (
	"testing"
	"time"

	"gorm.io/gorm/logger"
)

// TestConfigValidate 测试配置验证功能
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: &Config{
				DSN:             "user:pass@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=True&loc=Local",
				MaxOpenConns:    10,
				MaxIdleConns:    5,
				ConnMaxLifetime: time.Hour,
				ConnMaxIdleTime: 10 * time.Minute,
				SlowThreshold:   200 * time.Millisecond,
			},
			wantErr: false,
		},
		{
			name: "empty DSN should error",
			config: &Config{
				DSN:             "",
				MaxOpenConns:    10,
				MaxIdleConns:    5,
				ConnMaxLifetime: time.Hour,
				ConnMaxIdleTime: 10 * time.Minute,
			},
			wantErr:     true,
			errContains: "DSN is required",
		},
		{
			name: "negative slow threshold should error",
			config: &Config{
				DSN:             "user:pass@tcp(localhost:3306)/db",
				MaxOpenConns:    10,
				MaxIdleConns:    5,
				ConnMaxLifetime: time.Hour,
				ConnMaxIdleTime: 10 * time.Minute,
				SlowThreshold:   -1 * time.Millisecond,
			},
			wantErr:     true,
			errContains: "slow threshold must be non-negative",
		},
		{
			name: "zero values should be set to defaults",
			config: &Config{
				DSN: "user:pass@tcp(localhost:3306)/db",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 记录修改前的值用于验证默认值设置
			originalConfig := *tt.config

			err := tt.config.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			// 验证默认值设置
			if originalConfig.MaxOpenConns <= 0 && tt.config.MaxOpenConns != 100 {
				t.Errorf("Validate() MaxOpenConns = %v, want 100", tt.config.MaxOpenConns)
			}
			if originalConfig.MaxIdleConns <= 0 && tt.config.MaxIdleConns != 10 {
				t.Errorf("Validate() MaxIdleConns = %v, want 10", tt.config.MaxIdleConns)
			}
			if originalConfig.ConnMaxLifetime <= 0 && tt.config.ConnMaxLifetime != time.Hour {
				t.Errorf("Validate() ConnMaxLifetime = %v, want 1h", tt.config.ConnMaxLifetime)
			}
			if originalConfig.ConnMaxIdleTime <= 0 && tt.config.ConnMaxIdleTime != 10*time.Minute {
				t.Errorf("Validate() ConnMaxIdleTime = %v, want 10m", tt.config.ConnMaxIdleTime)
			}
		})
	}
}

// TestGetLogLevel 测试日志级别转换
func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		levelStr string
		want     logger.LogLevel
	}{
		{"warn", "warn", logger.Warn},
		{"WARN", "WARN", logger.Warn},
		{"  warn  ", "  warn  ", logger.Warn},
		{"info", "info", logger.Info},
		{"INFO", "INFO", logger.Info},
		{"unknown", "unknown", logger.Info},
		{"empty", "", logger.Info},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{LogLevelStr: tt.levelStr}
			if got := cfg.GetLogLevel(); got != tt.want {
				t.Errorf("GetLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStats 测试连接池统计功能
func TestStats(t *testing.T) {
	// 测试数据库未初始化的情况
	stats := Stats(nil)
	if stats == nil {
		t.Fatal("Stats() should return a map even when DB is nil")
	}
	if err, ok := stats["error"]; !ok || err == nil {
		t.Errorf("Stats() should contain error key when DB is nil, got: %v", stats)
	}
}

// TestPing 测试连接检查功能
func TestPing(t *testing.T) {
	// 测试数据库未初始化的情况
	err := Ping(nil)
	if err == nil {
		t.Error("Ping() should return error when DB is nil")
	}
}

// contains 检查字符串是否包含子字符串（辅助函数）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
