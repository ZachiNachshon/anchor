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

func Test_ResolverShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail to get resolver due to invalid config",
			Func: FailToGetResolverDueToInvalidConfig,
		},
		{
			Name: "get local resolver from config successfully",
			Func: GetLocalResolverFromConfigSuccessfully,
		},
		{
			Name: "get remote resolver from config successfully",
			Func: GetRemoteResolverFromConfigSuccessfully,
		},
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
			Name: "reset to revision on existing cloned repo successfully",
			Func: ResetToRevisionOnExistingClonedRepoSuccessfully,
		},
		{
			Name: "fail resetting to revision on existing cloned repo",
			Func: FailResettingToRevisionOnExistingClonedRepo,
		},
		{
			Name: "auto update to latest HEAD revision successfully",
			Func: AutoUpdateToLatestHeadRevisionSuccessfully,
		},
		{
			Name: "fail auto updating to latest HEAD revision",
			Func: FailAutoUpdatingToLatestHeadRevision,
		},
	}
	harness.RunTests(t, tests)
}

var FailToGetResolverDueToInvalidConfig = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				res, err := GetResolverBasedOnConfig(cfg.Config.Repository)
				assert.Nil(t, res, "expected invalid resolver")
				assert.NotNil(t, err, "expected to fail getting a repository resolver")
				assert.Contains(t, err.Error(), "missing required config value")
			})
		})
	})
}

var GetLocalResolverFromConfigSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    remote:
      url: https://github.com/ZachiNachshon/dummy-repo.git      
      revision: l33tf4k3c0mm1757r1n6 
      branch: some-branch
      clonePath: /best/path/ever
    local:
      path: /local/path/wins
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				res, err := GetResolverBasedOnConfig(cfg.Config.Repository)
				assert.Equal(t, fmt.Sprintf("%T", &LocalResolver{}), fmt.Sprintf("%T", res), "expected a local resolver")
				assert.Nil(t, err, "expected getting a resolver successfully")
			})
		})
	})
}

var GetRemoteResolverFromConfigSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    remote:
      url: https://github.com/ZachiNachshon/dummy-repo.git      
      revision: l33tf4k3c0mm1757r1n6 
      branch: some-branch
      clonePath: /best/path/ever
    local:
      path: ""
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				res, err := GetResolverBasedOnConfig(cfg.Config.Repository)
				assert.Equal(t, fmt.Sprintf("%T", &RemoteResolver{}), fmt.Sprintf("%T", res), "expected a remote resolver")
				assert.Nil(t, err, "expected getting a resolver successfully")
			})
		})
	})
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
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
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
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, url string, branch string, revision string) error {
					tryResetToRevisionCallCount++
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
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, url string, branch string, revision string) error {
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

var AutoUpdateToLatestHeadRevisionSuccessfully = func(t *testing.T) {
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
				tryResetToRevisionCallCount := 0
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, url string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return nil
				}
				tryFetchHeadRevisionCallCount := 0
				fakeRemoteActions.TryFetchHeadRevisionMock = func(clonePath string, url string, branch string) error {
					tryFetchHeadRevisionCallCount++
					return nil
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 0, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchHeadRevisionCallCount)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var FailAutoUpdatingToLatestHeadRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: /some/clone/path
     autoUpdate: true
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
				fakeRemoteActions.TryResetToRevisionMock = func(clonePath string, url string, branch string, revision string) error {
					tryResetToRevisionCallCount++
					return nil
				}
				tryFetchHeadRevisionCallCount := 0
				fakeRemoteActions.TryFetchHeadRevisionMock = func(clonePath string, url string, branch string) error {
					tryFetchHeadRevisionCallCount++
					return fmt.Errorf("failed to auto update")
				}
				rslvr := &RemoteResolver{
					RemoteConfig:  cfg.Config.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 0, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchHeadRevisionCallCount)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "failed to auto update", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}
