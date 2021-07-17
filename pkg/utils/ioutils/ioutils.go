package ioutils

import (
	"os"
	"path/filepath"
)

const repositoryName = "anchor"

func GetRepositoryAbsoluteRootPath() string {
	path, _ := os.Getwd()

	trailingPath := filepath.Base(path)
	if trailingPath == repositoryName {
		return path
	}

	for search := true; search; search = trailingPath != repositoryName {
		path = filepath.Dir(path)
		trailingPath = filepath.Base(path)
		search = false
	}

	return path
}

func IsValidPath(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func GetUserHomeFolder() (string, error) {
	if dirname, err := os.UserHomeDir(); err != nil {
		return "", err
	} else {
		return dirname, nil
	}
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
