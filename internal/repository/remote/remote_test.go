package remote

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/git"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"

	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RemoteShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail resolving registry components",
			Func: FailResolvingRegistryComponents,
		},
		{
			Name: "verify remote repository config values",
			Func: VerifyRemoteRepositoryConfigValues,
		},
		{
			Name: "not clone repo if already exists",
			Func: NotCloneRepoIfAlreadyExists,
		},
		{
			Name: "fail to clone repository",
			Func: FailToCloneRepository,
		},
		{
			Name: "clone repository successfully",
			Func: CloneRepositorySuccessfully,
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
			Name: "auto update: fail to get local origin commit hash",
			Func: AutoUpdateFailToGetLocalOriginCommitHash,
		},
		{
			Name: "auto update: fail to get remote HEAD commit hash",
			Func: AutoUpdateFailToGetRemoteHeadCommitHash,
		},
		{
			Name: "auto update: fail to reset to revision",
			Func: AutoUpdateFailToResetToRevision,
		},
		{
			Name: "auto update: do not fail when revision diff print fails",
			Func: AutoUpdateDoNotFailWhenRevisionDiffPrintFails,
		},
		{
			Name: "auto update: run a successful already up to date flow",
			Func: AutoUpdateRunSuccessfulAlreadyUpToDateFlow,
		},
		{
			Name: "load: fail on preparations",
			Func: LoadFailOnPreparations,
		},
		{
			Name: "load: fail to verify configuration",
			Func: LoadFailToVerifyConfiguration,
		},
		{
			Name: "load: fail to clone repository",
			Func: LoadFailToCloneRepository,
		},
		{
			Name: "load: fail to reset to revision",
			Func: LoadFailToResetToRevision,
		},
		{
			Name: "load: fail to auto update repository",
			Func: LoadFailToAutoUpdateRepository,
		},
		{
			Name: "load: fail to checkout from branch",
			Func: LoadFailToCheckoutFromBranch,
		},
		{
			Name: "load: remote repository successfully",
			Func: LoadRemoteRepositorySuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		remote := NewRemoteRepository(nil)

		err := remote.prepareFunc(remote, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", printer.Identifier))
		reg.Set(printer.Identifier, printer.CreateFakePrinter())

		err = remote.prepareFunc(remote, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", shell.Identifier))
		reg.Set(shell.Identifier, shell.CreateFakeShell())

		err = remote.prepareFunc(remote, ctx)
		assert.Nil(t, err)
	})
}

var VerifyRemoteRepositoryConfigValues = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.Url = ""
				remoteCfg.Branch = ""
				remoteCfg.ClonePath = ""

				remote := NewRemoteRepository(nil)

				err := remote.verifyRemoteRepositoryConfigFunc(nil)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "invalid remote repository configuration", err.Error())

				err = remote.verifyRemoteRepositoryConfigFunc(remoteCfg)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: url", err.Error())

				remoteCfg.Url = "/some/url"
				err = remote.verifyRemoteRepositoryConfigFunc(remoteCfg)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: branch", err.Error())

				remoteCfg.Branch = "some-branch"
				err = remote.verifyRemoteRepositoryConfigFunc(remoteCfg)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote repository config is missing value. name: clonePath", err.Error())

				remoteCfg.ClonePath = "/some/clone/path"
				err = remote.verifyRemoteRepositoryConfigFunc(remoteCfg)
				assert.Nil(t, err, "expected to succeed on remote resolver")
			})
		})
	})
}

var NotCloneRepoIfAlreadyExists = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeGit := git.CreateFakeGit()
				cloneCallCount := 0
				fakeGit.CloneMock = func(url string, branch string, clonePath string) error {
					cloneCallCount++
					return nil
				}
				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit

				err := remote.cloneRepoIfMissingFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.Nil(t, err, "expected not to fail")
				assert.Equal(t, 0, cloneCallCount, "expected not to be called")
			})
		})
	})
}

