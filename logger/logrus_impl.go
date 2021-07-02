package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var LogrusLoggerLoader = func(verbose bool, logFilePath string) error {
	level := "info"
	if verbose {
		level = "debug"
	}

	stdOutLog, err := createStdoutLogger(level)
	if err != nil {
		return err
	}

	// TODO: add retention for xx log files with log rotation to conserve disk space
	//       currently file based logger use debug level for visibility
	fileLog, err := createFileBasedLogger(logFilePath, "debug")
	if err != nil {
		return err
	}

	logrusLogger := &logrusLoggerImpl{
		stdoutLogger: stdOutLog,
		fileLogger:   fileLog,
	}

	SetLogger(logrusLogger)
	return nil
}

func createStdoutLogger(level string) (*logrus.Logger, error) {
	var log = logrus.New()
	log.Out = os.Stdout

	if err := setLogrusVerbosityLevel(level, log); err != nil {
		return nil, fmt.Errorf("failed to set verbosity on stdout logger")
	}

	textFormatter := logrus.TextFormatter{}
	textFormatter.FullTimestamp = true
	textFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.Formatter = &textFormatter

	return log, nil
}

func createFileBasedLogger(filePath string, level string) (*logrus.Logger, error) {
	var log = logrus.New()

	if file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return nil, fmt.Errorf("failed to open logger at path, cannot create file based logger. path: %s", filePath)
	} else {
		log.Out = file
	}

	if err := setLogrusVerbosityLevel(level, log); err != nil {
		return nil, fmt.Errorf("failed to set verbosity on file based logger. path: %s", filePath)
	}

	// Log as JSON instead of the default ASCII formatter.
	jsonFormatter := logrus.JSONFormatter{}
	jsonFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.Formatter = &jsonFormatter

	return log, nil
}

func setLogrusVerbosityLevel(level string, logger *logrus.Logger) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	logger.SetLevel(lvl)
	return nil
}

type logrusLoggerImpl struct {
	Logger
	stdoutLogger *logrus.Logger
	fileLogger   *logrus.Logger
}

func (lr *logrusLoggerImpl) Debug(msg string) {
	lr.stdoutLogger.Debug(msg)
	lr.fileLogger.Debug(msg)
}

func (lr *logrusLoggerImpl) Debugf(format string, args ...interface{}) {
	lr.stdoutLogger.Debugf(format, args...)
	lr.fileLogger.Debugf(format, args...)
}

func (lr *logrusLoggerImpl) Warning(msg string) {
	lr.stdoutLogger.Warning(msg)
	lr.fileLogger.Warning(msg)
}

func (lr *logrusLoggerImpl) Warningf(format string, args ...interface{}) {
	lr.stdoutLogger.Warningf(format, args...)
	lr.fileLogger.Warningf(format, args...)
}

func (lr *logrusLoggerImpl) Info(msg string) {
	lr.stdoutLogger.Info(msg)
	lr.fileLogger.Info(msg)
}

func (lr *logrusLoggerImpl) Infof(format string, args ...interface{}) {
	lr.stdoutLogger.Infof(format, args...)
	lr.fileLogger.Infof(format, args...)
}

func (lr *logrusLoggerImpl) Error(msg string) {
	lr.stdoutLogger.Error(msg)
	lr.fileLogger.Error(msg)
}

func (lr *logrusLoggerImpl) Errorf(format string, args ...interface{}) {
	lr.stdoutLogger.Errorf(format, args...)
	lr.fileLogger.Errorf(format, args...)
}

func (lr *logrusLoggerImpl) Fatal(msg string) {
	lr.stdoutLogger.Fatal(msg)
	lr.fileLogger.Fatal(msg)
}

func (lr *logrusLoggerImpl) Fatalf(format string, args ...interface{}) {
	lr.stdoutLogger.Fatalf(format, args...)
	lr.fileLogger.Fatalf(format, args...)
}

func (f *logrusLoggerImpl) setVerbosityLevel(level string) error {
	if err := setLogrusVerbosityLevel(level, f.stdoutLogger); err != nil {
		return fmt.Errorf("failed to set verbosity on stdout logger")
	}
	if err := setLogrusVerbosityLevel(level, f.fileLogger); err != nil {
		return fmt.Errorf("failed to set verbosity on file based logger")
	}
	return nil
}
