package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewZapLogger(t *testing.T) {
	cfg := DefaultConfig()
	logger, _ := NewZapLogger(cfg)

	logger.Info("logger init", zap.String("type", "zap"))
	logger.Error("logger init", zap.String("type", "zap"))
}
