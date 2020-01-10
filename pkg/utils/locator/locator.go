package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"strconv"

	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

var DirLocator Locator

type DirectoryIdentifier string

var excludedDirectories map[string]bool

var lineFormat = "| %-3v | %-35v %-15v %-15v %-15v\n"
var header = fmt.Sprintf(lineFormat, "#", "NAME", "DOCKERFILE", "K8S_MANIFEST", "AFFINITY")

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
	DockerFileIdentifier   DirectoryIdentifier = "Dockerfile"
	ManifestsIdentifier    DirectoryIdentifier = "k8s/manifest.yaml"
	AnchorIgnoreIdentifier DirectoryIdentifier = ".anchorignore"
)

type locator struct {
	dirs        map[string]*dirContent
	dirsNumeric map[int]*dirContent
}

type dirContent struct {
	name              string
	k8sManifest       string
	dockerfile        string
	affinity          string
	dockerContextRoot string
}

type ListOpts struct {
	OnlyK8sManifests bool
	AffinityFilter   string
}

func New() Locator {
	var locator = &locator{
		dirs:        make(map[string]*dirContent),
		dirsNumeric: make(map[int]*dirContent),
	}

	return locator
}

func (l *locator) Scan() error {
	i := 1
	err := filepath.Walk(common.GlobalOptions.DockerRepositoryPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignore root directory
			if info.IsDir() && path == common.GlobalOptions.DockerRepositoryPath {
				return nil
			}

			// Ignore all hidden directories
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			// Ignore all files under root folder
			dir := filepath.Dir(path)
			if !info.IsDir() && dir == common.GlobalOptions.DockerRepositoryPath {
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

			if dockerfilePath, ok := hasDockerfile(path); ok {

				dirContent := new(dirContent)
				dirContent.name = name
				dirContent.dockerfile = dockerfilePath

				if _, ok := hasAnchorIgnoreIdentifier(path); ok {
					return filepath.SkipDir
				}

				if k8sManifestPath, ok := hasKubernetesManifest(path); ok {
					dirContent.k8sManifest = k8sManifestPath
				}

				if affinity, ok := hasAffinity(path); ok {
					dirContent.affinity = affinity
				}

				ctxRoot := filepath.Dir(dockerfilePath)
				dirContent.dockerContextRoot = ctxRoot

				l.dirs[name] = dirContent

				// Maintain a 2nd map based on numeric keys for easier CLI selection
				l.dirsNumeric[i] = dirContent
				i += 1
			}
			return nil
		})

	if err != nil {
		return err
	}

	return nil
}

func (l *locator) Print(opts *ListOpts) {
	size := len(l.dirsNumeric)
	if size == 0 {
		return
	}

	var affinityFilter = ""
	var listK8sOnly = false
	if opts != nil {
		listK8sOnly = opts.OnlyK8sManifests
		affinityFilter = opts.AffinityFilter
	}

	table := "\n"
	table += header
	for lineNumber := 1; lineNumber <= size; {
		content := l.dirsNumeric[lineNumber]

		if len(affinityFilter) > 0 && content.affinity != affinityFilter {
			lineNumber += 1
			continue
		}

		hasDockerfile := "     no"
		if content.dockerfile != "" {
			hasDockerfile = "   yes"
		}

		hasK8sManifest := "    no"
		hasK8s := false
		if content.k8sManifest != "" {
			hasK8sManifest = "    yes"
			hasK8s = true
		}

		if listK8sOnly && !hasK8s {
			lineNumber += 1
			continue
		} else {
			l := fmt.Sprintf(lineFormat, lineNumber, content.name, hasDockerfile, hasK8sManifest, content.affinity)
			table += l
			lineNumber += 1
		}
	}

	logger.Info(table)
}

func (l *locator) Name(identifier string) (string, error) {
	if number, err := strconv.Atoi(identifier); err == nil {
		if content, ok := l.dirsNumeric[number]; ok {
			return content.name, nil
		}
	} else {
		if content, ok := l.dirs[identifier]; ok {
			return content.name, nil
		}
	}
	return "", errors.Errorf("Cannot find Dockerfile for %v", identifier)
}

func (l *locator) Dockerfile(identifier string) (string, error) {
	if number, err := strconv.Atoi(identifier); err == nil {
		if content, ok := l.dirsNumeric[number]; ok {
			return content.dockerfile, nil
		}
	} else {
		if content, ok := l.dirs[identifier]; ok {
			return content.dockerfile, nil
		}
	}
	return "", errors.Errorf("Cannot find Dockerfile for %v", identifier)
}

func (l *locator) DockerContext(identifier string) (string, error) {
	if number, err := strconv.Atoi(identifier); err == nil {
		if content, ok := l.dirsNumeric[number]; ok {
			return content.dockerContextRoot, nil
		}
	} else {
		if content, ok := l.dirs[identifier]; ok {
			return content.dockerContextRoot, nil
		}
	}
	return "", errors.Errorf("Cannot find Docker context directory for %v", identifier)
}

func (l *locator) Manifest(identifier string) (string, error) {
	if number, err := strconv.Atoi(identifier); err == nil {
		if content, ok := l.dirsNumeric[number]; ok {
			return content.k8sManifest, nil
		}
	} else {
		if content, ok := l.dirs[identifier]; ok {
			return content.k8sManifest, nil
		}
	}
	return "", errors.Errorf("Cannot find K8s manifest for %v", identifier)
}

func isExcluded(name string) bool {
	if excluded, exist := excludedDirectories[name]; exist && excluded {
		return true
	}
	return false
}

func hasDockerfile(path string) (string, bool) {
	dockerfilePath := fmt.Sprintf("%v/%v", path, DockerFileIdentifier)
	if _, err := os.Stat(dockerfilePath); err != nil {
		return "", false
	}
	return dockerfilePath, true
}

func hasKubernetesManifest(path string) (string, bool) {
	k8sManifestPath := fmt.Sprintf("%v/%v", path, ManifestsIdentifier)
	if _, err := os.Stat(k8sManifestPath); err != nil {
		return "", false
	}
	return k8sManifestPath, true
}

func hasAnchorIgnoreIdentifier(path string) (string, bool) {
	anchorIgnorePath := fmt.Sprintf("%v/%v", path, AnchorIgnoreIdentifier)
	if _, err := os.Stat(anchorIgnorePath); err != nil {
		return "", false
	}
	return anchorIgnorePath, true
}

func hasAffinity(path string) (string, bool) {
	parentPath := path[:strings.LastIndex(path, "/")]
	rootDirName := filepath.Base(common.GlobalOptions.DockerRepositoryPath)
	parentDirName := filepath.Base(parentPath)

	if rootDirName == parentDirName {
		return "", false
	}
	return parentDirName, true
}