var FailToCloneRepository = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopOnFailureCallCount := 0
				fakeSpinner.StopOnFailureMock = func(err error) {
					stopOnFailureCallCount++
				}
				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrepareCloneRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				fakeGit := git.CreateFakeGit()
				cloneCallCount := 0
				fakeGit.CloneMock = func(url string, branch string, clonePath string) error {
					cloneCallCount++
					assert.Equal(t, remoteCfg.Url, url)
					assert.Equal(t, remoteCfg.Branch, branch)
					assert.Equal(t, remoteCfg.ClonePath, clonePath)
					return fmt.Errorf("failed to clone")
				}
				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit
				remote.prntr = fakePrinter

				err := remote.cloneRepoIfMissingFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "failed to clone", err.Error())
				assert.Equal(t, 1, spinCallCount, "expected to be called exactly once")
				assert.Equal(t, 1, stopOnFailureCallCount, "expected to be called exactly once")
			})
		})
	})
}

var CloneRepositorySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopOnSuccessCallCount := 0
				fakeSpinner.StopOnSuccessMock = func() {
					stopOnSuccessCallCount++
				}
				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrepareCloneRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				fakeGit := git.CreateFakeGit()
				cloneCallCount := 0
				fakeGit.CloneMock = func(url string, branch string, clonePath string) error {
					assert.Equal(t, remoteCfg.Url, url)
					assert.Equal(t, remoteCfg.Branch, branch)
					assert.Equal(t, remoteCfg.ClonePath, clonePath)
					cloneCallCount++
					return nil
				}
				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit
				remote.prntr = fakePrinter

				err := remote.cloneRepoIfMissingFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, spinCallCount, "expected to be called exactly once")
				assert.Equal(t, 1, stopOnSuccessCallCount, "expected to be called exactly once")
			})
		})
	})
}

var ResetToRevisionOnFirstTrySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()
				remoteCfg.AutoUpdate = true

				fakeGit := git.CreateFakeGit()
				gitResetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					return nil
				}
				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit

				err := remote.resetToRevisionFunc(remote, remoteCfg.ClonePath, remoteCfg.Branch, remoteCfg.Revision)
				assert.Nil(t, err, "expected not to fail")
				assert.Equal(t, 1, gitResetCallCount)
			})
		})
	})
}

var FailToFetchAfterFirstTryToResetFails = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeGit := git.CreateFakeGit()
				gitResetCallCount := 0
				fakeGit.ResetMock = func(path string, revision string) error {
					gitResetCallCount++
					return fmt.Errorf("fail to reset to revision 1st try")
				}
				gitFetchCallCount := 0
				fakeGit.FetchShallowMock = func(path string, branch string) error {
					gitFetchCallCount++
					return fmt.Errorf("fail to fetch")
				}
				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit

				err := remote.resetToRevisionFunc(remote, remoteCfg.ClonePath, remoteCfg.Branch, remoteCfg.Revision)
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeGit := git.CreateFakeGit()
				gitResetCallCount := 0
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
				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit

				err := remote.resetToRevisionFunc(remote, remoteCfg.ClonePath, remoteCfg.Branch, remoteCfg.Revision)
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeGit := git.CreateFakeGit()
				gitResetCallCount := 0
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
				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit

				err := remote.resetToRevisionFunc(remote, remoteCfg.ClonePath, remoteCfg.Branch, remoteCfg.Revision)
				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "fail to reset to revision 2nd try", err.Error())
				assert.Equal(t, 2, gitResetCallCount)
				assert.Equal(t, 1, gitFetchCallCount)
			})
		})
	})
}

var AutoUpdateFailToGetLocalOriginCommitHash = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopOnFailureCallCount := 0
				fakeSpinner.StopOnFailureMock = func(err error) {
					stopOnFailureCallCount++
				}
				fakePrinter := printer.CreateFakePrinter()
				prepareSpinnerCallCount := 0
				fakePrinter.PrepareReadRemoteHeadCommitHashSpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					prepareSpinnerCallCount++
					return fakeSpinner
				}

				fakeGit := git.CreateFakeGit()
				gitLocalOriginCommitCallCount := 0
				fakeGit.GetLocalOriginCommitHashMock = func(path string, branch string) (string, error) {
					gitLocalOriginCommitCallCount++
					return "", fmt.Errorf("fail to get local origin commit hash")
				}

				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit
				remote.prntr = fakePrinter

				err := remote.autoUpdateRepositoryFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "fail to get local origin commit hash", err.Error())
				assert.Equal(t, 1, prepareSpinnerCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, gitLocalOriginCommitCallCount)
				assert.Equal(t, 1, stopOnFailureCallCount)
			})
		})
	})
}

