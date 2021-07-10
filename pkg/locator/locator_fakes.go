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
	ScanMock         func(anchorFilesLocalPath string) error
	ApplicationsMock func() []*models.ApplicationInfo
}

func (l *fakeLocatorImpl) LocalRepoPath() string {
	return l.localRepoPath
}

func (l *fakeLocatorImpl) Scan(anchorFilesLocalPath string) error {
	return l.ScanMock(anchorFilesLocalPath)
}

func (l *fakeLocatorImpl) Applications() []*models.ApplicationInfo {
	return l.ApplicationsMock()
}
