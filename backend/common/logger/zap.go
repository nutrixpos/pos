// Package logger implements a logger that uses the zap library.
//
// It provides a type ZapLogger that implements the logger.ILogger interface.
// It also provides a function NewZapLogger() that returns a new ZapLogger.
//
// The ZapLogger has 3 methods: Info, Warning, and Error. Each of these methods
// logs a message at a different level of severity.
//
// The Info method logs a message at the INFO level of severity.
// The Warning method logs a message at the WARNING level of severity.
// The Error method logs a message at the ERROR level of severity.
package logger

import (
	"go.uber.org/zap"
)

// ZapLogger is a logger that uses the zap library.
type ZapLogger struct {
	// logger is the zap logger.
	logger *zap.Logger
}

// Info logs a message at the INFO level of severity.
func (zl *ZapLogger) Info(message string, fields ...interface{}) {
	// Create a sugar logger which provides a nicer API.
	sugar := zl.logger.Sugar()
	// Log the message and fields at the INFO level.
	sugar.Infow(message, fields)
}

// Warning logs a message at the WARNING level of severity.
func (zl *ZapLogger) Warning(message string, fields ...interface{}) {
	// Create a sugar logger which provides a nicer API.
	sugar := zl.logger.Sugar()
	// Log the message and fields at the WARNING level.
	sugar.Warnw(message, fields)
}

// Error logs a message at the ERROR level of severity.
func (zl *ZapLogger) Error(message string, fields ...interface{}) {
	// Create a sugar logger which provides a nicer API.
	sugar := zl.logger.Sugar()
	// Log the message and fields at the ERROR level.
	sugar.Errorw(message, fields)
}

// NewZapLogger returns a new ZapLogger.
func NewZapLogger() ZapLogger {
	// Create a new zap logger at the production level.
	logger, _ := zap.NewProduction()

	// Create a new ZapLogger and return it.
	zap_logger := ZapLogger{
		logger: logger,
	}

	return zap_logger
}
