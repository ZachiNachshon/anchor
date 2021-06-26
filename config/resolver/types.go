package resolver

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/git"
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
	RemoteConfig *config.Remote
	Git          git.Git
}
