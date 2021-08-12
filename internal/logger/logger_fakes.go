package logger

var CreateFakeLoggerManager = func() *fakeLoggerManagerImpl {
	return &fakeLoggerManagerImpl{}
}

type fakeLoggerManagerImpl struct {
	LoggerManager
	CreateEmptyLoggerMock           func() (Logger, error)
	AppendStdoutLoggerMock          func(level string) (Logger, error)
	AppendFileLoggerMock            func(level string) (Logger, error)
	SetActiveLoggerMock             func(log *Logger) error
	SetVerbosityLevelMock           func(level string) error
	GetDefaultLoggerLogFilePathMock func() (string, error)
}

func (e *fakeLoggerManagerImpl) CreateEmptyLogger() (Logger, error) {
	return e.CreateEmptyLoggerMock()
}

func (e *fakeLoggerManagerImpl) AppendStdoutLogger(level string) (Logger, error) {
	return e.AppendStdoutLoggerMock(level)
}

func (e *fakeLoggerManagerImpl) AppendFileLogger(level string) (Logger, error) {
	return e.AppendFileLoggerMock(level)
}

func (e *fakeLoggerManagerImpl) SetActiveLogger(log *Logger) error {
	loggerInUse = *log
	return e.SetActiveLoggerMock(log)
}

func (e *fakeLoggerManagerImpl) SetVerbosityLevel(level string) error {
	return e.SetVerbosityLevelMock(level)
}

func (e *fakeLoggerManagerImpl) GetDefaultLoggerLogFilePath() (string, error) {
	return e.GetDefaultLoggerLogFilePathMock()
}
