package ioutils

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_IOUtilsShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "parse path with anchor sub directory",
			Func: ParsePathWithAnchorSubDirectory,
		},
		{
			Name: "parse path with consecutive anchor dir names",
			Func: ParsePathWithConsecutiveAnchorDirNames,
		},
		{
			Name: "stop at root folder when parsing repo absolute path",
			Func: StopAtRootFolderWhenParsingRepoAbsolutePath,
		},
		{
			Name: "return that a path is invalid",
			Func: ReturnThatPathIsInvalid,
		},
		{
			Name: "return that a path is valid",
			Func: ReturnThatPathIsValid,
		},
		{
			Name: "return user home directory",
			Func: ReturnUserHomeDirectory,
		},
		{
			Name: "return working directory",
			Func: ReturnWorkingDirectory,
		},
		{
			Name: "create a new file if not exists",
			Func: CreateNewFileIfNotExists,
		},
		{
			Name: "open file if exists",
			Func: OpenFileIfExists,
		},
		{
			Name: "create new file with folder hierarchy if not exists",
			Func: CreateNewFileWithFolderHierarchyIfNotExists,
		},
		{
			Name: "create new file with modes if not exists",
			Func: CreateNewFileWithModesIfNotExists,
		},
		{
			Name: "open file with modes if exists",
			Func: OpenFileWithModesIfExists,
		},
	}
	harness.RunTests(t, tests)
}

var ParsePathWithAnchorSubDirectory = func(t *testing.T) {
	workingDir, _ := os.Getwd()
	pathInTest := workingDir + "anchor/sub-dir/anchor"
	// <REPO_PATH>/anchor/sub-dir/anchor...
	anchorPath := GetRepositoryAbsoluteRootPath(pathInTest)
	assert.NotEmpty(t, anchorPath)
	assert.NotContains(t, anchorPath, "/sub-dir/anchor", "failed parsing path: %s", pathInTest)
	assert.Equal(t, repositoryName, filepath.Base(anchorPath), "failed parsing path: %s", pathInTest)
}

var ParsePathWithConsecutiveAnchorDirNames = func(t *testing.T) {
	workingDir, _ := os.Getwd()
	pathInTest := workingDir + "/anchor/anchor"
	// <REPO_PATH>/anchor/anchor/path...
	anchorPath := GetRepositoryAbsoluteRootPath(pathInTest)
	assert.NotEmpty(t, anchorPath)
	//assert.NotContains(t, anchorPath, "/anchor/anchor", "failed parsing path: %s", pathInTest)
	assert.Equal(t, repositoryName, filepath.Base(anchorPath), "failed parsing path: %s", pathInTest)
}

var StopAtRootFolderWhenParsingRepoAbsolutePath = func(t *testing.T) {
	pathInTest := "/user/src/github.com/noname/some/example/path"
	anchorPath := GetRepositoryAbsoluteRootPath(pathInTest)
	assert.Equal(t, "/", anchorPath)
	assert.NotEmpty(t, anchorPath)
}

var ReturnThatPathIsInvalid = func(t *testing.T) {
	invalidPath := "/invalid/path"
	isValid := IsValidPath(invalidPath)
	assert.False(t, isValid)
}

var ReturnThatPathIsValid = func(t *testing.T) {
	validPath, _ := os.Getwd()
	isValid := IsValidPath(validPath)
	assert.True(t, isValid)
}

var ReturnUserHomeDirectory = func(t *testing.T) {
	homeDir, err := GetUserHomeDirectory()
	assert.Nil(t, err, "expected to succeed")
	assert.NotEmpty(t, homeDir)
}

var ReturnWorkingDirectory = func(t *testing.T) {
	wd, err := GetWorkingDirectory()
	assert.Nil(t, err, "expected to succeed")
	assert.NotEmpty(t, wd)
}

var CreateNewFileIfNotExists = func(t *testing.T) {
	tempDir := os.TempDir()
	tempConfigFile := tempDir + "/tempFile.txt"
	f, err := CreateOrOpenFile(tempConfigFile)
	assert.Nil(t, err, "expected to succeed")
	assert.NotNil(t, f)
}

var OpenFileIfExists = func(t *testing.T) {
	tempDir := os.TempDir()
	tempConfigFile := tempDir + "/tempFile.txt"

	err := os.WriteFile(tempConfigFile, []byte("some test text"), 0)
	assert.Nil(t, err, "expected write updated config to a temp file successfully")

	f, err := CreateOrOpenFile(tempConfigFile)
	assert.Nil(t, err, "expected to succeed")
	assert.NotNil(t, f)
}

var CreateNewFileWithFolderHierarchyIfNotExists = func(t *testing.T) {
	tempDir := os.TempDir()
	tempConfigFile := tempDir + "new/folder/tempFile.txt"
	f, err := createFile(tempConfigFile)
	assert.Nil(t, err, "expected to succeed")
	assert.NotNil(t, f)
}

var CreateNewFileWithModesIfNotExists = func(t *testing.T) {
	tempDir := os.TempDir()
	tempConfigFile := tempDir + "/tempFile.txt"
	f, err := CreateOrOpenFileWithModes(tempConfigFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND)
	assert.Nil(t, err, "expected to succeed")
	assert.NotNil(t, f)
}

var OpenFileWithModesIfExists = func(t *testing.T) {
	tempDir := os.TempDir()
	tempConfigFile := tempDir + "/tempFile.txt"

	err := os.WriteFile(tempConfigFile, []byte("some test text"), 0)
	assert.Nil(t, err, "expected write updated config to a temp file successfully")

	f, err := CreateOrOpenFileWithModes(tempConfigFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND)
	assert.Nil(t, err, "expected to succeed")
	assert.NotNil(t, f)
}
