package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RemoteShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "resolve local repository successfully",
			Func: ResolveLocalRepositorySuccessfully,
		},
		{
			Name: "fail to resolve local repository due to invalid path",
			Func: FailToResolveLocalRepositoryDueToInvalidPath,
		},
		{
			Name: "fail to resolve local repository due to missing config",
			Func: FailToResolveLocalRepositoryDueToMissingConfig,
		},
		{
			Name: "fail to resolve remote repository due to invalid remote actions",
			Func: FailToResolveRemoteRepositoryDueToInvalidRemoteActions,
		},
		{
			Name: "fail to resolve remote repository due to invalid config",
			Func: FailToResolveRemoteRepositoryDueToInvalidConfig,
		},
		{
			Name: "fail to clone a fresh remote repository into clone path",
			Func: FailToCloneFreshRemoteRepositoryIntoClonePath,
		},
		{
			Name: "perform an initial fresh remote repository clone into a clone path successfully",
			Func: PerformInitialFreshRemoteRepositoryCloneIntoClonePathSuccessfully,
		},
		{
			Name: "clone repository and fail on checkout",
			Func: CloneRepositoryAndFailOnCheckout,
		},
		{
			Name: "reset to revision on existing cloned repo successfully",
			Func: ResetToRevisionOnExistingClonedRepoSuccessfully,
		},
		{
			Name: "fail resetting to revision on existing cloned repo",
			Func: FailResettingToRevisionOnExistingClonedRepo,
		},
		{
			Name: "auto update: fetch remote HEAD revision successfully",
			Func: AutoUpdateToRemoteHeadRevisionSuccessfully,
		},
		{
			Name: "auto update: avoid printing commit log since revision is up to date",
			Func: AvoidPrintingCommitLogSinceRevisionIsAlreadyUpToDate,
		},
		{
			Name: "auto update: fails to fetch local origin revision",
			Func: AutoUpdateFailsToFetchLocalOriginRevision,
		},
		{
			Name: "auto update: fails to fetch remote HEAD revision",
			Func: AutoUpdateFailsToFetchRemoteHeadRevision,
		},
		{
			Name: "auto update: fails to reset to revision",
			Func: AutoUpdateFailsToResetToRevision,
		},
		{
			Name: "auto update: fails to print revision diff does not generate an error",
			Func: AutoUpdateFailsToPrintRevisionDiffDoesNotGenerateAnError,
		},
	}
	harness.RunTests(t, tests)
}

var ResolveLocalRepositorySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
  repository:
    local:
      path: %s
