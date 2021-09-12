package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/atomics"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	Identifier string = "locator"
)

type Locator interface {
	Scan(anchorFilesLocalPath string) *errors.LocatorError
	AnchorFolders() []*models.AnchorFolderInfo
	AnchorFolder(name string) *models.AnchorFolderInfo
	AnchorFoldersAsMap() map[string]*models.AnchorFolderInfo
	AnchorFolderItems(parentFolderName string) []*models.AnchorFolderItemInfo
}

var excludedDirectories map[string]bool

func init() {
	excludedDirectories = make(map[string]bool)
	excludedDirectories[".git"] = true
	excludedDirectories[".idea"] = true
	excludedDirectories[".DS_Store"] = true
	excludedDirectories[".gradle"] = true
	excludedDirectories["build"] = true
	excludedDirectories["out"] = true
	excludedDirectories["target"] = true
	excludedDirectories["logs"] = true
}

const (
	anchorIgnoreFileName string = ".anchorignore"
)

type locatorImpl struct {
	Locator

	initialized                       atomics.AtomicBool
	anchorFoldersInfoMap              map[string]*models.AnchorFolderInfo
	longestAnchorFolderNameLength     int
	longestAnchorFolderItemNameLength int
}

func newAnchorFolderItem(name string, dirPath string) *models.AnchorFolderItemInfo {
	return &models.AnchorFolderItemInfo{
		Name:             name,
		DirPath:          dirPath,
		InstructionsPath: fmt.Sprintf("%s/%s", dirPath, globals.InstructionsFileName),
	}
}

func New() *locatorImpl {
	return &locatorImpl{
		anchorFoldersInfoMap: make(map[string]*models.AnchorFolderInfo),
	}
}

func (l *locatorImpl) Scan(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
	if l.isInitialized() {
		warnMsg := "scan for instructions can be called only once, using previous scan result"
		logger.Warning(warnMsg)
		return errors.NewAlreadyInitializedError(fmt.Errorf(warnMsg))
	}

	if !ioutils.IsValidPath(anchorfilesLocalPath) {
		return errors.NewLocatorError(fmt.Errorf("invalid anchorfile local path. path: %s", anchorfilesLocalPath))
	}

	err := scanForAnchorFolders(l, anchorfilesLocalPath, e, pa)
	if err != nil {
		return errors.NewLocatorError(err)
	}

	l.markInitialized()
	return nil
}

func (l *locatorImpl) AnchorFolders() []*models.AnchorFolderInfo {
	res := make([]*models.AnchorFolderInfo, 0, len(l.anchorFoldersInfoMap))
	for _, v := range l.anchorFoldersInfoMap {
		res = append(res, v)
	}
	return sortAnchorFolders(res)
}

func (l *locatorImpl) AnchorFolder(name string) *models.AnchorFolderInfo {
	if value, exists := l.anchorFoldersInfoMap[name]; exists {
		return value
	}
	return nil
}

func (l *locatorImpl) AnchorFoldersAsMap() map[string]*models.AnchorFolderInfo {
	return l.anchorFoldersInfoMap
}

func (l *locatorImpl) AnchorFolderItems(parentFolderName string) []*models.AnchorFolderItemInfo {
	if l.anchorFoldersInfoMap[parentFolderName] != nil &&
		l.anchorFoldersInfoMap[parentFolderName].Items != nil {
		items := l.anchorFoldersInfoMap[parentFolderName].Items
		result := make([]*models.AnchorFolderItemInfo, len(items))
		i := 0
		for _, item := range items {
			result[i] = item
			i++
		}
		return result
	}
	return nil
}

func (l *locatorImpl) isInitialized() bool {
	return l.initialized.Get()
}

func (l *locatorImpl) markInitialized() {
	l.initialized.Set(true)
}

func (l *locatorImpl) appendAnchorFolder(anchorFolder *models.AnchorFolderInfo) {
	l.anchorFoldersInfoMap[anchorFolder.Name] = anchorFolder
}

func (l *locatorImpl) appendAnchorFolderItem(anchorFolder *models.AnchorFolderInfo, anchorFolderItem *models.AnchorFolderItemInfo) {
	if anchorFolder.Items == nil {
		anchorFolder.Items = make(map[string]*models.AnchorFolderItemInfo)
	}
	anchorFolder.Items[anchorFolderItem.Name] = anchorFolderItem
}

