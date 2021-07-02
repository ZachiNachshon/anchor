package logger

import (
	"github.com/ZachiNachshon/anchor/common"
	"testing"
)

var FakeTestingLogger = func(ctx common.Context, t *testing.T, verbose bool) (Logger, error) {
	fakeLogger := &FakeLogger{
		testing: t,
		verbose: verbose,
	}
	return fakeLogger, nil
}

type FakeLogger struct {
	testing *testing.T
	verbose bool
}

func (f *FakeLogger) Debug(msg string) {
	if f.verbose {
		f.testing.Log(msg)
	}
}

func (f *FakeLogger) Debugf(format string, args ...interface{}) {
	if f.verbose {
		f.testing.Logf(format, args...)
	}
}

func (f *FakeLogger) Info(msg string) {
	if f.verbose {
		f.testing.Log(msg)
	}
}

func (f *FakeLogger) Infof(format string, args ...interface{}) {
	if f.verbose {
		f.testing.Logf(format, args...)
	}
}

func (f *FakeLogger) Warning(msg string) {
	if f.verbose {
		f.testing.Log(msg)
	}
}

func (f *FakeLogger) Warningf(format string, args ...interface{}) {
	if f.verbose {
		f.testing.Logf(format, args...)
	}
}

func (f *FakeLogger) Error(msg string) {
	f.testing.Error(msg)
}

func (f *FakeLogger) Errorf(format string, args ...interface{}) {
	f.testing.Errorf(format, args...)
}

func (f *FakeLogger) Fatal(msg string) {
	f.testing.Fatal(msg)
}

func (f *FakeLogger) Fatalf(format string, args ...interface{}) {
	f.testing.Fatalf(format, args...)
}
