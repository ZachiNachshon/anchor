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
	localRepoPath           string
	ScanMock                func(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError
	CommandFoldersMock      func() []*models.CommandFolderInfo
	CommandFolderByNameMock func(name string) *models.CommandFolderInfo
	CommandFoldersAsMapMock func() map[string]*models.CommandFolderInfo
	CommandFolderItemsMock  func(commandFolderName string) []*models.CommandFolderItemInfo
}

func (l *fakeLocatorImpl) Scan(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
	return l.ScanMock(anchorfilesLocalPath, e, pa)
}

func (l *fakeLocatorImpl) CommandFolders() []*models.CommandFolderInfo {
	return l.CommandFoldersMock()
}

func (l *fakeLocatorImpl) CommandFoldersAsMap() map[string]*models.CommandFolderInfo {
	return l.CommandFoldersAsMapMock()
}

func (l *fakeLocatorImpl) CommandFolderByName(name string) *models.CommandFolderInfo {
	return l.CommandFolderByNameMock(name)
}

func (l *fakeLocatorImpl) CommandFolderItems(commandFolderName string) []*models.CommandFolderItemInfo {
	return l.CommandFolderItemsMock(commandFolderName)
}
