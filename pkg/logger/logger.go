package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string, format string) error {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Set the log level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Customize the encoder config
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Build the logger
	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	// Replace the global logger
	zap.ReplaceGlobals(log)

	return nil
}

func Get() *zap.Logger {
	if log == nil {
		// Initialize with default values if not initialized
		_ = Init("info", "json")
	}
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
