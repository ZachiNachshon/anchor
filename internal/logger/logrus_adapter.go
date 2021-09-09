package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

type LoggerLogrusAdapter interface {
	AppendStdoutLogger(level string) (Logger, error)
	AppendFileBasedLogger(filePath string, level string) (Logger, error)
	SetStdoutVerbosityLevel(level string) error
	SetFileVerbosityLevel(level string) error
}

type logrusAdapterImpl struct {
	Logger
	LoggerLogrusAdapter
	// TODO: consider using logrus hook with different formatters - https://github.com/sirupsen/logrus/issues/784#issuecomment-403765306
	stdoutLogger *logrus.Logger
	fileLogger   *logrus.Logger
}

func NewAdapter() LoggerLogrusAdapter {
	return &logrusAdapterImpl{
		stdoutLogger: logrus.New(),
		fileLogger:   logrus.New(),
	}
}

func (lr *logrusAdapterImpl) AppendStdoutLogger(level string) (Logger, error) {
	lr.stdoutLogger.Out = os.Stdout
	if err := lr.SetStdoutVerbosityLevel(level); err != nil {
		return nil, fmt.Errorf("failed to set verbosity on stdout logger")
	}

	textFormatter := logrus.TextFormatter{}
	//textFormatter.FullTimestamp = true
	//textFormatter.TimestampFormat = "2006-01-02 15:04:05"
	textFormatter.DisableTimestamp = true
	lr.stdoutLogger.Formatter = &textFormatter

	return lr, nil
}

func (lr *logrusAdapterImpl) AppendFileBasedLogger(filePath string, level string) (Logger, error) {
	if file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return nil, fmt.Errorf("failed to open logger at path, cannot create file based logger. path: %s", filePath)
	} else {
		lr.fileLogger.Out = file
	}

	if err := lr.SetFileVerbosityLevel(level); err != nil {
		return nil, fmt.Errorf("failed to set verbosity on file based logger. path: %s", filePath)
	}

	// Log as JSON instead of the default ASCII formatter.
	jsonFormatter := logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// TODO: check why file/line num. fails to print to log file
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: true,
	}
	jsonFormatter.TimestampFormat = "2006-01-02 15:04:05"
	lr.fileLogger.Formatter = &jsonFormatter

	return lr, nil
}

func (lr *logrusAdapterImpl) SetVerbosityLevel(level string) error {
	err := lr.SetStdoutVerbosityLevel(level)
	if err != nil {
		return err
	}
	err = lr.SetFileVerbosityLevel(level)
	if err != nil {
		return err
	}
	return nil
}

func (lr *logrusAdapterImpl) SetStdoutVerbosityLevel(level string) error {
	if lvl, err := logrus.ParseLevel(level); err != nil {
		return err
	} else {
		lr.stdoutLogger.SetLevel(lvl)
		return nil
	}
}

func (lr *logrusAdapterImpl) SetFileVerbosityLevel(level string) error {
	if lvl, err := logrus.ParseLevel(level); err != nil {
		return err
	} else {
		lr.fileLogger.SetLevel(lvl)
		return nil
	}
}

func (lr *logrusAdapterImpl) Debug(msg string) {
	//lr.stdoutLogger.Debug(msg)
	lr.fileLogger.Debug(msg)
}

func (lr *logrusAdapterImpl) Debugf(format string, args ...interface{}) {
	//lr.stdoutLogger.Debugf(format, args...)
	lr.fileLogger.Debugf(format, args...)
}

func (lr *logrusAdapterImpl) Info(msg string) {
	//lr.stdoutLogger.Info(msg)
	lr.fileLogger.Info(msg)
}

func (lr *logrusAdapterImpl) Infof(format string, args ...interface{}) {
	//lr.stdoutLogger.Infof(format, args...)
	lr.fileLogger.Infof(format, args...)
}

func (lr *logrusAdapterImpl) Warning(msg string) {
	//lr.stdoutLogger.Warning(msg)
	lr.fileLogger.Warning(msg)
}

func (lr *logrusAdapterImpl) Warningf(format string, args ...interface{}) {
	//lr.stdoutLogger.Warningf(format, args...)
	lr.fileLogger.Warningf(format, args...)
}

func (lr *logrusAdapterImpl) Error(msg string) {
	//lr.stdoutLogger.Error(msg)
	lr.fileLogger.Error(msg)
}

func (lr *logrusAdapterImpl) Errorf(format string, args ...interface{}) {
	//lr.stdoutLogger.Errorf(format, args...)
	lr.fileLogger.Errorf(format, args...)
}

func (lr *logrusAdapterImpl) Fatal(msg string) {
	//lr.stdoutLogger.Fatal(msg)
	lr.fileLogger.Fatal(msg)
}

func (lr *logrusAdapterImpl) Fatalf(format string, args ...interface{}) {
	//lr.stdoutLogger.Fatalf(format, args...)
	lr.fileLogger.Fatalf(format, args...)
}
