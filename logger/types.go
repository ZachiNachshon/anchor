package logger

type Logger interface {
	Debug(msg string)
	Debugf(format string, args ...interface{})
	Info(msg string)
	Infof(format string, args ...interface{})
	Warning(msg string)
	Warningf(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, a ...interface{})
	Fatal(msg string)
	Fatalf(format string, a ...interface{})
	SetVerbosityLevel(level string) error
}