var AutoUpdateFailToGetRemoteHeadCommitHash = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopOnFailureCallCount := 0
				fakeSpinner.StopOnFailureMock = func(err error) {
					stopOnFailureCallCount++
				}
				fakePrinter := printer.CreateFakePrinter()
				prepareSpinnerCallCount := 0
				fakePrinter.PrepareReadRemoteHeadCommitHashSpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					prepareSpinnerCallCount++
					return fakeSpinner
				}

				fakeGit := git.CreateFakeGit()
				gitLocalOriginCommitCallCount := 0
				fakeGit.GetLocalOriginCommitHashMock = func(path string, branch string) (string, error) {
					gitLocalOriginCommitCallCount++
					return "", nil
				}

				gitRemoteHeadCommitCallCount := 0
				fakeGit.GetRemoteHeadCommitHashMock = func(path string, repoUrl string, branch string) (string, error) {
					gitRemoteHeadCommitCallCount++
					return "", fmt.Errorf("fail to get remote HEAD commit hash")
				}

				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit
				remote.prntr = fakePrinter

				err := remote.autoUpdateRepositoryFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "fail to get remote HEAD commit hash", err.Error())
				assert.Equal(t, 1, prepareSpinnerCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, gitLocalOriginCommitCallCount)
				assert.Equal(t, 1, gitRemoteHeadCommitCallCount)
				assert.Equal(t, 1, stopOnFailureCallCount)
			})
		})
	})
}

var AutoUpdateFailToResetToRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopOnFailureCallCount := 0
				fakeSpinner.StopOnFailureMock = func(err error) {
					stopOnFailureCallCount++
				}
				fakePrinter := printer.CreateFakePrinter()
				prepareSpinnerCallCount := 0
				fakePrinter.PrepareReadRemoteHeadCommitHashSpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					prepareSpinnerCallCount++
					return fakeSpinner
				}

				fakeGit := git.CreateFakeGit()
				gitLocalOriginCommitCallCount := 0
				fakeGit.GetLocalOriginCommitHashMock = func(path string, branch string) (string, error) {
					gitLocalOriginCommitCallCount++
					return "", nil
				}

				gitRemoteHeadCommitCallCount := 0
				fakeGit.GetRemoteHeadCommitHashMock = func(path string, repoUrl string, branch string) (string, error) {
					gitRemoteHeadCommitCallCount++
					return "", nil
				}

				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit
				remote.prntr = fakePrinter

				resetToRevisionCallCount := 0
				remote.resetToRevisionFunc = func(rr *remoteRepositoryImpl, clonePath string, branch string, revision string) error {
					resetToRevisionCallCount++
					return fmt.Errorf("fail to reset to revision")
				}

				err := remote.autoUpdateRepositoryFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.NotNil(t, err, "expected to fail")
				assert.Equal(t, "fail to reset to revision", err.Error())
				assert.Equal(t, 1, prepareSpinnerCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, gitLocalOriginCommitCallCount)
				assert.Equal(t, 1, gitRemoteHeadCommitCallCount)
				assert.Equal(t, 1, resetToRevisionCallCount)
				assert.Equal(t, 1, stopOnFailureCallCount)
			})
		})
	})
}

