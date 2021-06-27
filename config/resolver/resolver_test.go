package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/git"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
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
			Name: "fail to resolve remote repository due to invalid config",
			Func: FailToResolveRemoteRepositoryDueToInvalidConfig,
		},
		{
			Name: "fail to resolve remote repository due to invalid url",
			Func: FailToResolveRemoteRepositoryDueToInvalidUrl,
		},
		{
			Name: "fail to resolve remote repository due to invalid branch",
			Func: FailToResolveRemoteRepositoryDueToInvalidBranch,
		},
		{
			Name: "fail to resolve remote repository due to invalid clonePath",
			Func: FailToResolveRemoteRepositoryDueToInvalidClonePath,
		},
		{
			Name: "clone a remote repository into a clone path successfully",
			Func: CloneRemoteRepositoryIntoClonePathSuccessfully,
		},
		{
			Name: "fail to clone a remote repository due to git clone failure",
			Func: FailToCloneRemoteRepositoryDueToGitCloneFailure,
		},
		{
			Name: "clone a remote repository successfully on 2nd reset to revision after fetch",
			Func: CloneRemoteRepositorySuccessfullyOnSecondResetToRevisionAfterFetch,
		},
		{
			Name: "clone a remote repository successfully on 1st reset to revision",
			Func: CloneRemoteRepositorySuccessfullyOnFirstResetToRevision,
		},
		{
			Name: "fail to fetch a remote repository after 1st reset failures",
			Func: FailToFetchRemoteRepositoryAfterFirstResetFailure,
		},
		{
			Name: "fail to reset 2nd time after a successful fetch",
			Func: FailToResetSecondTimeAfterSuccessfulFetch,
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

var FailToResolveRemoteRepositoryDueToInvalidConfig = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
		remote:
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "invalid remote repository configuration", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToResolveRemoteRepositoryDueToInvalidUrl = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.Repository.Remote.Url = ""
				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: url", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToResolveRemoteRepositoryDueToInvalidBranch = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.Repository.Remote.Branch = ""
				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: branch", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToResolveRemoteRepositoryDueToInvalidClonePath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.Repository.Remote.ClonePath = ""
				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: clonePath", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var CloneRemoteRepositoryIntoClonePathSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    remote:
      url: https://github.com/ZachiNachshon/dummy-repo.git
      branch: some-branch
      clonePath: /not/exist/path/to/clone/new/repo
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				populateShellInRegistry(ctx)

				callCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(clonePath string) error {
					callCount++
					return nil
				}

				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
					Git:          fakeGit,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, 1, callCount, "expected git action to be called exactly once. name: clone")
				assert.Equal(t, "/not/exist/path/to/clone/new/repo", repoPath, "expected to have a repository path")
			})
		})
	})
}

var FailToCloneRemoteRepositoryDueToGitCloneFailure = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
  repository:
    remote:
      url: https://github.com/ZachiNachshon/dummy-repo.git
      branch: some-branch
      clonePath: /clone/repo/path
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				populateShellInRegistry(ctx)

				callCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(clonePath string) error {
					callCount++
					return fmt.Errorf("failed to perform git clone")
				}

				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
					Git:          fakeGit,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, 1, callCount, "expected git action to be called exactly once. name: clone")
				assert.Equal(t, "failed to perform git clone", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var CloneRemoteRepositorySuccessfullyOnFirstResetToRevision = func(t *testing.T) {
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
				populateShellInRegistry(ctx)

				cloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(clonePath string) error {
					cloneCallCount++
					return nil
				}

				resetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					resetCallCount++
					return nil
				}

				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
					Git:          fakeGit,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, 0, cloneCallCount, "expected git action not to be called. name: clone")
				assert.Equal(t, 1, resetCallCount, "expected git action to be called exactly once. name: reset")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var CloneRemoteRepositorySuccessfullyOnSecondResetToRevisionAfterFetch = func(t *testing.T) {
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
				populateShellInRegistry(ctx)

				cloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(clonePath string) error {
					cloneCallCount++
					return nil
				}

				resetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					resetCallCount++
					if resetCallCount == 1 {
						return fmt.Errorf("failed to git reset on 1st try")
					} else {
						return nil
					}
				}

				fetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, url string, branch string) error {
					fetchCallCount++
					return nil
				}

				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
					Git:          fakeGit,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, 0, cloneCallCount, "expected git action not to be called. name: clone")
				assert.Equal(t, 1, fetchCallCount, "expected git action to be called exactly once. name: fetch")
				assert.Equal(t, 2, resetCallCount, "expected git action to be called multiple times. name: reset")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var FailToFetchRemoteRepositoryAfterFirstResetFailure = func(t *testing.T) {
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
				populateShellInRegistry(ctx)

				cloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(clonePath string) error {
					cloneCallCount++
					return nil
				}

				resetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					resetCallCount++
					return fmt.Errorf("failed to git reset on 1st try")
				}

				fetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, url string, branch string) error {
					fetchCallCount++
					return fmt.Errorf("failed to git fetch shallow")
				}

				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
					Git:          fakeGit,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, 0, cloneCallCount, "expected git action not to be called. name: clone")
				assert.Equal(t, 1, resetCallCount, "expected git action to be called exactly once. name: reset")
				assert.Equal(t, 1, fetchCallCount, "expected git action to be called exactly once. name: fetch")
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToResetSecondTimeAfterSuccessfulFetch = func(t *testing.T) {
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
				populateShellInRegistry(ctx)

				cloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(clonePath string) error {
					cloneCallCount++
					return nil
				}

				resetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					resetCallCount++
					if resetCallCount == 1 {
						return fmt.Errorf("failed to git reset on 1st try")
					} else {
						return fmt.Errorf("failed to git reset on 2nd try")
					}
				}

				fetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, url string, branch string) error {
					fetchCallCount++
					return nil
				}

				rslvr := &RemoteResolver{
					RemoteConfig: cfg.Config.Repository.Remote,
					Git:          fakeGit,
				}
				repoPath, err := rslvr.ResolveRepository(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, 0, cloneCallCount, "expected git action not to be called. name: clone")
				assert.Equal(t, 1, fetchCallCount, "expected git action to be called exactly once. name: fetch")
				assert.Equal(t, 2, resetCallCount, "expected git action to be called exactly once. name: reset")
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

func populateShellInRegistry(ctx common.Context) {
	s := shell.New()
	shell.ToRegistry(ctx.Registry(), s)
}
