package resolver

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/git"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ResolverActionsShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail due to invalid remote repo config",
			Func: FailDueToInvalidRemoteRepoConfig,
		},
		{
			Name: "fail to verify remote repo config due to invalid url",
			Func: FailToVerifyRemoteRepoConfigDueToInvalidUrl,
		},
		{
			Name: "fail to verify remote repo config due to invalid branch",
			Func: FailToVerifyRemoteRepoConfigDueToInvalidBranch,
		},
		{
			Name: "fail to verify remote repo config due to invalid clonePath",
			Func: FailToVerifyRemoteRepoConfigDueToInvalidClonePath,
		},
		{
			Name: "fail to clone a remote repository due to git clone failure",
			Func: FailToCloneRemoteRepositoryDueToGitCloneFailure,
		},
		{
			Name: "clone a remote repository successfully when clone path is missing",
			Func: CloneRemoteRepositorySuccessfullyWhenClonePathIsMissing,
		},
		{
			Name: "do not clone a new repository when clone path exists",
			Func: DoNotCloneNewRepositoryWhenClonePathExists,
		},
		{
			Name: "reset to revision on 1st try successfully",
			Func: ResetToRevisionOnFirstTrySuccessfully,
		},
		{
			Name: "fail to fetch after 1st try to reset fails",
			Func: FailToFetchAfterFirstTryToResetFails,
		},
		{
			Name: "reset to revision on 2ns try successfully",
			Func: ResetToRevisionOnSecondTrySuccessfully,
		},
		{
			Name: "fail to reset to revision on 2nd try",
			Func: FailToResetToRevisionOnSecondTry,
		},
		{
			Name: "fetch remote HEAD revision successfully",
			Func: FetchRemoteHeadRevisionSuccessfully,
		},
		{
			Name: "fail to fetch remote HEAD revision",
			Func: FailedToFetchRemoteHeadRevision,
		},
		{
			Name: "fetch local origin revision successfully",
			Func: FetchLocalOriginRevisionSuccessfully,
		},
		{
			Name: "fail to fetch local origin revision",
			Func: FailedToFetchLocalOriginRevision,
		},
		{
			Name: "print revisions diff successfully",
			Func: PrintRevisionsDiffSuccessfully,
		},
		{
			Name: "fail to checkout a branch",
			Func: FailToCheckoutBranch,
		},
		{
			Name: "checkout to branch successfully",
			Func: CheckoutToBranchSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var FailDueToInvalidRemoteRepoConfig = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.ActiveContext.Context.Repository.Remote = nil
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.ActiveContext.Context.Repository.Remote)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "invalid remote repository configuration", err.Error())
			})
		})
	})
}

var FailToVerifyRemoteRepoConfigDueToInvalidUrl = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.ActiveContext.Context.Repository.Remote.Url = ""
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.ActiveContext.Context.Repository.Remote)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: url", err.Error())
			})
		})
	})
}

var FailToVerifyRemoteRepoConfigDueToInvalidBranch = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.ActiveContext.Context.Repository.Remote.Branch = ""
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.ActiveContext.Context.Repository.Remote)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: branch", err.Error())
			})
		})
	})
}

var FailToVerifyRemoteRepoConfigDueToInvalidClonePath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath = ""
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.ActiveContext.Context.Repository.Remote)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: clonePath", err.Error())
			})
		})
	})
}

var FailToCloneRemoteRepositoryDueToGitCloneFailure = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository:  
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /new/clone/path
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				gitCloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(url string, branch string, clonePath string) error {
					gitCloneCallCount++
					return fmt.Errorf("failed to perform git clone")
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.CloneRepositoryIfMissing(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Url,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, 1, gitCloneCallCount)
				assert.Equal(t, "failed to perform git clone", err.Error())
			})
		})
	})
}

var CloneRemoteRepositorySuccessfullyWhenClonePathIsMissing = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /new/clone/path
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				gitCloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(url string, branch string, clonePath string) error {
					gitCloneCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.CloneRepositoryIfMissing(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Url,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.Nil(t, err, "expected not to fail")
				assert.Equal(t, 1, gitCloneCallCount)
			})
		})
	})
}