`, ctx.AnchorFilesPath())
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				rslvr := &LocalResolver{
					LocalConfig: cfg.Config.Repository.Local,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Nil(t, err, "expected to succeed on local resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected a valid local repository path")
			})
		})
	})
}

var FailToResolveLocalRepositoryDueToInvalidPath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    local:
      path: /invalid/path
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				rslvr := &LocalResolver{
					LocalConfig: cfg.Config.Repository.Local,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on local resolver")
				assert.Contains(t, err.Error(), "local anchorfiles repository path is invalid")
				assert.Equal(t, "", repoPath, "expected not to have a repository path path")
			})
		})
	})
}

var FailToResolveLocalRepositoryDueToMissingConfig = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    local:
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				rslvr := &LocalResolver{
					LocalConfig: cfg.Config.Repository.Local,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on local resolver")
				assert.Equal(t, "invalid local repository configuration", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToResolveRemoteRepositoryDueToInvalidRemoteActions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote actions weren't defined for remote resolver, cannot proceed", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToResolveRemoteRepositoryDueToInvalidConfig = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    remote:
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return fmt.Errorf("invalid config")
				}

				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "invalid config", err.Error())
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToCloneFreshRemoteRepositoryIntoClonePath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return fmt.Errorf("failed to clone")
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "failed to clone", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var PerformInitialFreshRemoteRepositoryCloneIntoClonePathSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				checkoutCallCount := 0
				fakeRemoteActions.TryCheckoutToBranchMock = func(clonePath string, branch string) error {
					checkoutCallCount++
					return nil
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var CloneRepositoryAndFailOnCheckout = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				checkoutCallCount := 0
				fakeRemoteActions.TryCheckoutToBranchMock = func(clonePath string, branch string) error {
					checkoutCallCount++
					return fmt.Errorf("failed to checkout branch")
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "failed to checkout branch", err.Error())
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var ResetToRevisionOnExistingClonedRepoSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     revision: l33tf4k3c0mm1757r1n6
     clonePath: %s
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryResetToRevisionCallCount := 0
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return nil
				}
				checkoutCallCount := 0
				fakeRemoteActions.TryCheckoutToBranchMock = func(clonePath string, branch string) error {
					checkoutCallCount++
					return nil
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var FailResettingToRevisionOnExistingClonedRepo = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     revision: l33tf4k3c0mm1757r1n6
     clonePath: /some/clone/path
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryResetToRevisionCallCount := 0
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return fmt.Errorf("failed resetting to revision")
				}

				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "failed resetting to revision", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var AutoUpdateToRemoteHeadRevisionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
     autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryFetchLocalOriginCallCount := 0
				fakeRemoteActions.TryFetchLocalOriginRevisionMock = func(clonePath string, branch string) (string, error) {
					tryFetchLocalOriginCallCount++
					return "local-origin-rev", nil
				}

				tryFetchRemoteHeadCallCount := 0
				fakeRemoteActions.TryFetchRemoteHeadRevisionMock = func(clonePath string, repoUrl string, branch string) (string, error) {
					tryFetchRemoteHeadCallCount++
					return "head-revision", nil
				}
				tryResetToRevisionCallCount := 0
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return nil
				}
				printRevisionsDiffCallCount := 0
				fakeRemoteActions.PrintRevisionsDiffMock = func(path string, prevRevision string, newRevision string) error {
					printRevisionsDiffCallCount++
					return nil
				}
				checkoutCallCount := 0
				fakeRemoteActions.TryCheckoutToBranchMock = func(clonePath string, branch string) error {
					checkoutCallCount++
					return nil
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 1, printRevisionsDiffCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var AvoidPrintingCommitLogSinceRevisionIsAlreadyUpToDate = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
     autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryFetchLocalOriginCallCount := 0
				fakeRemoteActions.TryFetchLocalOriginRevisionMock = func(clonePath string, branch string) (string, error) {
					tryFetchLocalOriginCallCount++
					return "head-revision", nil
				}

				tryFetchRemoteHeadCallCount := 0
				fakeRemoteActions.TryFetchRemoteHeadRevisionMock = func(clonePath string, repoUrl string, branch string) (string, error) {
					tryFetchRemoteHeadCallCount++
					return "head-revision", nil
				}
				tryResetToRevisionCallCount := 0
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return nil
				}
				printRevisionsDiffCallCount := 0
				fakeRemoteActions.PrintRevisionsDiffMock = func(path string, prevRevision string, newRevision string) error {
					printRevisionsDiffCallCount++
					return nil
				}
				checkoutCallCount := 0
				fakeRemoteActions.TryCheckoutToBranchMock = func(clonePath string, branch string) error {
					checkoutCallCount++
					return nil
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 0, printRevisionsDiffCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var AutoUpdateFailsToFetchLocalOriginRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
     autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryFetchLocalOriginCallCount := 0
				fakeRemoteActions.TryFetchLocalOriginRevisionMock = func(clonePath string, branch string) (string, error) {
					tryFetchLocalOriginCallCount++
					return "", fmt.Errorf("fail to fetch local origin revision")
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "fail to fetch local origin revision", err.Error())
				assert.Equal(t, "", repoPath, "expected to have invalid repository path")
			})
		})
	})
}

var AutoUpdateFailsToFetchRemoteHeadRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
     autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryFetchLocalOriginCallCount := 0
				fakeRemoteActions.TryFetchLocalOriginRevisionMock = func(clonePath string, branch string) (string, error) {
					tryFetchLocalOriginCallCount++
					return "local-origin-revision", nil
				}

				tryFetchRemoteHeadCallCount := 0
				fakeRemoteActions.TryFetchRemoteHeadRevisionMock = func(clonePath string, repoUrl string, branch string) (string, error) {
					tryFetchRemoteHeadCallCount++
					return "", fmt.Errorf("fail to fetch remote HEAD revision")
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "fail to fetch remote HEAD revision", err.Error())
				assert.Equal(t, "", repoPath, "expected to have invalid repository path")
			})
		})
	})
}

var AutoUpdateFailsToResetToRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
     autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryFetchLocalOriginCallCount := 0
				fakeRemoteActions.TryFetchLocalOriginRevisionMock = func(clonePath string, branch string) (string, error) {
					tryFetchLocalOriginCallCount++
					return "local-origin-revision", nil
				}

				tryFetchRemoteHeadCallCount := 0
				fakeRemoteActions.TryFetchRemoteHeadRevisionMock = func(clonePath string, repoUrl string, branch string) (string, error) {
					tryFetchRemoteHeadCallCount++
					return "head-revision", nil
				}
				tryResetToRevisionCallCount := 0
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return fmt.Errorf("failed to reset to revision")
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "failed to reset to revision", err.Error())
				assert.Equal(t, "", repoPath, "expected to have invalid repository path")
			})
		})
	})
}

var AutoUpdateFailsToPrintRevisionDiffDoesNotGenerateAnError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			harness.HarnessAnchorfilesRemoteGitTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: %s
     autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				verifyConfigCallCount := 0
				fakeRemoteActions := CreateFakeRemoteActions()
				fakeRemoteActions.VerifyRemoteRepositoryConfigMock = func(remoteCfg *config.Remote) error {
					verifyConfigCallCount++
					return nil
				}
				cloneRepoIfMissingCallCount := 0
				fakeRemoteActions.CloneRepositoryIfMissingMock = func(clonePath string, url string, branch string) error {
					cloneRepoIfMissingCallCount++
					return nil
				}
				tryFetchLocalOriginCallCount := 0
				fakeRemoteActions.TryFetchLocalOriginRevisionMock = func(clonePath string, branch string) (string, error) {
					tryFetchLocalOriginCallCount++
					return "local-origin-revision", nil
				}

				tryFetchRemoteHeadCallCount := 0
				fakeRemoteActions.TryFetchRemoteHeadRevisionMock = func(clonePath string, repoUrl string, branch string) (string, error) {
					tryFetchRemoteHeadCallCount++
					return "head-revision", nil
				}
				tryResetToRevisionCallCount := 0
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return nil
				}
				printRevisionsDiffCallCount := 0
				fakeRemoteActions.PrintRevisionsDiffMock = func(path string, prevRevision string, newRevision string) error {
					printRevisionsDiffCallCount++
					return fmt.Errorf("failed to print revision diff")
				}
				checkoutCallCount := 0
				fakeRemoteActions.TryCheckoutToBranchMock = func(clonePath string, branch string) error {
					checkoutCallCount++
					return nil
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 1, printRevisionsDiffCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}
