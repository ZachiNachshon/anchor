package locator

var FakeLocatorLoader = func(anchorFilesPath string) *fakeLocator {
	return &fakeLocator{
		localRepoPath: anchorFilesPath,
	}
}

type fakeLocator struct {
	Locator
	localRepoPath string
	PrintMock     func()
}

func (l *fakeLocator) LocalRepoPath() string {
	return l.localRepoPath
}

func (l *fakeLocator) Print() {
	l.PrintMock()
}
