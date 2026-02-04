package mysql

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type zapGormLogger struct {
	logger                    *zap.Logger
	logLevel                  logger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func NewZapGormLogger(logger *zap.Logger, level logger.LogLevel, slowThreshold time.Duration, ignoreRecordNotFoundError bool) logger.Interface {
	return &zapGormLogger{
		logger:                    logger,
		logLevel:                  level,
		slowThreshold:             slowThreshold,
		ignoreRecordNotFoundError: ignoreRecordNotFoundError,
	}
}

func (l *zapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		l.logger.Info(msg, zap.Any("data", data))
	}
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.logger.Warn(msg, zap.Any("data", data))
	}
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		l.logger.Error(msg, zap.Any("data", data))
	}
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	logFields := []zap.Field{
		zap.String("duration", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
		zap.String("sql", sql),
	}

	// 提取并记录 context 中的信息
	if ctx != nil {
		if requestID := ctx.Value("request_id"); requestID != nil {
			logFields = append(logFields, zap.Any("request_id", requestID))
		}
	}

	if rows != -1 {
		logFields = append(logFields, zap.Int64("rows", rows))
	}

	switch {
	case err != nil && l.logLevel >= logger.Error && (!l.ignoreRecordNotFoundError || err != logger.ErrRecordNotFound):
		l.logger.Error("sql execution failed", append(logFields, zap.Error(err))...)
	case l.slowThreshold != 0 && elapsed > l.slowThreshold && l.logLevel >= logger.Warn:
		l.logger.Warn("slow sql detected", logFields...)
	case l.logLevel >= logger.Info:
		l.logger.Info("sql executed", logFields...)
	}
}