var AutoUpdateDoNotFailWhenRevisionDiffPrintFails = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopOnSuccessCallCount := 0
				fakeSpinner.StopOnSuccessMock = func() {
					stopOnSuccessCallCount++
				}
				fakePrinter := printer.CreateFakePrinter()
				prepareSpinnerCallCount := 0
				fakePrinter.PrepareReadRemoteHeadCommitHashSpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					prepareSpinnerCallCount++
					return fakeSpinner
				}

				fakeGit := git.CreateFakeGit()
				gitLocalOriginCommitCallCount := 0
				fakeGit.GetLocalOriginCommitHashMock = func(path string, branch string) (string, error) {
					gitLocalOriginCommitCallCount++
					return "12345", nil
				}

				gitRemoteHeadCommitCallCount := 0
				fakeGit.GetRemoteHeadCommitHashMock = func(path string, repoUrl string, branch string) (string, error) {
					gitRemoteHeadCommitCallCount++
					return "abcdef", nil
				}

				logRevisionDiffCallCount := 0
				fakeGit.LogRevisionsDiffPrettyMock = func(path string, prevRevision string, newRevision string) error {
					logRevisionDiffCallCount++
					return fmt.Errorf("failed to log revision diff pretty")
				}

				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit
				remote.prntr = fakePrinter

				resetToRevisionCallCount := 0
				remote.resetToRevisionFunc = func(rr *remoteRepositoryImpl, clonePath string, branch string, revision string) error {
					resetToRevisionCallCount++
					return nil
				}

				err := remote.autoUpdateRepositoryFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, prepareSpinnerCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, gitLocalOriginCommitCallCount)
				assert.Equal(t, 1, gitRemoteHeadCommitCallCount)
				assert.Equal(t, 1, resetToRevisionCallCount)
				assert.Equal(t, 1, logRevisionDiffCallCount)
				assert.Equal(t, 1, stopOnSuccessCallCount)
			})
		})
	})
}

var AutoUpdateRunSuccessfulAlreadyUpToDateFlow = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.ClonePath = ctx.AnchorFilesPath()

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopOnCustomSuccessCallCount := 0
				fakeSpinner.StopOnSuccessWithCustomMessageMock = func(message string) {
					stopOnCustomSuccessCallCount++
				}
				fakePrinter := printer.CreateFakePrinter()
				prepareSpinnerCallCount := 0
				fakePrinter.PrepareReadRemoteHeadCommitHashSpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					prepareSpinnerCallCount++
					return fakeSpinner
				}

				fakeGit := git.CreateFakeGit()
				gitLocalOriginCommitCallCount := 0
				fakeGit.GetLocalOriginCommitHashMock = func(path string, branch string) (string, error) {
					gitLocalOriginCommitCallCount++
					return "12345", nil
				}

				gitRemoteHeadCommitCallCount := 0
				fakeGit.GetRemoteHeadCommitHashMock = func(path string, repoUrl string, branch string) (string, error) {
					gitRemoteHeadCommitCallCount++
					return "12345", nil
				}

				remote := NewRemoteRepository(remoteCfg)
				remote.git = fakeGit
				remote.prntr = fakePrinter

				resetToRevisionCallCount := 0
				remote.resetToRevisionFunc = func(rr *remoteRepositoryImpl, clonePath string, branch string, revision string) error {
					resetToRevisionCallCount++
					return nil
				}

				err := remote.autoUpdateRepositoryFunc(remote, remoteCfg.Url, remoteCfg.Branch, remoteCfg.ClonePath)
				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, prepareSpinnerCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, gitLocalOriginCommitCallCount)
				assert.Equal(t, 1, gitRemoteHeadCommitCallCount)
				assert.Equal(t, 1, resetToRevisionCallCount)
				assert.Equal(t, 1, stopOnCustomSuccessCallCount)
			})
		})
	})
}

var LoadFailOnPreparations = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remote := NewRemoteRepository(remoteCfg)
				prepareCallCount := 0
				remote.prepareFunc = func(rr *remoteRepositoryImpl, ctx common.Context) error {
					prepareCallCount++
					return fmt.Errorf("fail to prepare")
				}

				clonePath, err := remote.Load(ctx)
				assert.NotNil(t, err, "expected to fail")
				assert.Empty(t, clonePath)
				assert.Equal(t, "fail to prepare", err.Error())
				assert.Equal(t, 1, prepareCallCount)
			})
		})
	})
}

