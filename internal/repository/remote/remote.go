package remote

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
)

type RemoteRepository struct {
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

func (rr *RemoteRepository) Load(ctx common.Context) (string, error) {
	if rr.RemoteActions == nil {
		return "", fmt.Errorf("remote actions weren't defined for remote resolver, cannot proceed")
	}

	if err := rr.RemoteActions.VerifyRemoteRepositoryConfig(rr.RemoteConfig); err != nil {
		return "", err
	}

	clonePath := rr.RemoteConfig.ClonePath
	url := rr.RemoteConfig.Url
	branch := rr.RemoteConfig.Branch
	if err := rr.RemoteActions.CloneRepositoryIfMissing(clonePath, url, branch); err != nil {
		return "", err
	}

	if len(rr.RemoteConfig.Revision) > 0 {
		if err := rr.RemoteActions.TryResetToRevision(
			clonePath,
			branch,
			rr.RemoteConfig.Revision); err != nil {
			return "", err
		}
		logger.Infof("Updated anchorfiles repo to revision. hash: %s", rr.RemoteConfig.Revision)

		if rr.RemoteConfig.AutoUpdate {
			msg := fmt.Sprintf("Mutually exclusive config values found: autoUpdate / revision. "+
				"To allow auto update from '%s' branch latest HEAD, remove the revision from config.",
				branch)

			logger.Warning(msg)
		}

	} else if rr.RemoteConfig.AutoUpdate {
		logger.Debug("Checking anchorfiles local origin revision...")
		originRevision, err := rr.RemoteActions.TryFetchLocalOriginRevision(clonePath, branch)
		if err != nil {
			return "", err
		}

		logger.Info("Checking anchorfiles remote HEAD revision...")
		headRevision, err := rr.RemoteActions.TryFetchRemoteHeadRevision(clonePath, rr.RemoteConfig.Url, branch)
		if err != nil {
			return "", err
		}

		if err = rr.RemoteActions.TryResetToRevision(
			clonePath,
			branch,
			headRevision); err != nil {
			return "", err
		}

		if originRevision != headRevision {
			logger.Infof("Fetched remote HEAD revision. hash: %s", headRevision)
			err = rr.RemoteActions.PrintRevisionsDiff(clonePath, originRevision, headRevision)
			if err != nil {
				logger.Debugf("failed to print revisions diff. error: %s", err.Error())
				// Do not return an error if print fails
				//return "", err
			}
		} else {
			logger.Info("Already up to date !")
		}
	}

	err := rr.RemoteActions.TryCheckoutToBranch(clonePath, branch)
	if err != nil {
		return "", err
	}

	return clonePath, nil
}
