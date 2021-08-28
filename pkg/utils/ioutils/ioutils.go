package ioutils

import (
	"fmt"
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
			fmt.Printf("found: %v, dirPath: %s\n", found, dirPath)
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

//func CreateDirectory(path string) error {
//	if _, err := os.Stat(path); os.IsNotExist(err) {
//		return os.Mkdir(path, os.ModePerm)
//	}
//	return nil
//}

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
		if file, err := os.Create(path); err != nil {
			return nil, err
		} else {
			return file, nil
		}
	} else {
		if file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
			return nil, err
		} else {
			return file, err
		}
	}
}
