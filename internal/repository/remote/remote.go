package remote

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/git"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

type RemoteRepository interface {
	Load(ctx common.Context) (string, error)
}

type remoteRepositoryImpl struct {
	RemoteRepository

	remoteConfig *config.Remote
	prntr        printer.Printer
	git          git.Git

	prepareFunc                      func(rr *remoteRepositoryImpl, ctx common.Context) error
	verifyRemoteRepositoryConfigFunc func(remoteCfg *config.Remote) error
	cloneRepoIfMissingFunc           func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error
	resetToRevisionFunc              func(rr *remoteRepositoryImpl, clonePath string, branch string, revision string) error
	autoUpdateRepositoryFunc         func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error
	readRemoteHeadRevisionFunc       func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) (string, error)
}

func NewRemoteRepository(remoteConfig *config.Remote) *remoteRepositoryImpl {
	return &remoteRepositoryImpl{
		remoteConfig:                     remoteConfig,
		prepareFunc:                      prepare,
		verifyRemoteRepositoryConfigFunc: verifyRemoteRepositoryConfig,
		cloneRepoIfMissingFunc:           cloneRepoIfMissing,
		resetToRevisionFunc:              resetToRevision,
		readRemoteHeadRevisionFunc:       readRemoteHeadRevisionFunc,
		autoUpdateRepositoryFunc:         autoUpdateRepository,
	}
}

func (rr *remoteRepositoryImpl) Load(ctx common.Context) (string, error) {
	if err := rr.prepareFunc(rr, ctx); err != nil {
		return "", err
	}

	if err := rr.verifyRemoteRepositoryConfigFunc(rr.remoteConfig); err != nil {
		return "", err
	}

	clonePath := rr.remoteConfig.ClonePath
	url := rr.remoteConfig.Url
	branch := rr.remoteConfig.Branch
	rr.prntr.PrintEmptyLines(1)

	if err := rr.cloneRepoIfMissingFunc(rr, url, branch, clonePath); err != nil {
		return "", err
	}

	if len(rr.remoteConfig.Revision) > 0 {
		if err := rr.resetToRevisionFunc(rr, clonePath, branch, rr.remoteConfig.Revision); err != nil {
			return "", err
		}
		if rr.remoteConfig.AutoUpdate {
			msg := fmt.Sprintf("Mutually exclusive config values found: autoUpdate / revision. "+
				"To allow auto update from '%s' branch latest HEAD, remove the revision from config.",
				branch)
			logger.Warning(msg)
			rr.prntr.PrintWarning(msg)
		}
	} else if rr.remoteConfig.AutoUpdate {
		if err := rr.autoUpdateRepositoryFunc(rr, url, branch, clonePath); err != nil {
			return "", err
		}
	}

	if err := rr.git.Checkout(clonePath, branch); err != nil {
		return "", err
	}

	return clonePath, nil
}

func prepare(rr *remoteRepositoryImpl, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		rr.prntr = resolved.(printer.Printer)
	}

	if resolved, err := ctx.Registry().SafeGet(shell.Identifier); err != nil {
		return err
	} else {
		rr.git = git.New(resolved.(shell.Shell))
	}
	return nil
}

func verifyRemoteRepositoryConfig(remoteCfg *config.Remote) error {
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

func cloneRepoIfMissing(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
	if !ioutils.IsValidPath(clonePath) {
		spnr := rr.prntr.PrepareCloneRepositorySpinner(url, branch)
		spnr.Spin()
		logger.Infof("Fetching anchorfiles repository for the first time...")
		if err := rr.git.Clone(url, branch, clonePath); err != nil {
			spnr.StopOnFailure(err)
			return err
		}
		spnr.StopOnSuccess()
	}
	return nil
}

func resetToRevision(rr *remoteRepositoryImpl, clonePath string, branch string, revision string) error {
	if err := rr.git.Reset(clonePath, revision); err != nil {
		spnr := rr.prntr.PrepareResetToRevisionSpinner(revision)
		spnr.Spin()
		// TODO: identify a "revision does not exists" error code before fetching again
		if err = rr.git.FetchShallow(clonePath, branch); err != nil {
			spnr.StopOnFailureWithCustomMessage(fmt.Sprintf("Failed fetching repository (branch: %s)", branch))
			return err
		} else {
			if err = rr.git.Reset(clonePath, revision); err != nil {
				spnr.StopOnFailureWithCustomMessage(fmt.Sprintf(
					"Failed resetting to revision after fetching repository (branch: %s, revision: %s)",
					branch, revision))
				return err
			}
		}
		spnr.StopOnSuccess()
	} else {
		// No need for a spinner if resetting to an already fetched revision
		rr.prntr.PrintSuccess(fmt.Sprintf("Reset to revision %s", revision))
	}

	logger.Infof("Updated anchorfiles repo to revision. commit-hash: %s", rr.remoteConfig.Revision)
	return nil
}

func readRemoteHeadRevisionFunc(rr *remoteRepositoryImpl, url string, branch string, clonePath string) (string, error) {
	spnr := rr.prntr.PrepareReadRemoteHeadCommitHashSpinner(rr.remoteConfig.Url, branch)
	spnr.Spin()
	if headRevision, err := rr.git.GetRemoteHeadCommitHash(clonePath, url, branch); err != nil {
		spnr.StopOnFailure(err)
		return "", err
	} else {
		spnr.StopOnSuccess()
		return headRevision, nil
	}
}

func autoUpdateRepository(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
	logger.Info("Checking anchorfiles local origin revision...")
	originRevision, err := rr.git.GetLocalOriginCommitHash(clonePath, branch)
	if err != nil {
		return err
	}

	logger.Info("Checking anchorfiles remote HEAD revision...")
	headRevision, err := rr.readRemoteHeadRevisionFunc(rr, url, branch, clonePath)
	if err != nil {
		return err
	}

	logger.Infof("Trying to reset to revision. commit-hash: %s", headRevision)
	if err = rr.resetToRevisionFunc(rr, clonePath, branch, headRevision); err != nil {
		return err
	}

	if originRevision != headRevision {
		logger.Infof("Fetched remote HEAD revision. commit-hash: %s", headRevision)
		err = rr.git.LogRevisionsDiffPretty(clonePath, originRevision, headRevision)
		if err != nil {
			logger.Errorf("failed to print revisions diff. error: %s", err.Error())
			// Do not return an error if print fails
			//return "", err
		}
	} else {
		alreadyUpToDateMsg := "Remote repository is already up to date !"
		logger.Info(alreadyUpToDateMsg)
		rr.prntr.PrintSuccess(alreadyUpToDateMsg)
	}
	return nil
}
