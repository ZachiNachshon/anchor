package resolver

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
)

const (
	RemoteUrlFlagName        = "repository.remote.url"
	RemoteBranchFlagName     = "repository.remote.branch"
	RemoteRevisionFlagName   = "repository.remote.revision"
	RemoteClonePathFlagName  = "repository.remote.clonePath"
	RemoteAutoUpdateFlagName = "repository.remote.autoUpdate"
	LocalPathFlagName        = "repository.local.path"
)

type Resolver interface {
	ResolveRepository(ctx common.Context) (string, error)
}

type LocalResolver struct {
	Resolver
	LocalConfig *config.Local
}

type RemoteResolver struct {
	Resolver
	RemoteConfig  *config.Remote
	RemoteActions RemoteResolverActions
}

type RemoteResolverActions interface {
	TryResetToRevision(
		clonePath string,
		branch string,
		revision string) error

	TryFetchRemoteHeadRevision(
		clonePath string,
		repoUrl string,
		branch string) (string, error)

	TryFetchLocalOriginRevision(
		clonePath string,
		branch string) (string, error)

	TryCheckoutToBranch(
		clonePath string,
		branch string) error

	CloneRepositoryIfMissing(
		clonePath string,
		url string,
		branch string) error

	VerifyRemoteRepositoryConfig(remoteCfg *config.Remote) error

	PrintRevisionsDiff(path string, prevRevision string, newRevision string) error
}
