package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// Init initializes the logger
func Init(config *Config) error {
	logger, err := NewZapLogger(config)
	if err != nil {
		return err
	}

	globalLogger = logger
	return nil
}

// getLogger gets the logger instance
func getLogger() *zap.Logger {
	if globalLogger == nil {
		once.Do(func() {
			// Initialize with default configuration
			config := DefaultConfig()
			globalLogger, _ = NewZapLogger(config)
		})
	}
	return globalLogger
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return getLogger()
}

// Sync syncs the log buffer
func Sync() error {
	if logger := getLogger(); logger != nil {
		return logger.Sync()
	}
	return nil
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

// Fatal logs a fatal message
func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}
