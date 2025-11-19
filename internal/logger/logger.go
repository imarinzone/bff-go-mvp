package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewProductionLogger creates a production logger
func NewProductionLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config.Build()
}

// NewDevelopmentLogger creates a development logger
func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config.Build()
}

// NewLogger creates a logger based on environment (defaults to production)
func NewLogger(env string) (*zap.Logger, error) {
	if env == "development" || env == "dev" {
		return NewDevelopmentLogger()
	}
	return NewProductionLogger()
}

