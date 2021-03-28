package log

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}
func NewLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
