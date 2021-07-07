package shell

var CreateFakeShell = func() *fakeShell {
	return &fakeShell{}
}

type fakeShell struct {
	Shell
	ExecuteMock             func(script string) error
	ExecuteSilentlyMock     func(script string) error
	ExecuteScriptMock       func(dir string, relativeScriptPath string, args ...string) error
	ExecuteWithOutputMock   func(script string) (string, error)
	ExecuteTTYMock          func(script string) error
	ExecuteInBackgroundMock func(script string) error
}

func (s *fakeShell) Execute(script string) error {
	return s.ExecuteMock(script)
}

func (s *fakeShell) ExecuteSilently(script string) error {
	return s.ExecuteSilentlyMock(script)
}

func (s *fakeShell) ExecuteScript(dir string, relativeScriptPath string, args ...string) error {
	return s.ExecuteScriptMock(dir, relativeScriptPath, args...)
}

func (s *fakeShell) ExecuteWithOutput(script string) (string, error) {
	return s.ExecuteWithOutputMock(script)
}

func (s *fakeShell) ExecuteTTY(script string) error {
	return s.ExecuteTTYMock(script)
}

func (s *fakeShell) ExecuteInBackground(script string) error {
	return s.ExecuteInBackgroundMock(script)
}
