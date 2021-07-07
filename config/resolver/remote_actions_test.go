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
			Name: "fail to fetch HEAD revision",
			Func: FailedToFetchHeadRevision,
		},
		{
			Name: "fail to reset to HEAD revision",
			Func: FailedToResetToHeadRevision,
		},
		{
			Name: "fetch latest HEAD revision successfully",
			Func: FetchLatestHeadRevisionSuccessfully,
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
				cfg.Config.Repository.Remote = nil
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.Repository.Remote)
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
				cfg.Config.Repository.Remote.Url = ""
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.Repository.Remote)
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
				cfg.Config.Repository.Remote.Branch = ""
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.Repository.Remote)
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
				cfg.Config.Repository.Remote.ClonePath = ""
				fakeGit := git.CreateFakeGit()
				actions := NewRemoteActions(fakeGit)
				err := actions.VerifyRemoteRepositoryConfig(cfg.Config.Repository.Remote)
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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Url,
					cfg.Config.Repository.Remote.Branch)

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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Url,
					cfg.Config.Repository.Remote.Branch)

				assert.Nil(t, err, "expected not to fail")
				assert.Equal(t, 1, gitCloneCallCount)
			})
		})
	})
}

var DoNotCloneNewRepositoryWhenClonePathExists = func(t *testing.T) {
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
				gitCloneCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.CloneMock = func(url string, branch string, clonePath string) error {
					gitCloneCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.CloneRepositoryIfMissing(
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Url,
					cfg.Config.Repository.Remote.Branch)

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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch,
					cfg.Config.Repository.Remote.Revision)

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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch,
					cfg.Config.Repository.Remote.Revision)

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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch,
					cfg.Config.Repository.Remote.Revision)

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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch,
					cfg.Config.Repository.Remote.Revision)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "fail to reset to revision 2nd try", err.Error())
				assert.Equal(t, 2, gitResetCallCount)
				assert.Equal(t, 1, gitFetchCallCount)
			})
		})
	})
}

var FailedToFetchHeadRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				getHeadRevCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.GetHeadCommitHashMock = func(path string, branch string) (string, error) {
					getHeadRevCallCount++
					return "", fmt.Errorf("failed to get latest HEAD rev")
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryFetchHeadRevision(
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "failed to get latest HEAD rev", err.Error())
				assert.Equal(t, 1, getHeadRevCallCount)
			})
		})
	})
}

var FailedToResetToHeadRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				getHeadRevCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.GetHeadCommitHashMock = func(path string, branch string) (string, error) {
					getHeadRevCallCount++
					return "l33tf4k3c0mm1757r1n6", nil
				}
				// Used for inner call to actions.TryResetToRevision
				gitResetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					return fmt.Errorf("failed to reset on 1st attempt")
				}
				// Used for inner call to actions.TryResetToRevision
				gitFetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, branch string) error {
					gitFetchCallCount++
					return fmt.Errorf("failed to fetch")
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryFetchHeadRevision(
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch)

				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "failed to fetch", err.Error())
				assert.Equal(t, 1, getHeadRevCallCount)
				assert.Equal(t, 1, gitResetCallCount)
				assert.Equal(t, 1, gitFetchCallCount)
			})
		})
	})
}

var FetchLatestHeadRevisionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
 repository:
   remote:
     url: https://github.com/ZachiNachshon/dummy-repo.git
     branch: some-branch
     clonePath: /path/to/clone
`
			with.Config(ctx, yamlConfigText, func(cfg config.AnchorConfig) {
				getHeadRevCallCount := 0
				fakeGit := git.CreateFakeGit()
				fakeGit.GetHeadCommitHashMock = func(path string, branch string) (string, error) {
					getHeadRevCallCount++
					return "l33tf4k3c0mm1757r1n6", nil
				}
				// Used for inner call to actions.TryResetToRevision
				gitResetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					return nil
				}
				actions := NewRemoteActions(fakeGit)
				err := actions.TryFetchHeadRevision(
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch)

				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, getHeadRevCallCount)
				assert.Equal(t, 1, gitResetCallCount)
			})
		})
	})
}

var FailToCheckoutBranch = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := `
config:
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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch)

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
					cfg.Config.Repository.Remote.ClonePath,
					cfg.Config.Repository.Remote.Branch)

				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, checkoutCallCount)
			})
		})
	})
}
