package remote

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/printer"

	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RemoteShould(t *testing.T) {
	tests := []harness.TestsHarness{

		{
			Name: "fail to resolve remote repository due to invalid remote actions",
			Func: FailToResolveRemoteRepositoryDueToInvalidRemoteActions,
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

var FailToResolveRemoteRepositoryDueToInvalidRemoteActions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
				repo := &RemoteRepository{
					RemoteConfig: cfg.Config.ActiveContext.Context.Repository.Remote,
				}
				repoPath, err := repo.Load(ctx)
				assert.NotNil(t, err, "expected to fail on remote resolver")
				assert.Equal(t, "remote actions weren't defined for remote resolver, cannot proceed", err.Error())
				assert.Equal(t, "", repoPath, "expected not to have a repository path")
			})
		})
	})
}

var FailToCloneFreshRemoteRepositoryIntoClonePath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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
				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := repo.Load(ctx)
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

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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
				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := repo.Load(ctx)
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

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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
				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := repo.Load(ctx)
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
           revision: l33tf4k3c0mm1757r1n6
           clonePath: %s
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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
				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := repo.Load(ctx)
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
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
           url: https://github.com/ZachiNachshon/dummy-repo.git
           branch: some-branch
           revision: l33tf4k3c0mm1757r1n6
           clonePath: /some/clone/path
`
			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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

				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
				}
				repoPath, err := repo.Load(ctx)
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
           autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopSuccessCallCount := 0
				fakeSpinner.StopOnSuccessMock = func() {
					stopSuccessCallCount++
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				fakePrinter.PrepareAutoUpdateRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
					Printer:       fakePrinter,
				}
				repoPath, err := repo.Load(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 1, printRevisionsDiffCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, stopSuccessCallCount)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var AvoidPrintingCommitLogSinceRevisionIsAlreadyUpToDate = func(t *testing.T) {
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
           autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopCustomSuccessCallCount := 0
				fakeSpinner.StopOnSuccessWithCustomMessageMock = func(message string) {
					assert.Contains(t, message, "already up to date")
					stopCustomSuccessCallCount++
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				fakePrinter.PrepareAutoUpdateRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
					Printer:       fakePrinter,
				}
				repoPath, err := repo.Load(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 0, printRevisionsDiffCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, stopCustomSuccessCallCount)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}

var AutoUpdateFailsToFetchLocalOriginRevision = func(t *testing.T) {
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
           autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopCustomSuccessCallCount := 0
				fakeSpinner.StopOnFailureMock = func(err error) {
					stopCustomSuccessCallCount++
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				fakePrinter.PrepareAutoUpdateRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
					Printer:       fakePrinter,
				}
				repoPath, err := repo.Load(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, stopCustomSuccessCallCount)
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
           autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopCustomSuccessCallCount := 0
				fakeSpinner.StopOnFailureMock = func(err error) {
					stopCustomSuccessCallCount++
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				fakePrinter.PrepareAutoUpdateRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
					Printer:       fakePrinter,
				}
				repoPath, err := repo.Load(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, stopCustomSuccessCallCount)
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
           autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopCustomSuccessCallCount := 0
				fakeSpinner.StopOnFailureMock = func(err error) {
					stopCustomSuccessCallCount++
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				fakePrinter.PrepareAutoUpdateRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
					Printer:       fakePrinter,
				}
				repoPath, err := repo.Load(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, stopCustomSuccessCallCount)
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
           autoUpdate: true
`, ctx.AnchorFilesPath())

			with.Config(ctx, yamlConfigText, func(cfg *config.AnchorConfig) {
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

				fakeSpinner := printer.CreateFakePrinterSpinner()
				spinCallCount := 0
				fakeSpinner.SpinMock = func() {
					spinCallCount++
				}
				stopSuccessCallCount := 0
				fakeSpinner.StopOnSuccessMock = func() {
					stopSuccessCallCount++
				}

				fakePrinter := printer.CreateFakePrinter()
				fakePrinter.PrintEmptyLinesMock = func(count int) {}
				fakePrinter.PrepareAutoUpdateRepositorySpinnerMock = func(url string, branch string) printer.PrinterSpinner {
					return fakeSpinner
				}

				repo := &RemoteRepository{
					RemoteConfig:  cfg.Config.ActiveContext.Context.Repository.Remote,
					RemoteActions: fakeRemoteActions,
					Printer:       fakePrinter,
				}
				repoPath, err := repo.Load(ctx)
				assert.Equal(t, 1, verifyConfigCallCount)
				assert.Equal(t, 1, cloneRepoIfMissingCallCount)
				assert.Equal(t, 1, tryResetToRevisionCallCount)
				assert.Equal(t, 1, tryFetchRemoteHeadCallCount)
				assert.Equal(t, 1, tryFetchLocalOriginCallCount)
				assert.Equal(t, 1, printRevisionsDiffCallCount)
				assert.Equal(t, 1, checkoutCallCount)
				assert.Equal(t, 1, spinCallCount)
				assert.Equal(t, 1, stopSuccessCallCount)
				assert.Nil(t, err, "expected to succeed on remote resolver")
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath, "expected to have a repository path")
			})
		})
	})
}
