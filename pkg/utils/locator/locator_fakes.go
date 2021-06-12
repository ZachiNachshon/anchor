package locator

import "github.com/ZachiNachshon/anchor/models"

var FakeLocatorLoader = func(anchorFilesPath string) *fakeLocator {
	return &fakeLocator{
		localRepoPath: anchorFilesPath,
	}
}

type fakeLocator struct {
	Locator
	localRepoPath    string
	ScanMock         func() error
	ApplicationsMock func() []*models.AppContent
}

func (l *fakeLocator) LocalRepoPath() string {
	return l.localRepoPath
}

func (l *fakeLocator) Scan() error {
	return l.ScanMock()
}

func (l *fakeLocator) Applications() []*models.AppContent {
	return l.ApplicationsMock()
}
