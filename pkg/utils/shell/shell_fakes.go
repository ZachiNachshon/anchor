package shell

var CreateFakeShell = func() *fakeShell {
	return &fakeShell{}
}

type fakeShell struct {
	Shell
	ExecuteMock                           func(script string) error
	ExecuteWithOutputToFileMock           func(script string, outputFilePath string) error
	ExecuteSilentlyMock                   func(script string) error
	ExecuteScriptFileMock                 func(dir string, relativeScriptPath string, args ...string) error
	ExecuteScriptFileWithOutputToFileMock func(dir string, relativeScriptPath string, outputFilePath string, args ...string) error
	ExecuteWithOutputMock                 func(script string) (string, error)
	ExecuteTTYMock                        func(script string) error
	ExecuteInBackgroundMock               func(script string) error
	ClearScreenMock                       func() error
}

func (s *fakeShell) Execute(script string) error {
	return s.ExecuteMock(script)
}

func (s *fakeShell) ExecuteWithOutputToFile(script string, outputFilePath string) error {
	return s.ExecuteWithOutputToFileMock(script, outputFilePath)
}

func (s *fakeShell) ExecuteSilently(script string) error {
	return s.ExecuteSilentlyMock(script)
}

func (s *fakeShell) ExecuteScriptFile(dir string, relativeScriptPath string, args ...string) error {
	return s.ExecuteScriptFileMock(dir, relativeScriptPath, args...)
}

func (s *fakeShell) ExecuteScriptFileWithOutputToFile(
	workingDirectory string,
	relativeScriptPath string,
	outputFilePath string,
	args ...string) error {

	return s.ExecuteScriptFileWithOutputToFileMock(workingDirectory, relativeScriptPath, outputFilePath, args...)
}

func (s *fakeShell) ExecuteReturnOutput(script string) (string, error) {
	return s.ExecuteWithOutputMock(script)
}

func (s *fakeShell) ExecuteTTY(script string) error {
	return s.ExecuteTTYMock(script)
}

func (s *fakeShell) ExecuteInBackground(script string) error {
	return s.ExecuteInBackgroundMock(script)
}

func (s *fakeShell) ClearScreen() error {
	return s.ClearScreenMock()
}
