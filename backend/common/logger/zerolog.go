package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZeroLog struct {
}

func (zl *ZeroLog) Info(message string, fields ...interface{}) {
	log.Info().Msg(message)
}

func (zl *ZeroLog) Warning(message string, fields ...interface{}) {
	log.Warn().Msg(message)

}

func (zl *ZeroLog) Error(message string, fields ...interface{}) {
	log.Error().Msg(message)
}

func NewZeroLog() ZeroLog {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := ZeroLog{}

	return log
}
