package logger

// ILogger defines an interface for logging with different levels of severity.
type ILogger interface {
	// Info logs informational messages.
	Info(string, ...interface{})
	// Warning logs warning messages.
	Warning(string, ...interface{})
	// Error logs error messages.
	Error(string, ...interface{})
}
