package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/git"
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
		if err := ra.git.Clone(url, branch, clonePath); err != nil {
			return err
		}
	}
	return nil
}

func (ra *remoteActionsImpl) TryResetToRevision(
	clonePath string,
	url string,
	branch string,
	revision string) error {

	if err := ra.git.Reset(clonePath, revision); err != nil {
		// TODO: identify a "revision does not exists" error code before fetching again
		if err = ra.git.FetchShallow(clonePath, url, branch); err != nil {
			return err
		} else {
			if err = ra.git.Reset(clonePath, revision); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ra *remoteActionsImpl) TryFetchHeadRevision(
	clonePath string,
	url string,
	branch string) error {

	if headRevision, err := ra.git.GetHeadCommitHash(branch); err != nil {
		return err
	} else {
		if err = ra.TryResetToRevision(
			clonePath,
			url,
			branch,
			headRevision); err != nil {
			return err
		}
	}
	return nil
}
