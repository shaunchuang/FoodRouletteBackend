package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// Init 初始化 logger
func Init(level string) error {
	var config zap.Config

	switch level {
	case "debug":
		config = zap.NewDevelopmentConfig()
	case "production":
		config = zap.NewProductionConfig()
	default:
		config = zap.NewDevelopmentConfig()
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// Info 記錄 info 等級日誌
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Error 記錄 error 等級日誌
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Debug 記錄 debug 等級日誌
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn 記錄 warn 等級日誌
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Fatal 記錄 fatal 等級日誌並結束程式
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sync 刷新日誌緩衝區
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}