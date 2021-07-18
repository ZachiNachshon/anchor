package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
)

func (rr *RemoteResolver) ResolveRepository(ctx common.Context) (string, error) {
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

		anchorConfig := config.FromContext(ctx)
		repoUrl := anchorConfig.Config.Repository.Remote.Url
		logger.Info("Checking anchorfiles remote HEAD revision...")
		headRevision, err := rr.RemoteActions.TryFetchRemoteHeadRevision(clonePath, repoUrl, branch)
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
			logger.Info("Already up to date")
		}
	}

	err := rr.RemoteActions.TryCheckoutToBranch(clonePath, branch)
	if err != nil {
		return "", err
	}

	return clonePath, nil
}
