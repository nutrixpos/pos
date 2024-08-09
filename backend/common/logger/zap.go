package logger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func (zl *ZapLogger) Info(message string, fields ...interface{}) {
	sugar := zl.logger.Sugar()
	sugar.Infow(message, fields)
}

func (zl *ZapLogger) Warning(message string, fields ...interface{}) {
	sugar := zl.logger.Sugar()
	sugar.Warnw(message, fields)
}

func (zl *ZapLogger) Error(message string, fields ...interface{}) {
	sugar := zl.logger.Sugar()
	sugar.Errorw(message, fields)
}

func NewZapLogger() ZapLogger {
	logger, _ := zap.NewProduction()

	zap_logger := ZapLogger{
		logger: logger,
	}

	return zap_logger
}
