package resolver

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
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
		url string,
		branch string,
		revision string) error

	TryFetchHeadRevision(
		clonePath string,
		url string,
		branch string) error

	CloneRepositoryIfMissing(clonePath string) error
	VerifyRemoteRepositoryConfig(remoteCfg *config.Remote) error
}
