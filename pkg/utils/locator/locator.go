package locator

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryIdentifier string

const (
	DOCKER_FILE_IDENTIFIER DirectoryIdentifier = "Dockerfile"
	MANIFESTS_IDENTIFIER   DirectoryIdentifier = "k8s/manifest.yaml"
)

type locator struct {
}

func NewLocator() Locator {
	return &locator{}
}

func (l *locator) Dockerfile(name string) (string, error) {
	expected := fmt.Sprintf("%v/%v/%v", common.GlobalOptions.DockerRepositoryPath, name, DOCKER_FILE_IDENTIFIER)
	dirNames, _ := GetDirNamesNoPath(false, DOCKER_FILE_IDENTIFIER)

	for _, e := range dirNames {
		if strings.EqualFold(expected, e) {
			return e, nil
		}
	}

	return "", errors.Errorf("Cannot find Dockerfile for %v", name)
}

func (l *locator) DockerfileDir(name string) (string, error) {
	if path, err := l.Dockerfile(name); err != nil {
		return "", err
	} else {
		return filepath.Dir(path), nil
	}
}

func (l *locator) Manifest(name string) (string, error) {
	expected := fmt.Sprintf("%v/%v/%v", common.GlobalOptions.DockerRepositoryPath, name, MANIFESTS_IDENTIFIER)
	dirNames, _ := GetDirNamesNoPath(false, MANIFESTS_IDENTIFIER)

	for _, e := range dirNames {
		if strings.EqualFold(expected, e) {
			return e, nil
		}
	}

	return "", errors.Errorf("Cannot find K8s manifest for %v", name)
}

func (l *locator) ManifestDir(name string) (string, error) {
	if path, err := l.Manifest(name); err != nil {
		return "", err
	} else {
		return filepath.Dir(path), nil
	}
}

func (l *locator) GetRootFromManifestFile(path string) string {
	// Dirname to k8s
	dirName := filepath.Dir(path)
	// Dirname to container dir name
	dirName = filepath.Dir(dirName)
	return dirName
}

func (l *locator) GetRootFromDockerfile(path string) string {
	// Dirname to container dir name
	dirName := filepath.Dir(path)
	return dirName
}

func GetDirNamesNoPath(verbose bool, identifier DirectoryIdentifier) ([]string, error) {
	var dirNames = make([]string, 0)
	err := filepath.Walk(common.GlobalOptions.DockerRepositoryPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Continue to the next path
			if !strings.Contains(path, string(identifier)) {
				return nil
			}

			if filePath, err := filepath.Abs(path); err != nil {
				return err
			} else {
				dirName := extractDirName(filePath, identifier)

				if verbose {
					logger.Info("  " + dirName)
				}

				dirNames = append(dirNames, filePath)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return dirNames, nil
}

func extractDirName(path string, identifier DirectoryIdentifier) string {
	dirName := strings.TrimPrefix(path, common.GlobalOptions.DockerRepositoryPath+"/")
	dirName = strings.TrimSuffix(dirName, "/"+string(identifier))
	return dirName
}
