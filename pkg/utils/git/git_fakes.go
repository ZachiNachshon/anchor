package git

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

var CreateFakeGit = func() *fakeGitImpl {
	return &fakeGitImpl{
	}
}

type fakeGitImpl struct {
	Git
	CloneMock        func(s shell.Shell, clonePath string) error
	InitMock         func(s shell.Shell, path string) error
	AddOriginMock    func(s shell.Shell, path string, url string) error
	FetchShallowMock func(s shell.Shell, path string, url string, branch string) error
	ResetMock        func(s shell.Shell, path string, revision string) error
	CleanMock        func(s shell.Shell, path string) error
}

func (g *fakeGitImpl) Clone(s shell.Shell, clonePath string) error {
	return g.CloneMock(s, clonePath)
}

func (g *fakeGitImpl) Init(s shell.Shell, path string) error {
	return g.InitMock(s, path)
}

func (g *fakeGitImpl) AddOrigin(s shell.Shell, path string, url string) error {
	return g.AddOriginMock(s, path, url)
}

func (g *fakeGitImpl) FetchShallow(s shell.Shell, path string, url string, branch string) error {
	return g.FetchShallowMock(s, path, url, branch)
}

func (g *fakeGitImpl) Reset(s shell.Shell, path string, revision string) error {
	return g.ResetMock(s, path, revision)
}

func (g *fakeGitImpl) Clean(s shell.Shell, path string) error {
	return g.CleanMock(s, path)
}
