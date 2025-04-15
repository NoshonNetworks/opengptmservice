package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global logger instance
	log *zap.Logger
)

// NewLogger creates a new logger instance
func NewLogger(level, format string) *zap.Logger {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Build the logger
	var err error
	log, err = config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return log
}

// Get returns the global logger instance
func Get() *zap.Logger {
	return log
}

// Sync flushes any buffered log entries
func Sync() {
	if log != nil {
		log.Sync()
	}
}