// Longest anchor folder name is Used for prompter left padding
func (l *locatorImpl) setLongestAnchorFolderName(name string) {
	if l.longestAnchorFolderNameLength < len(name) {
		l.longestAnchorFolderNameLength = len(name)
	}
}

// Longest anchor folder item name is Used for prompter left padding
func (l *locatorImpl) setLongestAnchorFolderItemName(name string) {
	if l.longestAnchorFolderItemNameLength < len(name) {
		l.longestAnchorFolderItemNameLength = len(name)
	}
}

func scanForAnchorFolders(l *locatorImpl, anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) error {
	filepath.Walk(anchorfilesLocalPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignore root directory
			if info.IsDir() && path == anchorfilesLocalPath {
				return nil
			}

			// Ignore all hidden directories
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			// Ignore all files under root folder
			dir := filepath.Dir(path)
			if !info.IsDir() && dir == anchorfilesLocalPath {
				return nil
			}

			// Skip if dir/file should be excluded
			name := info.Name()
			if isExcluded(name) {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			if _, ok := hasAnchorIgnoreIdentifier(path); ok {
				return filepath.SkipDir
			}

			resolvedAnchorFolder := tryResolveAnchorFolder(l, path, name, e, pa)
			if resolvedAnchorFolder != nil {
				scanForAnchorFolderItems(l, resolvedAnchorFolder)
			}

			return nil
		})

	return nil
}

func tryResolveAnchorFolder(l *locatorImpl, dirPath string, name string, e extractor.Extractor, pa parser.Parser) *models.AnchorFolderInfo {
	if isAnchorFolder(dirPath) {
		logger.Debugf("Locate anchor folder. name: %s", name)
		if anchorFolder, err := e.ExtractAnchorFolderInfo(dirPath, pa); err != nil {
			return nil
		} else {
			l.appendAnchorFolder(anchorFolder)
			l.setLongestAnchorFolderName(name)
			return anchorFolder
		}
	}
	return nil
}

func scanForAnchorFolderItems(l *locatorImpl, anchorFolder *models.AnchorFolderInfo) {
	filepath.Walk(anchorFolder.DirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignore all hidden directories
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			// Skip if dir/file should be excluded
			name := info.Name()
			if isExcluded(name) {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			if _, ok := hasAnchorIgnoreIdentifier(path); ok {
				return filepath.SkipDir
			}

			resolvedAnchorFolderItem := tryResolveAnchorFolderItems(l, path, name)
			if resolvedAnchorFolderItem != nil {
				l.appendAnchorFolderItem(anchorFolder, resolvedAnchorFolderItem)
			}

			return nil
		})
}

func tryResolveAnchorFolderItems(l *locatorImpl, dirPath string, name string) *models.AnchorFolderItemInfo {
	if isAnchorFolderItem(dirPath) {
		logger.Debugf("Locate anchor folder item. name: %s", name)
		// Avoid unmarshalling instructions since it should be executed upon selection
		anchorFolderItem := newAnchorFolderItem(name, dirPath)
		l.setLongestAnchorFolderItemName(name)
		return anchorFolderItem
	}
	return nil
}

func isExcluded(name string) bool {
	if excluded, exist := excludedDirectories[name]; exist && excluded {
		return true
	}
	return false
}

func isAnchorFolder(dirPath string) bool {
	commandFilePath := fmt.Sprintf("%s/%s", dirPath, globals.AnchorCommandFileName)
	return ioutils.IsValidPath(commandFilePath)
}

func isAnchorFolderItem(dirPath string) bool {
	instructionsFilePath := fmt.Sprintf("%s/%s", dirPath, globals.InstructionsFileName)
	return ioutils.IsValidPath(instructionsFilePath)
}

func hasAnchorIgnoreIdentifier(path string) (string, bool) {
	anchorIgnorePath := fmt.Sprintf("%s/%s", path, anchorIgnoreFileName)
	if _, err := os.Stat(anchorIgnorePath); err != nil {
		return "", false
	}
	return anchorIgnorePath, true
}

func sortAnchorFolders(anchorFoldersInfo []*models.AnchorFolderInfo) []*models.AnchorFolderInfo {
	sort.Slice(anchorFoldersInfo, func(i, j int) bool {
		return anchorFoldersInfo[i].Name < anchorFoldersInfo[j].Name
	})
	return anchorFoldersInfo
}