var LoadFailToVerifyConfiguration = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remote := NewRemoteRepository(remoteCfg)
				prepareCallCount := 0
				remote.prepareFunc = func(rr *remoteRepositoryImpl, ctx common.Context) error {
					prepareCallCount++
					return nil
				}
				verifyRepoConfigCallCount := 0
				remote.verifyRemoteRepositoryConfigFunc = func(remoteCfg *config.Remote) error {
					verifyRepoConfigCallCount++
					return fmt.Errorf("fail to verify remote repo configuration")
				}

				clonePath, err := remote.Load(ctx)
				assert.NotNil(t, err, "expected to fail")
				assert.Empty(t, clonePath)
				assert.Equal(t, "fail to verify remote repo configuration", err.Error())
				assert.Equal(t, 1, prepareCallCount)
				assert.Equal(t, 1, verifyRepoConfigCallCount)
			})
		})
	})
}

var LoadFailToCloneRepository = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remote := NewRemoteRepository(remoteCfg)
				prepareCallCount := 0
				remote.prepareFunc = func(rr *remoteRepositoryImpl, ctx common.Context) error {
					prepareCallCount++
					return nil
				}
				verifyRepoConfigCallCount := 0
				remote.verifyRemoteRepositoryConfigFunc = func(remoteCfg *config.Remote) error {
					verifyRepoConfigCallCount++
					return nil
				}
				cloneRepoCallCount := 0
				remote.cloneRepoIfMissingFunc = func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
					cloneRepoCallCount++
					return fmt.Errorf("fail to clone repo")
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				remote.prntr = fakePrinter

				clonePath, err := remote.Load(ctx)
				assert.NotNil(t, err, "expected to fail")
				assert.Empty(t, clonePath)
				assert.Equal(t, "fail to clone repo", err.Error())
				assert.Equal(t, 1, prepareCallCount)
				assert.Equal(t, 1, verifyRepoConfigCallCount)
				assert.Equal(t, 1, cloneRepoCallCount)
			})
		})
	})
}

var LoadFailToResetToRevision = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.Revision = "12345"
				remoteCfg.AutoUpdate = false

				remote := NewRemoteRepository(remoteCfg)
				prepareCallCount := 0
				remote.prepareFunc = func(rr *remoteRepositoryImpl, ctx common.Context) error {
					prepareCallCount++
					return nil
				}
				verifyRepoConfigCallCount := 0
				remote.verifyRemoteRepositoryConfigFunc = func(remoteCfg *config.Remote) error {
					verifyRepoConfigCallCount++
					return nil
				}
				cloneRepoCallCount := 0
				remote.cloneRepoIfMissingFunc = func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
					cloneRepoCallCount++
					return nil
				}
				resetToRevisionCallCount := 0
				remote.resetToRevisionFunc = func(rr *remoteRepositoryImpl, clonePath string, branch string, revision string) error {
					resetToRevisionCallCount++
					return fmt.Errorf("failed to reset to revision")
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				remote.prntr = fakePrinter

				clonePath, err := remote.Load(ctx)
				assert.NotNil(t, err, "expected to fail")
				assert.Empty(t, clonePath)
				assert.Equal(t, "failed to reset to revision", err.Error())
				assert.Equal(t, 1, prepareCallCount)
				assert.Equal(t, 1, verifyRepoConfigCallCount)
				assert.Equal(t, 1, cloneRepoCallCount)
				assert.Equal(t, 1, resetToRevisionCallCount)
			})
		})
	})
}

