package logger

import (
	"strings"
	"testing"
)

var CreateFakeTestingLogger = func(t *testing.T, verbose bool) (Logger, error) {
	fakeLogger := &fakeLogger{
		testing: t,
		verbose: verbose,
	}
	return fakeLogger, nil
}

type fakeLogger struct {
	Logger
	LoggerLogrusAdapter

	testing *testing.T
	verbose bool

	AppendStdoutLoggerMock      func(level string) (Logger, error)
	AppendFileBasedLoggerMock   func(filePath string, level string) (Logger, error)
	SetStdoutVerbosityLevelMock func(level string) error // has concrete implementation
	SetFileVerbosityLevelMock   func(level string) error // has concrete implementation
}

func (fl *fakeLogger) AppendStdoutLogger(level string) (Logger, error) {
	return fl.AppendStdoutLoggerMock(level)
}

func (fl *fakeLogger) AppendFileBasedLogger(filePath string, level string) (Logger, error) {
	return fl.AppendFileBasedLoggerMock(filePath, level)
}

func (fl *fakeLogger) SetStdoutVerbosityLevel(level string) error {
	fl.verbose = strings.EqualFold(level, "debug")
	return nil
}

func (fl *fakeLogger) SetFileVerbosityLevel(level string) error {
	fl.verbose = strings.EqualFold(level, "debug")
	return nil
}

func (fl *fakeLogger) Debug(msg string) {
	if fl.verbose {
		fl.testing.Log(msg)
	}
}

func (fl *fakeLogger) Debugf(format string, args ...interface{}) {
	if fl.verbose {
		fl.testing.Logf(format, args...)
	}
}

func (fl *fakeLogger) Info(msg string) {
	if fl.verbose {
		fl.testing.Log(msg)
	}
}

func (fl *fakeLogger) Infof(format string, args ...interface{}) {
	if fl.verbose {
		fl.testing.Logf(format, args...)
	}
}

func (fl *fakeLogger) Warning(msg string) {
	if fl.verbose {
		fl.testing.Log(msg)
	}
}

func (fl *fakeLogger) Warningf(format string, args ...interface{}) {
	if fl.verbose {
		fl.testing.Logf(format, args...)
	}
}

func (fl *fakeLogger) Error(msg string) {
	fl.testing.Logf(msg)
}

func (fl *fakeLogger) Errorf(format string, args ...interface{}) {
	fl.testing.Logf(format, args...)
}

func (fl *fakeLogger) Fatal(msg string) {
	fl.testing.Fatal(msg)
}

func (fl *fakeLogger) Fatalf(format string, args ...interface{}) {
	fl.testing.Fatalf(format, args...)
}
