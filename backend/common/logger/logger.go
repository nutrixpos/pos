package logger

type ILogger interface {
	Info(string, ...interface{})
	Warning(string, ...interface{})
	Error(string, ...interface{})
}
