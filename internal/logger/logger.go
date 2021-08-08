package logger

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

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

const (
	defaultLoggerLogFilePathFormat    = "%s/.config/anchor/anchor.log"
	defaultScriptOutputFilePathFormat = "%s/.config/anchor/scripts-output.log"
)

var GetDefaultLoggerLogFilePath = func() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultLoggerLogFilePathFormat, homeFolder), nil
	}
}

var GetDefaultScriptOutputLogFilePath = func() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultScriptOutputFilePathFormat, homeFolder), nil
	}
}

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
