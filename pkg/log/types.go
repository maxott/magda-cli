package log

type Tags map[string]interface{}

// Log interface for different log levels.
type Logger interface {
	Debug(msg string)
	Debugf(format string, args ...interface{})
	
	Info(msg string)
	Infof(format string, args ...interface{})

	Warn(msg string)
	Warnf(format string, args ...interface{})

	Error(err error, msg string) error
	Errorf(err error, format string, args ...interface{}) error

	Fatal(msg string)
	Fatalf(format string, args ...interface{})

	With(key string, value interface{}) Logger
	WithTags(tags Tags) Logger
}
