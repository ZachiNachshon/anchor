package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/atomics"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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
	appDirs              map[string]*models.AppContent
	longestAppNameLength int
}

func newAppContent(name string, path string) *models.AppContent {
	return &models.AppContent{
		Name:             name,
		DirPath:          path,
		InstructionsPath: fmt.Sprintf("%s/%s", path, instructionsFileName),
	}
}

func New() Locator {
	return &locatorImpl{
		appDirs: make(map[string]*models.AppContent),
	}
}

func (l *locatorImpl) Scan(anchorfilesLocalPath string) error {
	if l.isInitialized() {
		logger.Warning("scan can be called only once, using previous scan result")
		return nil
	}

	if !ioutils.IsValidPath(anchorfilesLocalPath) {
		return fmt.Errorf("invalid anchorfile local path. path: %s", anchorfilesLocalPath)
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
		return err
	}

	l.markInitialized()
	return nil
}

func (l *locatorImpl) Applications() []*models.AppContent {
	res := make([]*models.AppContent, 0, len(l.appDirs))
	for _, v := range l.appDirs {
		res = append(res, v)
	}
	return sortApplications(res)
}

func (l *locatorImpl) ApplicationsAsMap() map[string]*models.AppContent {
	return l.appDirs
}

func (l *locatorImpl) Application(name string) *models.AppContent {
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
		logger.Debugf("locate application. Name: %s", name)
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

func sortApplications(apps []*models.AppContent) []*models.AppContent {
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Name < apps[j].Name
	})
	return apps
}

//func (l *locator) Print() {
//size := len(l.dirsNumeric)
//if size == 0 {
//	return
//}
//
//var affinityFilter = ""
//var listK8sOnly = false
//
//table := "\n"
//table += header
//for lineNumber := 1; lineNumber <= size; {
//	content := l.dirsNumeric[lineNumber]
//
//	if len(affinityFilter) > 0 && content.affinity != affinityFilter {
//		lineNumber += 1
//		continue
//	}
//
//	hasDockerfile := "     no"
//	if content.dockerfile != "" {
//		hasDockerfile = "   yes"
//	}
//
//	hasK8sManifest := "    no"
//	hasK8s := false
//	if content.k8sManifest != "" {
//		hasK8sManifest = "    yes"
//		hasK8s = true
//	}
//
//	if listK8sOnly && !hasK8s {
//		lineNumber += 1
//		continue
//	} else {
//		l := fmt.Sprintf(lineFormat, lineNumber, content.Name, hasDockerfile, hasK8sManifest, content.affinity)
//		table += l
//		lineNumber += 1
//	}
//}
//
//logger.Info(table)
//}
