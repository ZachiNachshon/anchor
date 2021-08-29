package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/models"

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
	Applications() []*models.ApplicationInfo
	ApplicationsAsMap() map[string]*models.ApplicationInfo
	Application(name string) *models.ApplicationInfo
}

type DirectoryIdentifier string

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
	app                  DirectoryIdentifier = "app"
	cli                  DirectoryIdentifier = "cli"
	controller           DirectoryIdentifier = "controller"
	docker               DirectoryIdentifier = "docker"
	k8s                  DirectoryIdentifier = "k8s"
	instructionsFileName string              = "instructions.yaml"
	anchorIgnoreFileName string              = ".anchorignore"
)

type locatorImpl struct {
	Locator
	initialized          atomics.AtomicBool
	anchorfilesLocalPath string
	appDirs              map[string]*models.ApplicationInfo
	longestAppNameLength int
}

func newAppContent(name string, path string) *models.ApplicationInfo {
	return &models.ApplicationInfo{
		Name:             name,
		DirPath:          path,
		InstructionsPath: fmt.Sprintf("%s/%s", path, instructionsFileName),
	}
}

func New() *locatorImpl {
	return &locatorImpl{
		appDirs: make(map[string]*models.ApplicationInfo),
	}
}

func (l *locatorImpl) Scan(anchorfilesLocalPath string) *errors.LocatorError {
	if l.isInitialized() {
		warnMsg := "scan can be called only once, using previous scan result"
		logger.Warning(warnMsg)
		return errors.NewAlreadyInitializedError(fmt.Errorf(warnMsg))
	}

	if !ioutils.IsValidPath(anchorfilesLocalPath) {
		return errors.NewLocatorError(fmt.Errorf("invalid anchorfile local path. path: %s", anchorfilesLocalPath))
	}

	err := filepath.Walk(anchorfilesLocalPath,
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

			if tryResolveApplication(l, path, name) {
				return nil
			}

			return nil
		})

	if err != nil {
		return errors.NewLocatorError(err)
	}

	l.markInitialized()
	return nil
}

func (l *locatorImpl) Applications() []*models.ApplicationInfo {
	res := make([]*models.ApplicationInfo, 0, len(l.appDirs))
	for _, v := range l.appDirs {
		res = append(res, v)
	}
	return sortApplications(res)
}

func (l *locatorImpl) ApplicationsAsMap() map[string]*models.ApplicationInfo {
	return l.appDirs
}

func (l *locatorImpl) Application(name string) *models.ApplicationInfo {
	if value, exists := l.appDirs[name]; exists {
		return value
	}
	return nil
}

func (l *locatorImpl) isInitialized() bool {
	return l.initialized.Get()
}

func (l *locatorImpl) markInitialized() {
	l.initialized.Set(true)
}

// Longest application name is Used for prompter left padding
func setLongestApplicationName(l *locatorImpl, name string) {
	if l.longestAppNameLength < len(name) {
		l.longestAppNameLength = len(name)
	}
}

func tryResolveApplication(l *locatorImpl, path string, name string) bool {
	if isApp := isApplication(path); isApp {
		logger.Debugf("Locate application. Name: %s", name)
		appContent := newAppContent(name, path)
		l.appDirs[name] = appContent
		setLongestApplicationName(l, name)
		return true
	}
	return false
}

func isExcluded(name string) bool {
	if excluded, exist := excludedDirectories[name]; exist && excluded {
		return true
	}
	return false
}

func isApplication(path string) bool {
	dirPath := filepath.Dir(path)
	return filepath.Base(dirPath) == string(app)
}

func hasAnchorIgnoreIdentifier(path string) (string, bool) {
	anchorIgnorePath := fmt.Sprintf("%s/%s", path, anchorIgnoreFileName)
	if _, err := os.Stat(anchorIgnorePath); err != nil {
		return "", false
	}
	return anchorIgnorePath, true
}

func sortApplications(apps []*models.ApplicationInfo) []*models.ApplicationInfo {
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Name < apps[j].Name
	})
	return apps
}
