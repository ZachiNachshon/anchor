package repository

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/internal/repository/local"
	"github.com/ZachiNachshon/anchor/internal/repository/remote"
	"github.com/ZachiNachshon/anchor/pkg/git"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

type Repository interface {
	Load(ctx common.Context) (string, error)
}

var GetRepositoryOriginByConfig = func(repoConfig *config.Repository) (Repository, error) {
	// Checks if repository config attribute is empty
	if repoConfig == nil {
		return nil, fmt.Errorf("missing required config value. name: repository")
	}

	if repoConfig.Local != nil && len(repoConfig.Local.Path) > 0 {
		logger.Debug("Using local anchorfiles repository")
		return &local.LocalRepository{
			LocalConfig: repoConfig.Local,
		}, nil

	} else if repoConfig.Remote != nil {
		logger.Debugf("Using remote anchorfiles repository")
		g := git.New(shell.New())
		return &remote.RemoteRepository{
			RemoteConfig:  repoConfig.Remote,
			RemoteActions: remote.NewRemoteActions(g),
		}, nil
	}
	return nil, fmt.Errorf("could not resolve anchorfiles local repository path or git tracked remote repository")
}
