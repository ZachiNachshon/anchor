package logger

var loggerInUse Logger

func SetLogger(log Logger) {
	loggerInUse = log
}

func Debug(msg string) {
	loggerInUse.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	loggerInUse.Debugf(format, args...)
}

func Info(msg string) {
	loggerInUse.Info(msg)
}

func Infof(format string, args ...interface{}) {
	loggerInUse.Infof(format, args...)
}

func Warning(msg string) {
	loggerInUse.Warning(msg)
}

func Warningf(format string, a ...interface{}) {
	loggerInUse.Warningf(format, a...)
}

func Error(msg string) {
	loggerInUse.Error(msg)
}

func Errorf(format string, a ...interface{}) {
	loggerInUse.Errorf(format, a...)
}

func Fatal(msg string) {
	loggerInUse.Fatal(msg)
}

func Fatalf(format string, a ...interface{}) {
	loggerInUse.Fatalf(format, a...)
}
