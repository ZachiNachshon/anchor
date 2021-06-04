package shell

var FakeShellLoader = func() *fakeShell {
	return &fakeShell{}
}

type fakeShell struct {
	Shell
	ExecuteMock func(script string) error
}

func (s *fakeShell) Execute(script string) error {
	return s.ExecuteMock(script)
}
