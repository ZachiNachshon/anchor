package logger

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

const (
	defaultLoggerLogFilePathFormat    = "%s/.config/anchor/anchor.log"
	defaultScriptOutputFilePathFormat = "%s/.config/anchor/scripts-output.log"
)

const (
	Identifier string = "logger-manager"
)

var loggerInUse Logger

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
}

type LoggerManager interface {
	CreateEmptyLogger() (Logger, error)
	AppendStdoutLogger(level string) (Logger, error)
	AppendFileLogger(level string) (Logger, error)
	SetActiveLogger(log *Logger) error
	SetVerbosityLevel(level string) error
	GetDefaultLoggerLogFilePath() (string, error)
}

type loggerManagerImpl struct {
	LoggerManager
	adapter LoggerLogrusAdapter
}

func NewManager() LoggerManager {
	return &loggerManagerImpl{}
}

func (l *loggerManagerImpl) CreateEmptyLogger() (Logger, error) {
	l.adapter = NewAdapter()
	return l.adapter.(Logger), nil
}

func (l *loggerManagerImpl) AppendStdoutLogger(level string) (Logger, error) {
	return l.adapter.AppendStdoutLogger(level)
}

func (l *loggerManagerImpl) AppendFileLogger(level string) (Logger, error) {
	logFilePath, err := l.GetDefaultLoggerLogFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve logger file path. error: %s", err)
	}
	return l.adapter.AppendFileBasedLogger(logFilePath, level)
}

func (l *loggerManagerImpl) SetActiveLogger(log *Logger) error {
	loggerInUse = *log
	return nil
}

func (l *loggerManagerImpl) SetVerbosityLevel(level string) error {
	err := l.adapter.SetStdoutVerbosityLevel(level)
	if err != nil {
		return err
	}
	// TODO: currently we'll keep file logger in debug level
	return nil
}

func (l *loggerManagerImpl) GetDefaultLoggerLogFilePath() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultLoggerLogFilePathFormat, homeFolder), nil
	}
}

func GetDefaultScriptOutputLogFilePath() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultScriptOutputFilePathFormat, homeFolder), nil
	}
}

func FromContext(ctx common.Context) *Logger {
	return ctx.Logger().(*Logger)
}

func SetInContext(ctx common.Context, log *Logger) {
	ctx.(common.LoggerSetter).SetLogger(log)
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
