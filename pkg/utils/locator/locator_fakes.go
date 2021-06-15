package locator

import "github.com/ZachiNachshon/anchor/models"

var CreateFakeLocator = func(anchorFilesPath string) *fakeLocatorImpl {
	return &fakeLocatorImpl{
		localRepoPath: anchorFilesPath,
	}
}

type fakeLocatorImpl struct {
	Locator
	localRepoPath    string
	ScanMock         func() error
	ApplicationsMock func() []*models.AppContent
}

func (l *fakeLocatorImpl) LocalRepoPath() string {
	return l.localRepoPath
}

func (l *fakeLocatorImpl) Scan() error {
	return l.ScanMock()
}

func (l *fakeLocatorImpl) Applications() []*models.AppContent {
	return l.ApplicationsMock()
}
