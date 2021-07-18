package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/git"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

func NewRemoteActions(git git.Git) *remoteActionsImpl {
	return &remoteActionsImpl{
		git: git,
	}
}

type remoteActionsImpl struct {
	RemoteResolverActions
	git git.Git
}

func (ra *remoteActionsImpl) VerifyRemoteRepositoryConfig(remoteCfg *config.Remote) error {
	if remoteCfg == nil {
		return fmt.Errorf("invalid remote repository configuration")
	}
	errFormat := "remote repository config is missing value. name: %s"

	if len(remoteCfg.Url) == 0 {
		return fmt.Errorf(errFormat, "url")
	}

	if len(remoteCfg.Branch) == 0 {
		return fmt.Errorf(errFormat, "branch")
	}

	if len(remoteCfg.ClonePath) == 0 {
		return fmt.Errorf(errFormat, "clonePath")
	}

	return nil
}

func (ra *remoteActionsImpl) CloneRepositoryIfMissing(
	clonePath string,
	url string,
	branch string) error {

	if !ioutils.IsValidPath(clonePath) {
		logger.Infof("Fetching anchorfiles repository for the first time...")
		if err := ra.git.Clone(url, branch, clonePath); err != nil {
			return err
		}
	}
	return nil
}

func (ra *remoteActionsImpl) TryResetToRevision(
	clonePath string,
	branch string,
	revision string) error {

	if err := ra.git.Reset(clonePath, revision); err != nil {
		// TODO: identify a "revision does not exists" error code before fetching again
		if err = ra.git.FetchShallow(clonePath, branch); err != nil {
			return err
		} else {
			if err = ra.git.Reset(clonePath, revision); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ra *remoteActionsImpl) TryFetchRemoteHeadRevision(
	clonePath string,
	repoUrl string,
	branch string) (string, error) {

	if headRevision, err := ra.git.GetRemoteHeadCommitHash(clonePath, repoUrl, branch); err != nil {
		return "", err
	} else {
		return headRevision, nil
	}
}

func (ra *remoteActionsImpl) TryFetchLocalOriginRevision(
	clonePath string,
	branch string) (string, error) {

	if originRevision, err := ra.git.GetLocalOriginCommitHash(clonePath, branch); err != nil {
		return "", err
	} else {
		return originRevision, nil
	}
}

func (ra *remoteActionsImpl) TryCheckoutToBranch(clonePath string, branch string) error {
	return ra.git.Checkout(clonePath, branch)
}

func (ra *remoteActionsImpl) PrintRevisionsDiff(path string, prevRevision string, newRevision string) error {
	return ra.git.LogRevisionsDiffPretty(path, prevRevision, newRevision)
}