var DoNotCloneNewRepositoryWhenClonePathExists = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.HarnessAnchorfilesTestRepo(ctx)
			yamlConfigText := fmt.Sprintf(`
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository:
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: %s
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				gitCloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(url string, branch string, clonePath string) error {
					gitCloneCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.CloneRepositoryIfMissing(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Url,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.Nil(t, err, "expected not to fail")
				assert.Equal(t, 0, gitCloneCallCount)
			})
		})
	})
}

var ResetToRevisionOnFirstTrySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       revision: l33tf4k3c0mm1757r1n6
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				gitResetCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryResetToRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch,
					cfg.Config.ActiveContext.Context.Repository.Remote.Revision)

				assert.Nil(t, err, "expected not to fail")
				assert.Equal(t, 1, gitResetCallCount)
			})
		})
	})
}

var FailToFetchAfterFirstTryToResetFails = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       revision: l33tf4k3c0mm1757r1n6
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				gitResetCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					return fmt.Errorf("fail to reset to revision 1st try")
				}
				gitFetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, branch string) error {
					gitFetchCallCount++
					return fmt.Errorf("fail to fetch")
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryResetToRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch,
					cfg.Config.ActiveContext.Context.Repository.Remote.Revision)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "fail to fetch", err.Error())
				assert.Equal(t, 1, gitResetCallCount)
				assert.Equal(t, 1, gitFetchCallCount)
			})
		})
	})
}

var ResetToRevisionOnSecondTrySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       revision: l33tf4k3c0mm1757r1n6
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				gitResetCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					if gitResetCallCount == 1 {
						return fmt.Errorf("fail to reset to revision 1st try")
					}
					return nil
				}
				gitFetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, branch string) error {
					gitFetchCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryResetToRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch,
					cfg.Config.ActiveContext.Context.Repository.Remote.Revision)

				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 2, gitResetCallCount)
				assert.Equal(t, 1, gitFetchCallCount)
			})
		})
	})
}

var FailToResetToRevisionOnSecondTry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       revision: l33tf4k3c0mm1757r1n6
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				gitResetCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					if gitResetCallCount == 1 {
						return fmt.Errorf("fail to reset to revision 1st try")
					}
					return fmt.Errorf("fail to reset to revision 2nd try")
				}
				gitFetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, branch string) error {
					gitFetchCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryResetToRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch,
					cfg.Config.ActiveContext.Context.Repository.Remote.Revision)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "fail to reset to revision 2nd try", err.Error())
				assert.Equal(t, 2, gitResetCallCount)
				assert.Equal(t, 1, gitFetchCallCount)
			})
		})
	})
}

var FailedToFetchRemoteHeadRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				getHeadRevCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.GetRemoteHeadCommitHashMock = func(path string, repoUrl string, branch string) (string, error) {
					getHeadRevCallCount++
					return "", fmt.Errorf("failed to get latest HEAD rev")
				}
				actions := NewRemoteActions(fakeGit)
				revision, err := actions.TryFetchRemoteHeadRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Url,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "failed to get latest HEAD rev", err.Error())
				assert.Equal(t, "", revision, "expected to receive empty revision")
				assert.Equal(t, 1, getHeadRevCallCount)
			})
		})
	})
}

var FetchRemoteHeadRevisionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				getHeadRevCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.GetRemoteHeadCommitHashMock = func(path string, repoUrl string, branch string) (string, error) {
					getHeadRevCallCount++
					return "l33tf4k3c0mm1757r1n6", nil
				}

				actions := NewRemoteActions(fakeGit)
				revision, err := actions.TryFetchRemoteHeadRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Url,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, getHeadRevCallCount)
				assert.Equal(t, "l33tf4k3c0mm1757r1n6", revision, "expected to receive valid revision")
			})
		})
	})
}

var FetchLocalOriginRevisionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				getLocalOriginRevCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.GetLocalOriginCommitHashMock = func(path string, branch string) (string, error) {
					getLocalOriginRevCallCount++
					return "l33tf4k3c0mm1757r1n6", nil
				}
				actions := NewRemoteActions(fakeGit)
				revision, err := actions.TryFetchLocalOriginRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, getLocalOriginRevCallCount)
				assert.Equal(t, "l33tf4k3c0mm1757r1n6", revision, "expected to receive valid revision")
			})
		})
	})
}

var FailedToFetchLocalOriginRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				getLocalOriginRevCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.GetLocalOriginCommitHashMock = func(path string, branch string) (string, error) {
					getLocalOriginRevCallCount++
					return "", fmt.Errorf("failed to get local origin rev")
				}
				actions := NewRemoteActions(fakeGit)
				revision, err := actions.TryFetchLocalOriginRevision(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "failed to get local origin rev", err.Error())
				assert.Equal(t, "", revision, "expected to receive empty revision")
				assert.Equal(t, 1, getLocalOriginRevCallCount)
			})
		})
	})
}

var PrintRevisionsDiffSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				prevRevision := "l33tf4k3c0mm1757r1n6"
				headRevision := "head-revision"
				logRevDiffPrettyCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.LogRevisionsDiffPrettyMock = func(path string, prevRevision string, newRevision string) error {
					logRevDiffPrettyCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.PrintRevisionsDiff(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					prevRevision, headRevision)

				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, logRevDiffPrettyCallCount)
			})
		})
	})
}

var FailToCheckoutBranch = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				checkoutCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CheckoutMock = func(path string, branch string) error {
					checkoutCallCount++
					return fmt.Errorf("failed to checkout branch")
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryCheckoutToBranch(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "failed to checkout branch", err.Error())
				assert.Equal(t, 1, checkoutCallCount)
			})
		})
	})
}

var CheckoutToBranchSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 currentContext: test-cfg-ctx
 contexts:
  - name: test-cfg-ctx
    context:
     repository: 
      remote:
       url: https://github.com/ZachiNachshon/dummy-repo.git
       branch: some-branch
       clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				checkoutCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CheckoutMock = func(path string, branch string) error {
					checkoutCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryCheckoutToBranch(
					cfg.Config.ActiveContext.Context.Repository.Remote.ClonePath,
					cfg.Config.ActiveContext.Context.Repository.Remote.Branch)

				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, checkoutCallCount)
			})
		})
	})
}
