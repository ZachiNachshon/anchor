package harness

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"testing"
)

type TestsHarness struct {
	Name string
	Func func(t *testing.T)
}

func RunTests(t *testing.T, tests []TestsHarness) {
	for _, tt := range tests {
		t.Run(tt.Name, tt.Func)
	}
}

func HarnessAnchorfilesTestRepo(ctx common.Context) {
	repoRootPath := ioutils.GetRepositoryAbsoluteRootPath()
	if repoRootPath == "" {
		logger.Fatalf("failed to resolve the absolute path of the repository root.")
	}
	anchorfilesPathTest := fmt.Sprintf("%s/test/data/anchorfiles", repoRootPath)
	ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(anchorfilesPathTest)
}

func HarnessAnchorfilesRemoteGitTestRepo(ctx common.Context) {
	repoRootPath := ioutils.GetRepositoryAbsoluteRootPath()
	if repoRootPath == "" {
		logger.Fatalf("failed to resolve the absolute path of the repository root.")
	}
	anchorfilesPathTest := fmt.Sprintf("%s/test/data/anchorfiles-git-based", repoRootPath)
	ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(anchorfilesPathTest)
}
