package ioutils

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_IOUtilsShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "parse a valid anchor repository absolute path",
			Func: ParseValidAnchorRepositoryAbsolutePath,
		},
		{
			Name: "avoid infinite loop when parsing repo absolute path",
			Func: AvoidInfiniteLoopWhenParsingRepoAbsolutePath,
		},
		{
			Name: "return that a path is invalid",
			Func: ReturnThatPathIsInvalid,
		},
	}
	harness.RunTests(t, tests)
}

var ParseValidAnchorRepositoryAbsolutePath = func(t *testing.T) {
	pathInTest := "/user/src/github.com/anchor/internal/cmd/anchor/app/status"
	anchorPath := GetRepositoryAbsoluteRootPath(pathInTest)
	assert.Equal(t, repositoryName, filepath.Base(anchorPath))
	assert.NotEmpty(t, anchorPath)
}

var AvoidInfiniteLoopWhenParsingRepoAbsolutePath = func(t *testing.T) {
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
