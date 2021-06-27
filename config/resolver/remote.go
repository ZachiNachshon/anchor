package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
)

func (rr *RemoteResolver) ResolveRepository(ctx common.Context) (string, error) {
	if err := rr.verifyRemoteRepositoryConfig(); err != nil {
		return "", err
	}

	clonePath := rr.RemoteConfig.ClonePath
	if !rr.IsClonedPathExists(clonePath) {
		if err := rr.Git.Clone(clonePath); err != nil {
			return "", err
		}
	}

	if len(rr.RemoteConfig.Revision) > 0 {
		if err := rr.Git.Reset(clonePath, rr.RemoteConfig.Revision); err != nil {
			// TODO: identify a "revision does not exists" error code before fetching again
			if err = rr.Git.FetchShallow(clonePath, rr.RemoteConfig.Url, rr.RemoteConfig.Branch); err != nil {
				return "", err
			} else {
				if err = rr.Git.Reset(clonePath, rr.RemoteConfig.Revision); err != nil {
					return "", err
				}
			}
		}
	}

	return clonePath, nil
}

func (rr *RemoteResolver) IsClonedPathExists(path string) bool {
	return ioutils.IsValidPath(path)
}

func (rr *RemoteResolver) verifyRemoteRepositoryConfig() error {
	if rr.RemoteConfig == nil {
		return fmt.Errorf("invalid remote repository configuration")
	}
	remoteCfg := rr.RemoteConfig
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
