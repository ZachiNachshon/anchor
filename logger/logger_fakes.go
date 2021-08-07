package logger

import (
	"testing"
)

var FakeTestingLogger = func(t *testing.T, verbose bool) (Logger, error) {
	fakeLogger := &fakeLogger{
		testing: t,
		verbose: verbose,
	}
	return fakeLogger, nil
}

type fakeLogger struct {
	Logger
	testing *testing.T
	verbose bool
}

func (f *fakeLogger) Debug(msg string) {
	if f.verbose {
		f.testing.Log(msg)
	}
}

func (f *fakeLogger) Debugf(format string, args ...interface{}) {
	if f.verbose {
		f.testing.Logf(format, args...)
	}
}

func (f *fakeLogger) Info(msg string) {
	if f.verbose {
		f.testing.Log(msg)
	}
}

func (f *fakeLogger) Infof(format string, args ...interface{}) {
	if f.verbose {
		f.testing.Logf(format, args...)
	}
}

func (f *fakeLogger) Warning(msg string) {
	if f.verbose {
		f.testing.Log(msg)
	}
}

func (f *fakeLogger) Warningf(format string, args ...interface{}) {
	if f.verbose {
		f.testing.Logf(format, args...)
	}
}

func (f *fakeLogger) Error(msg string) {
	f.testing.Error(msg)
}

func (f *fakeLogger) Errorf(format string, args ...interface{}) {
	f.testing.Errorf(format, args...)
}

func (f *fakeLogger) Fatal(msg string) {
	f.testing.Fatal(msg)
}

func (f *fakeLogger) Fatalf(format string, args ...interface{}) {
	f.testing.Fatalf(format, args...)
}

func (f *fakeLogger) SetVerbosityLevel(level string) error {
	f.verbose = level == "debug"
	return nil
}
