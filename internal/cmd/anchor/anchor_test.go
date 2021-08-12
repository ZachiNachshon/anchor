package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AnchorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail to scan repo due to missing locator in registry",
			Func: FailToScanRepoDueToMissingLocatorFromRegistry,
		},
		{
			Name: "fail to scan repository",
			Func: FailToScanRepository,
		},
		{
			Name: "scan repository successfully",
			Func: ScanRepositorySuccessfully,
		},
		{
			Name: "use existing current context as the active config context",
			Func: UseExistingCurrentContentAsTheActiveConfigContext,
		},
		{
			Name: "fail to prompt for config context selection",
			Func: FailToPromptForConfigContextSelection,
		},
		{
			Name: "fail when config context was not selected",
			Func: FailWhenConfigContextWasNotSelected,
		},
		{
			Name: "clear screen upon config context selection",
			Func: ClearScreenUponConfigContextSelection,
		},
		{
			Name: "avoid failure when clear screen fails",
			Func: AvoidFailureWhenClearScreenFails,
		},
		{
			Name: "failed to resolve repository origin",
			Func: FailedToResolveRepositoryOrigin,
		},
		{
			Name: "failed to load repository",
			Func: FailedToLoadRepository,
		},
		{
			Name: "load repository files successfully",
			Func: LoadRepositoryFilesSuccessfully,
		},
		{
			Name: "run pre-run sequence successfully",
			Func: RunPreRunSequenceSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var FailToScanRepoDueToMissingLocatorFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		err := scanAnchorfilesRepositoryTree(ctx, "/some/path")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve from registry")
	})
}

var FailToScanRepository = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		path := "/some/path"
		fakeLocator := locator.CreateFakeLocator(path)
		fakeLocator.ScanMock = func(anchorFilesLocalPath string) error {
			return fmt.Errorf("failed to scan")
		}
		ctx.Registry().Set(locator.Identifier, fakeLocator)
		err := scanAnchorfilesRepositoryTree(ctx, path)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to scan")
	})
}

var ScanRepositorySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		path := "/some/path"
		fakeLocator := locator.CreateFakeLocator(path)
		fakeLocator.ScanMock = func(anchorFilesLocalPath string) error {
			return nil
		}
		ctx.Registry().Set(locator.Identifier, fakeLocator)
		err := scanAnchorfilesRepositoryTree(ctx, path)
		assert.Nil(t, err)
	})
}

var UseExistingCurrentContentAsTheActiveConfigContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			currCfgCtx := "test-curr-ctx"
			var items = config.TemplateItems{
				CurrentContext:   currCfgCtx,
				FirstContextName: currCfgCtx,
			}
			with.Config(ctx, config.GetCustomTestConfigText(items), func(cfg *config.AnchorConfig) {
				err := loadConfigContext(ctx, nil, nil)
				assert.Nil(t, err)
				assert.Equal(t, currCfgCtx, cfg.Config.ActiveContext.Name)
			})
		})
	})
}

var FailToPromptForConfigContextSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCfgCtxCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptConfigContextMock = func(cfgContexts []*config.Context) (*config.Context, error) {
				promptCfgCtxCallCount++
				return nil, fmt.Errorf("failed to prompt")
			}
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				cfg.Config.CurrentContext = ""
				err := loadConfigContext(ctx, fakePrompter, nil)
				assert.Equal(t, 1, promptCfgCtxCallCount, "expected action to be called exactly once. name: promptConfigContext")
				assert.NotNil(t, err)
				assert.Equal(t, "failed to prompt", err.Error())
			})
		})
	})
}

var FailWhenConfigContextWasNotSelected = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCfgCtxCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptConfigContextMock = func(cfgContexts []*config.Context) (*config.Context, error) {
				promptCfgCtxCallCount++
				return &config.Context{
					Name: prompter.CancelActionName,
				}, nil
			}
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				cfg.Config.CurrentContext = ""
				err := loadConfigContext(ctx, fakePrompter, nil)
				assert.Equal(t, 1, promptCfgCtxCallCount, "expected action to be called exactly once. name: promptConfigContext")
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "cannot proceed without selecting a configuration context")
			})
		})
	})
}

var ClearScreenUponConfigContextSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCfgCtxCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptConfigContextMock = func(cfgContexts []*config.Context) (*config.Context, error) {
				promptCfgCtxCallCount++
				return &config.Context{
					Name: cfgContexts[0].Name,
				}, nil
			}
			clearScreenCallCount := 0
			fakeShell := shell.CreateFakeShell()
			fakeShell.ClearScreenMock = func() error {
				clearScreenCallCount++
				return nil
			}
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				cfg.Config.CurrentContext = ""
				err := loadConfigContext(ctx, fakePrompter, fakeShell)
				assert.Equal(t, 1, promptCfgCtxCallCount, "expected action to be called exactly once. name: promptConfigContext")
				assert.Equal(t, 1, clearScreenCallCount, "expected action to be called exactly once. name: clearScreen")
				assert.Nil(t, err)
			})
		})
	})
}

var AvoidFailureWhenClearScreenFails = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCfgCtxCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptConfigContextMock = func(cfgContexts []*config.Context) (*config.Context, error) {
				promptCfgCtxCallCount++
				return &config.Context{
					Name: cfgContexts[0].Name,
				}, nil
			}
			clearScreenCallCount := 0
			fakeShell := shell.CreateFakeShell()
			fakeShell.ClearScreenMock = func() error {
				clearScreenCallCount++
				return fmt.Errorf("failed to clear screen")
			}
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				cfg.Config.CurrentContext = ""
				err := loadConfigContext(ctx, fakePrompter, fakeShell)
				assert.Equal(t, 1, promptCfgCtxCallCount, "expected action to be called exactly once. name: promptConfigContext")
				assert.Equal(t, 1, clearScreenCallCount, "expected action to be called exactly once. name: clearScreen")
				assert.Nil(t, err)
			})
		})
	})
}

var FailedToResolveRepositoryOrigin = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				cfg.Config.ActiveContext.Context.Repository = nil
				_, err := loadRepository(ctx)
				assert.NotNil(t, err)
			})
		})
	})
}

var FailedToLoadRepository = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				cfg.Config.ActiveContext.Context.Repository.Local.Path = "/invalid/path"
				_, err := loadRepository(ctx)
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "local anchorfiles repository path is invalid")
			})
		})
	})
}

var LoadRepositoryFilesSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				cfg.Config.ActiveContext.Context.Repository.Local.Path = ctx.AnchorFilesPath()
				repoPath, err := loadRepository(ctx)
				assert.Nil(t, err)
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath)
			})
		})
	})
}

var RunPreRunSequenceSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				fakePrompter := prompter.CreateFakePrompter()
				fakePrompter.PromptConfigContextMock = func(cfgContexts []*config.Context) (*config.Context, error) {
					return &config.Context{
						Name: cfgContexts[0].Name,
					}, nil
				}
				ctx.Registry().Set(prompter.Identifier, fakePrompter)

				fakeShell := shell.CreateFakeShell()
				fakeShell.ClearScreenMock = func() error {
					return nil
				}
				ctx.Registry().Set(shell.Identifier, fakeShell)

				fakeLocator := locator.CreateFakeLocator("/path/to/repo")
				fakeLocator.ScanMock = func(anchorFilesLocalPath string) error {
					return nil
				}
				ctx.Registry().Set(locator.Identifier, fakeLocator)

				with.HarnessAnchorfilesTestRepo(ctx)
				cfg.Config.ActiveContext.Context.Repository.Local.Path = ctx.AnchorFilesPath()
				preRunSeq := AnchorPreRunSequence()
				err := preRunSeq.Run(ctx)
				assert.Nil(t, err)
			})
		})
	})
}
