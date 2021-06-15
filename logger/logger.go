package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var loggerInUse Logger

func SetLogger(log Logger) {
	loggerInUse = log
}

func SetVerbosityLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	logrus.SetLevel(lvl)
	return nil
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

var LogrusLoggerLoader = func(verbose bool) error {
	logrusLogger := &logrusLoggerImpl{}

	level := "info"
	if verbose {
		level = "debug"
	}

	if err := SetVerbosityLevel(level); err != nil {
		return err
	}

	// Log as JSON instead of the default ASCII formatter.
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	textFormatter := logrus.TextFormatter{}
	textFormatter.FullTimestamp = true
	textFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(&textFormatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
	//logrus.SetFormatter(&logrus.TextFormatter{
	//	ForceColors: true,
	//	TimestampFormat : "2006-01-02 15:04:05",
	//	FullTimestamp:true,
	//})

	SetLogger(logrusLogger)
	return nil
}

type logrusLoggerImpl struct{}

func (lr *logrusLoggerImpl) Debug(msg string) {
	logrus.Debug(msg)
}

func (lr *logrusLoggerImpl) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func (lr *logrusLoggerImpl) Warning(msg string) {
	logrus.Warning(msg)
}

func (lr *logrusLoggerImpl) Warningf(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

func (lr *logrusLoggerImpl) Info(msg string) {
	logrus.Info(msg)
}

func (lr *logrusLoggerImpl) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (lr *logrusLoggerImpl) Error(msg string) {
	logrus.Error(msg)
}

func (lr *logrusLoggerImpl) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (lr *logrusLoggerImpl) Fatal(msg string) {
	logrus.Fatal(msg)
}

func (lr *logrusLoggerImpl) Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
