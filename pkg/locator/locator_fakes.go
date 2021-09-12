package locator

import (
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/pkg/models"
)

var CreateFakeLocator = func(anchorFilesPath string) *fakeLocatorImpl {
	return &fakeLocatorImpl{
		localRepoPath: anchorFilesPath,
	}
}

type fakeLocatorImpl struct {
	Locator
	localRepoPath          string
	ScanMock               func(anchorFilesLocalPath string) *errors.LocatorError
	AnchorFoldersMock      func() []*models.AnchorFolderInfo
	AnchorFolderMock       func(name string) *models.AnchorFolderInfo
	AnchorFoldersAsMapMock func() map[string]*models.AnchorFolderInfo
	AnchorFolderItemsMock  func(parentFolderName string) []*models.AnchorFolderItemInfo
}

func (l *fakeLocatorImpl) LocalRepoPath() string {
	return l.localRepoPath
}

func (l *fakeLocatorImpl) Scan(anchorFilesLocalPath string) *errors.LocatorError {
	return l.ScanMock(anchorFilesLocalPath)
}

func (l *fakeLocatorImpl) AnchorFolders() []*models.AnchorFolderInfo {
	return l.AnchorFoldersMock()
}

func (l *fakeLocatorImpl) AnchorFoldersAsMap() map[string]*models.AnchorFolderInfo {
	return l.AnchorFoldersAsMapMock()
}

func (l *fakeLocatorImpl) AnchorFolder(name string) *models.AnchorFolderInfo {
	return l.AnchorFolderMock(name)
}

func (l *fakeLocatorImpl) AnchorFolderItems(parentFolderName string) []*models.AnchorFolderItemInfo {
	return l.AnchorFolderItemsMock(parentFolderName)
}
