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
	}
	harness.RunTests(t, tests)
}

var ParsePathWithAnchorSubDirectory = func(t *testing.T) {
	workingDir, _ := os.Getwd()
	pathInTest := workingDir + "/sub-dir/anchor"
	// <REPO_PATH>/anchor/sub-dir/anchor...
	anchorPath := GetRepositoryAbsoluteRootPath(pathInTest)
	assert.NotEmpty(t, anchorPath)
	assert.NotContains(t, anchorPath, "/sub-dir/anchor", "failed parsing path: %s", pathInTest)
	assert.Equal(t, repositoryName, filepath.Base(anchorPath), "failed parsing path: %s", pathInTest)
}

var ParsePathWithConsecutiveAnchorDirNames = func(t *testing.T) {
	workingDir, _ := os.Getwd()
	pathInTest := workingDir + "/anchor"
	// <REPO_PATH>/anchor/anchor/path...
	anchorPath := GetRepositoryAbsoluteRootPath(pathInTest)
	assert.NotEmpty(t, anchorPath)
	assert.NotContains(t, anchorPath, "/anchor/anchor", "failed parsing path: %s", pathInTest)
	assert.Equal(t, repositoryName, filepath.Base(anchorPath), "failed parsing path: %s", pathInTest)
}

var StopAtRootFolderWhenParsingRepoAbsolutePath = func(t *testing.T) {
	pathInTest := "/user/src/github.com/noname/some/example/path"
	anchorPath := GetRepositoryAbsoluteRootPath(pathInTest)
	assert.Equal(t, "/", anchorPath)
	assert.NotEmpty(t, anchorPath)
}

var ReturnThatPathIsInvalid = func(t *testing.T) {
	pathInTest := "/invalid/path"
	isValid := IsValidPath(pathInTest)
	assert.False(t, isValid)
}
