package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *zap.Logger
	sugarLogger  *zap.SugaredLogger
	once         sync.Once
)

// Init 初始化日志系统
func Init(cfg *Config) error {
	var err error
	once.Do(func() {
		err = initLogger(cfg)
	})
	return err
}

func initLogger(cfg *Config) error {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel // 默认info级别
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
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

	// 文件写入器
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	// 构建核心
	var core zapcore.Core
	if cfg.Console {
		// 同时输出到文件和控制台
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleWriter := zapcore.AddSync(os.Stdout)
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, level),
			zapcore.NewCore(consoleEncoder, consoleWriter, level),
		)
	} else {
		// 仅输出到文件
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, level)
	}

	// 创建logger
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	sugarLogger = globalLogger.Sugar()

	return nil
}

// Debug 记录debug级别日志
func Debug(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Debug(msg, fields...)
	}
}

// Info 记录info级别日志
func Info(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Info(msg, fields...)
	}
}

// Warn 记录warn级别日志
func Warn(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Warn(msg, fields...)
	}
}

// Error 记录error级别日志
func Error(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Error(msg, fields...)
	}
}

// Fatal 记录fatal级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Fatal(msg, fields...)
	}
}

// With 创建带字段的logger
func With(fields ...zap.Field) *zap.Logger {
	if globalLogger != nil {
		return globalLogger.With(fields...)
	}
	return nil
}

// Sync 同步日志缓冲区
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// GetLogger 获取原始zap.Logger
func GetLogger() *zap.Logger {
	return globalLogger
}

// GetSugarLogger 获取SugaredLogger
func GetSugarLogger() *zap.SugaredLogger {
	return sugarLogger
}