var LoadFailToAutoUpdateRepository = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.Revision = ""
				remoteCfg.AutoUpdate = true

				remote := NewRemoteRepository(remoteCfg)
				prepareCallCount := 0
				remote.prepareFunc = func(rr *remoteRepositoryImpl, ctx common.Context) error {
					prepareCallCount++
					return nil
				}
				verifyRepoConfigCallCount := 0
				remote.verifyRemoteRepositoryConfigFunc = func(remoteCfg *config.Remote) error {
					verifyRepoConfigCallCount++
					return nil
				}
				cloneRepoCallCount := 0
				remote.cloneRepoIfMissingFunc = func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
					cloneRepoCallCount++
					return nil
				}
				autoUpdateRepoCallCount := 0
				remote.autoUpdateRepositoryFunc = func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
					autoUpdateRepoCallCount++
					return fmt.Errorf("failed to auto update repository")
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				remote.prntr = fakePrinter

				clonePath, err := remote.Load(ctx)
				assert.NotNil(t, err, "expected to fail")
				assert.Empty(t, clonePath)
				assert.Equal(t, "failed to auto update repository", err.Error())
				assert.Equal(t, 1, prepareCallCount)
				assert.Equal(t, 1, verifyRepoConfigCallCount)
				assert.Equal(t, 1, cloneRepoCallCount)
				assert.Equal(t, 1, autoUpdateRepoCallCount)
			})
		})
	})
}

var LoadFailToCheckoutFromBranch = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.Revision = ""
				remoteCfg.AutoUpdate = false

				remote := NewRemoteRepository(remoteCfg)
				prepareCallCount := 0
				remote.prepareFunc = func(rr *remoteRepositoryImpl, ctx common.Context) error {
					prepareCallCount++
					return nil
				}
				verifyRepoConfigCallCount := 0
				remote.verifyRemoteRepositoryConfigFunc = func(remoteCfg *config.Remote) error {
					verifyRepoConfigCallCount++
					return nil
				}
				cloneRepoCallCount := 0
				remote.cloneRepoIfMissingFunc = func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
					cloneRepoCallCount++
					return nil
				}

				fakeGit := git.CreateFakeGit()
				gitCheckoutCallCount := 0
				fakeGit.CheckoutMock = func(path string, branch string) error {
					gitCheckoutCallCount++
					return fmt.Errorf("failed to checkout")
				}
				remote.git = fakeGit

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				remote.prntr = fakePrinter

				clonePath, err := remote.Load(ctx)
				assert.NotNil(t, err, "expected to fail")
				assert.Empty(t, clonePath)
				assert.Equal(t, "failed to checkout", err.Error())
				assert.Equal(t, 1, prepareCallCount)
				assert.Equal(t, 1, verifyRepoConfigCallCount)
				assert.Equal(t, 1, cloneRepoCallCount)
				assert.Equal(t, 1, gitCheckoutCallCount)
			})
		})
	})
}

var LoadRemoteRepositorySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				remoteCfg := cfg.Config.ActiveContext.Context.Repository.Remote
				remoteCfg.Revision = ""
				remoteCfg.AutoUpdate = false

				remote := NewRemoteRepository(remoteCfg)
				prepareCallCount := 0
				remote.prepareFunc = func(rr *remoteRepositoryImpl, ctx common.Context) error {
					prepareCallCount++
					return nil
				}
				verifyRepoConfigCallCount := 0
				remote.verifyRemoteRepositoryConfigFunc = func(remoteCfg *config.Remote) error {
					verifyRepoConfigCallCount++
					return nil
				}
				cloneRepoCallCount := 0
				remote.cloneRepoIfMissingFunc = func(rr *remoteRepositoryImpl, url string, branch string, clonePath string) error {
					cloneRepoCallCount++
					return nil
				}

				fakeGit := git.CreateFakeGit()
				gitCheckoutCallCount := 0
				fakeGit.CheckoutMock = func(path string, branch string) error {
					gitCheckoutCallCount++
					return nil
				}
				remote.git = fakeGit

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				remote.prntr = fakePrinter

				clonePath, err := remote.Load(ctx)
				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, clonePath, remoteCfg.ClonePath)
				assert.Equal(t, 1, prepareCallCount)
				assert.Equal(t, 1, verifyRepoConfigCallCount)
				assert.Equal(t, 1, cloneRepoCallCount)
				assert.Equal(t, 1, gitCheckoutCallCount)
			})
		})
	})
}
