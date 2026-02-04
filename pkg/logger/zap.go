package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewZapLogger creates a new zap logger instance based on configuration
func NewZapLogger(cfg *Config) (*zap.Logger, error) {
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Parse log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	var cores []zapcore.Core

	// File output configuration
	if cfg.Filename != "" {
		fileCore := createFileCore(cfg, level)
		cores = append(cores, fileCore)
	}

	// Console output configuration
	if cfg.Console {
		consoleCore := createConsoleCore(cfg, level)
		cores = append(cores, consoleCore)
	}

	// If no output is configured, default to console
	if len(cores) == 0 {
		consoleCore := createConsoleCore(cfg, level)
		cores = append(cores, consoleCore)
	}

	// Combine multiple cores
	core := zapcore.NewTee(cores...)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), nil
}

// createFileCore creates a core for file output
func createFileCore(cfg *Config, level zapcore.Level) zapcore.Core {
	// Lumberjack log rotation
	lj := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	// File encoder (no colors)
	encoder := zapcore.NewConsoleEncoder(getFileConsoleEncoderConfig())

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(lj),
		zap.NewAtomicLevelAt(level),
	)
}

// createConsoleCore creates a core for console output
func createConsoleCore(cfg *Config, level zapcore.Level) zapcore.Core {
	// Console encoder (with colors if enabled)
	encoder := zapcore.NewConsoleEncoder(getConsoleEncoderConfig(cfg.EnableColor))

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)
}

// getFileConsoleEncoderConfig gets console encoder configuration for file output (no colors)
func getFileConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// getJSONEncoderConfig gets JSON encoder configuration
func getJSONEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// getConsoleEncoderConfig gets console encoder configuration
func getConsoleEncoderConfig(enableColor bool) zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if enableColor {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	return config
}
