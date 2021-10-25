package ioutils

import (
	"os"
	"path/filepath"
)

const repositoryName = "anchor"

func GetRepositoryAbsoluteRootPath(path string) string {
	// Parsing absolute repo root path should:
	//  - Avoid internal use of anchor as package name
	//    $GOPATH/src/github.com/anchor/internal/cmd/anchor/app/status
	//  - Handle consecutive repo name
	//    $GITHUB_WORKSPACE/anchor/anchor/.git
	dirPath := path
	dirName := filepath.Base(dirPath)
	for found := false; !found && dirPath != "/"; {
		if dirName == repositoryName {
			found = IsValidPath(dirPath+"/go.mod") || IsValidPath(dirPath+"/.git")
		}
		if found {
			break
		}
		dirPath = filepath.Dir(dirPath)
		dirName = filepath.Base(dirPath)
	}

	return dirPath
}

func IsValidPath(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func GetUserHomeDirectory() (string, error) {
	if dirname, err := os.UserHomeDir(); err != nil {
		return "", err
	} else {
		return dirname, nil
	}
}

func GetWorkingDirectory() (string, error) {
	return os.Getwd()
}

func CreateOrOpenFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return createFile(path)
	} else {
		if file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
			return nil, err
		} else {
			return file, err
		}
	}
}

func CreateOrOpenFileWithModes(path string, modeFlags int) (*os.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return createFile(path)
	} else {
		if file, err := os.OpenFile(path, modeFlags, 0666); err != nil {
			return nil, err
		} else {
			return file, err
		}
	}
}

func createFile(path string) (*os.File, error) {
	dir, _ := filepath.Split(path)
	if err := createDirectory(dir); err != nil {
		return nil, err
	}
	if file, err := os.Create(path); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func createDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}
