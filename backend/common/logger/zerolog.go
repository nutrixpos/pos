// Package logger implements a logger that uses the zerolog library.
//
// It provides a type ZeroLog that implements the logger.ILogger interface.
// It also provides a function NewZeroLog() that returns a new ZeroLog.
//
// The ZeroLog has 3 methods: Info, Warning, and Error. Each of these methods
// logs a message at a different level of severity.
//
// The Info method logs a message at the INFO level of severity.
// The Warning method logs a message at the WARNING level of severity.
// The Error method logs a message at the ERROR level of severity.

package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ZeroLog is a logger that uses the zerolog library.
type ZeroLog struct {
}

// Info logs a message at the INFO level of severity.
func (zl *ZeroLog) Info(message string, fields ...interface{}) {
	log.Info().Msg(message)
}

// Warning logs a message at the WARNING level of severity.
func (zl *ZeroLog) Warning(message string, fields ...interface{}) {
	log.Warn().Msg(message)
}

// Error logs a message at the ERROR level of severity.
func (zl *ZeroLog) Error(message string, fields ...interface{}) {
	log.Error().Msg(message)
}

// NewZeroLog returns a new ZeroLog.
func NewZeroLog() ZeroLog {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := ZeroLog{}

	return log
}
