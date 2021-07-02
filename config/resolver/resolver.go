package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/git"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func GetResolverBasedOnConfig(repoConfig *config.Repository) (Resolver, error) {
	// Checks if repository config attribute is empty
	if repoConfig == nil {
		return nil, fmt.Errorf("missing required config value. name: repository")
	}

	if repoConfig.Local != nil && len(repoConfig.Local.Path) > 0 {
		logger.Debug("Using local anchorfiles repository")
		return &LocalResolver{
			LocalConfig: repoConfig.Local,
		}, nil
	} else if repoConfig.Remote != nil {
		logger.Debugf("Using remote anchorfiles repository")
		g := git.New(shell.New())
		return &RemoteResolver{
			RemoteConfig:  repoConfig.Remote,
			RemoteActions: NewRemoteActions(g),
		}, nil
	}
	return nil, fmt.Errorf("could not resolve anchorfiles local repository path or git tracked remote repository")
}
