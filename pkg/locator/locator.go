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
	Scan(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError
	CommandFolders() []*models.CommandFolderInfo
	CommandFolderByName(name string) *models.CommandFolderInfo
	CommandFoldersAsMap() map[string]*models.CommandFolderInfo
	CommandFolderItems(commandFolderName string) []*models.CommandFolderItemInfo
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

	initialized                    atomics.AtomicBool
	commandFoldersInfoMap          map[string]*models.CommandFolderInfo
	longestcommandFolderNameLength int
	longestcommandItemNameLength   int
}

func newCommandFolderItem(name string, dirPath string) *models.CommandFolderItemInfo {
	return &models.CommandFolderItemInfo{
		Name:             name,
		DirPath:          dirPath,
		InstructionsPath: fmt.Sprintf("%s/%s", dirPath, globals.InstructionsFileName),
	}
}

func New() Locator {
	return &locatorImpl{
		commandFoldersInfoMap: make(map[string]*models.CommandFolderInfo),
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

	err := scanForCommandFolders(l, anchorfilesLocalPath, e, pa)
	if err != nil {
		return errors.NewLocatorError(err)
	}

	l.markInitialized()
	return nil
}

func (l *locatorImpl) CommandFolders() []*models.CommandFolderInfo {
	res := make([]*models.CommandFolderInfo, 0, len(l.commandFoldersInfoMap))
	for _, v := range l.commandFoldersInfoMap {
		res = append(res, v)
	}
	return sortCommandFolders(res)
}

func (l *locatorImpl) CommandFolderByName(name string) *models.CommandFolderInfo {
	if value, exists := l.commandFoldersInfoMap[name]; exists {
		return value
	}
	return nil
}

func (l *locatorImpl) CommandFoldersAsMap() map[string]*models.CommandFolderInfo {
	return l.commandFoldersInfoMap
}

func (l *locatorImpl) CommandFolderItems(commandFolderName string) []*models.CommandFolderItemInfo {
	if l.commandFoldersInfoMap[commandFolderName] != nil &&
		l.commandFoldersInfoMap[commandFolderName].Items != nil {
		items := l.commandFoldersInfoMap[commandFolderName].Items
		result := make([]*models.CommandFolderItemInfo, len(items))
		i := 0
		for _, item := range items {
			result[i] = item
			i++
		}
		return sortCommandFoldersItems(result)
	}
	return nil
}

func (l *locatorImpl) isInitialized() bool {
	return l.initialized.Get()
}

func (l *locatorImpl) markInitialized() {
	l.initialized.Set(true)
}

func (l *locatorImpl) appendCommandFolder(commandFolder *models.CommandFolderInfo) {
	l.commandFoldersInfoMap[commandFolder.Name] = commandFolder
}

func (l *locatorImpl) appendCommandFolderItem(commandFolder *models.CommandFolderInfo, commandFolderItem *models.CommandFolderItemInfo) {
	if commandFolder.Items == nil {
		commandFolder.Items = make(map[string]*models.CommandFolderItemInfo)
	}
	commandFolder.Items[commandFolderItem.Name] = commandFolderItem
}

// Longest anchor folder name is Used for prompter left padding
func (l *locatorImpl) setLongestcommandFolderName(name string) {
	if l.longestcommandFolderNameLength < len(name) {
		l.longestcommandFolderNameLength = len(name)
	}
}

// Longest anchor folder item name is Used for prompter left padding
func (l *locatorImpl) setLongestcommandItemName(name string) {
	if l.longestcommandItemNameLength < len(name) {
		l.longestcommandItemNameLength = len(name)
	}
}

func scanForCommandFolders(l *locatorImpl, anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) error {
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

			resolvedCommandFolder := tryResolveCommandFolder(l, path, name, e, pa)
			if resolvedCommandFolder != nil {
				scanForCommandFolderItems(l, resolvedCommandFolder)
			}

			return nil
		})

	return nil
}

func tryResolveCommandFolder(l *locatorImpl, dirPath string, name string, e extractor.Extractor, pa parser.Parser) *models.CommandFolderInfo {
	if isCommandFolder(dirPath) {
		logger.Debugf("Locate anchor folder. name: %s", name)
		if commandFolder, err := e.ExtractCommandFolderInfo(dirPath, pa); err != nil {
			return nil
		} else if commandFolder != nil {
			l.appendCommandFolder(commandFolder)
			l.setLongestcommandFolderName(name)
			return commandFolder
		}
	}
	return nil
}

func scanForCommandFolderItems(l *locatorImpl, commandFolder *models.CommandFolderInfo) {
	filepath.Walk(commandFolder.DirPath,
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

			resolvedCommandFolderItem := tryResolveCommandFolderItems(l, path, name)
			if resolvedCommandFolderItem != nil {
				l.appendCommandFolderItem(commandFolder, resolvedCommandFolderItem)
			}

			return nil
		})
}

func tryResolveCommandFolderItems(l *locatorImpl, dirPath string, name string) *models.CommandFolderItemInfo {
	if isCommandFolderItem(dirPath) {
		logger.Debugf("Locate anchor folder item. name: %s", name)
		// Avoid unmarshalling instructions since it should be executed upon selection
		commandFolderItem := newCommandFolderItem(name, dirPath)
		l.setLongestcommandItemName(name)
		return commandFolderItem
	}
	return nil
}

func isExcluded(name string) bool {
	if excluded, exist := excludedDirectories[name]; exist && excluded {
		return true
	}
	return false
}

func isCommandFolder(dirPath string) bool {
	commandFilePath := fmt.Sprintf("%s/%s", dirPath, globals.AnchorCommandFileName)
	return ioutils.IsValidPath(commandFilePath)
}

func isCommandFolderItem(dirPath string) bool {
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

func sortCommandFolders(commandFoldersInfo []*models.CommandFolderInfo) []*models.CommandFolderInfo {
	sort.Slice(commandFoldersInfo, func(i, j int) bool {
		return commandFoldersInfo[i].Name < commandFoldersInfo[j].Name
	})
	return commandFoldersInfo
}

func sortCommandFoldersItems(commandFoldersItemInfo []*models.CommandFolderItemInfo) []*models.CommandFolderItemInfo {
	sort.Slice(commandFoldersItemInfo, func(i, j int) bool {
		return commandFoldersItemInfo[i].Name < commandFoldersItemInfo[j].Name
	})
	return commandFoldersItemInfo
}
