package mysql

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm/logger"
)

// 创建一个可观察的zap logger用于测试
func createObservableLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.DebugLevel)
	return zap.New(core), logs
}

func TestNewZapGormLogger(t *testing.T) {
	zapLogger, _ := createObservableLogger()
	gormLogger := NewZapGormLogger(zapLogger, logger.Info, time.Second, true)

	if gormLogger == nil {
		t.Fatal("Expected non-nil logger")
	}

	_, ok := gormLogger.(*zapGormLogger)
	if !ok {
		t.Fatal("Expected logger to be of type *zapGormLogger")
	}
}

func TestZapGormLogger_LogMode(t *testing.T) {
	zapLogger, _ := createObservableLogger()
	gormLogger := NewZapGormLogger(zapLogger, logger.Info, time.Second, true).(*zapGormLogger)

	newLogger := gormLogger.LogMode(logger.Warn)
	newZapGormLogger, ok := newLogger.(*zapGormLogger)
	if !ok {
		t.Fatal("Expected logger to be of type *zapGormLogger")
	}

	if newZapGormLogger.logLevel != logger.Warn {
		t.Errorf("Expected log level to be Warn, got %v", newZapGormLogger.logLevel)
	}

	// 验证原logger未被修改
	if gormLogger.logLevel != logger.Info {
		t.Errorf("Expected original logger level to remain Info, got %v", gormLogger.logLevel)
	}
}

func TestZapGormLogger_Trace(t *testing.T) {
	tests := []struct {
		name                      string
		logLevel                  logger.LogLevel
		slowThreshold             time.Duration
		ignoreRecordNotFoundError bool
		err                       error
		elapsed                   time.Duration
		sql                       string
		rows                      int64
		expectedLogLevel          zapcore.Level
		expectedMessage           string
		wantLog                   bool
	}{
		{
			name:             "Successful query with Info level",
			logLevel:         logger.Info,
			slowThreshold:    time.Second,
			err:              nil,
			elapsed:          100 * time.Millisecond,
			sql:              "SELECT * FROM users",
			rows:             10,
			expectedLogLevel: zapcore.InfoLevel,
			expectedMessage:  "sql executed",
			wantLog:          true,
		},
		{
			name:             "Slow query warning",
			logLevel:         logger.Warn,
			slowThreshold:    50 * time.Millisecond,
			err:              nil,
			elapsed:          100 * time.Millisecond,
			sql:              "SELECT * FROM users",
			rows:             10,
			expectedLogLevel: zapcore.WarnLevel,
			expectedMessage:  "slow sql detected",
			wantLog:          true,
		},
		{
			name:             "Error query",
			logLevel:         logger.Error,
			slowThreshold:    time.Second,
			err:              errors.New("connection failed"),
			elapsed:          100 * time.Millisecond,
			sql:              "SELECT * FROM users",
			rows:             -1,
			expectedLogLevel: zapcore.ErrorLevel,
			expectedMessage:  "sql execution failed",
			wantLog:          true,
		},
		{
			name:                      "RecordNotFound error ignored",
			logLevel:                  logger.Error,
			slowThreshold:             time.Second,
			ignoreRecordNotFoundError: true,
			err:                       logger.ErrRecordNotFound,
			elapsed:                   100 * time.Millisecond,
			sql:                       "SELECT * FROM users WHERE id = 1",
			rows:                      0,
			wantLog:                   false,
		},
		{
			name:                      "RecordNotFound error not ignored",
			logLevel:                  logger.Error,
			slowThreshold:             time.Second,
			ignoreRecordNotFoundError: false,
			err:                       logger.ErrRecordNotFound,
			elapsed:                   100 * time.Millisecond,
			sql:                       "SELECT * FROM users WHERE id = 1",
			rows:                      0,
			expectedLogLevel:          zapcore.ErrorLevel,
			expectedMessage:           "sql execution failed",
			wantLog:                   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zapLogger, logs := createObservableLogger()
			gormLogger := NewZapGormLogger(zapLogger, tt.logLevel, tt.slowThreshold, tt.ignoreRecordNotFoundError)

			begin := time.Now().Add(-tt.elapsed)
			fc := func() (string, int64) {
				return tt.sql, tt.rows
			}

			gormLogger.Trace(context.Background(), begin, fc, tt.err)

			if tt.wantLog {
				if logs.Len() == 0 {
					t.Error("Expected log entry, got none")
				} else {
					entry := logs.All()[0]
					if entry.Message != tt.expectedMessage {
						t.Errorf("Expected message %q, got %q", tt.expectedMessage, entry.Message)
					}
					if entry.Level != tt.expectedLogLevel {
						t.Errorf("Expected level %v, got %v", tt.expectedLogLevel, entry.Level)
					}

					// 验证日志字段
					fields := entry.Context
					foundSQL := false
					foundDuration := false
					for _, field := range fields {
						if field.Key == "sql" && field.String == tt.sql {
							foundSQL = true
						}
						if field.Key == "duration" {
							foundDuration = true
						}
					}
					if !foundSQL {
						t.Errorf("Expected SQL field with value %q", tt.sql)
					}
					if !foundDuration {
						t.Error("Expected duration field")
					}
				}
			} else {
				if logs.Len() > 0 {
					t.Errorf("Expected no log entry, got %d", logs.Len())
				}
			}
		})
	}
}
