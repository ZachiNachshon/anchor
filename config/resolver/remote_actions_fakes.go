package resolver

import (
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/git"
)

var CreateFakeRemoteActions = func() *fakeRemoteActionsImpl {
	return &fakeRemoteActionsImpl{
		git: git.CreateFakeGit(),
	}
}

type fakeRemoteActionsImpl struct {
	RemoteResolverActions
	git git.Git

	TryResetToRevisionMock func(
		clonePath string,
		url string,
		branch string,
		revision string) error

	TryFetchHeadRevisionMock func(
		clonePath string,
		url string,
		branch string) error

	CloneRepositoryIfMissingMock func(
		clonePath string,
		url string,
		branch string) error

	VerifyRemoteRepositoryConfigMock func(remoteCfg *config.Remote) error
}

func (ra *fakeRemoteActionsImpl) VerifyRemoteRepositoryConfig(remoteCfg *config.Remote) error {
	return ra.VerifyRemoteRepositoryConfigMock(remoteCfg)
}

func (ra *fakeRemoteActionsImpl) TryResetToRevision(
	clonePath string,
	url string,
	branch string,
	revision string) error {

	return ra.TryResetToRevisionMock(clonePath, url, branch, revision)
}

func (ra *fakeRemoteActionsImpl) TryFetchHeadRevision(
	clonePath string,
	url string,
	branch string) error {

	return ra.TryFetchHeadRevisionMock(clonePath, url, branch)
}

func (ra *fakeRemoteActionsImpl) CloneRepositoryIfMissing(
	clonePath string,
	url string,
	branch string) error {

	return ra.CloneRepositoryIfMissingMock(clonePath, url, branch)
}
