package locator

import (
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
)

var CreateFakeLocator = func() *fakeLocatorImpl {
	return &fakeLocatorImpl{}
}

type fakeLocatorImpl struct {
	Locator
	localRepoPath          string
	ScanMock               func(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError
	AnchorFoldersMock      func() []*models.AnchorFolderInfo
	AnchorFolderMock       func(name string) *models.AnchorFolderInfo
	AnchorFoldersAsMapMock func() map[string]*models.AnchorFolderInfo
	AnchorFolderItemsMock  func(parentFolderName string) []*models.AnchorFolderItemInfo
}

func (l *fakeLocatorImpl) Scan(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
	return l.ScanMock(anchorfilesLocalPath, e, pa)
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
