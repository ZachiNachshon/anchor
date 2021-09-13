package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
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
			Name: "fail to scan repository",
			Func: FailToScanRepository,
		},
		{
			Name: "scan repository successfully",
			Func: ScanRepositorySuccessfully,
		},
		{
			Name: "fail to resolve config context due to missing config manager",
			Func: FailToResolveConfigContextDueToMissingConfigManager,
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
			Name: "fail to run pre-run sequence collaborators",
			Func: FailToRunPreRunSequenceCollaborators,
		},
		{
			Name: "run pre-run sequence successfully",
			Func: RunPreRunSequenceSuccessfully,
		},
		{
			Name: "prepare registry items successfully",
			Func: PrepareRegistryItemsSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var FailToScanRepository = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		c := createNoOpAnchorCollaborators()
		fakeLocator := locator.CreateFakeLocator()
		fakeLocator.ScanMock = func(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
			return errors.NewLocatorError(fmt.Errorf("failed to scan"))
		}
		c.l = fakeLocator
		err := scanAnchorfilesRepositoryTree(c, ctx, "/some/path")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to scan")
	})
}

var ScanRepositorySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		c := createNoOpAnchorCollaborators()
		fakeLocator := locator.CreateFakeLocator()
		fakeLocator.ScanMock = func(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
			return nil
		}
		c.l = fakeLocator
		err := scanAnchorfilesRepositoryTree(c, ctx, "/some/path")
		assert.Nil(t, err)
	})
}

var FailToResolveConfigContextDueToMissingConfigManager = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				ctx.Registry().Set(config.Identifier, nil)
				c := createNoOpAnchorCollaborators()
				err := resolveConfigContext(c, ctx)
				assert.NotNil(t, err)
				assert.Equal(t, fmt.Sprintf("failed to retrieve from registry. name: %s", config.Identifier), err.Error())
			})
		})
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
				c := createNoOpAnchorCollaborators()
				err := resolveConfigContext(c, ctx)
				assert.Nil(t, err)
				assert.Equal(t, currCfgCtx, cfg.Config.ActiveContext.Name)
			})
		})
	})
}

var FailToPromptForConfigContextSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				c := createNoOpAnchorCollaborators()
				promptCfgCtxCallCount := 0
				fakePrompter := prompter.CreateFakePrompter()
				fakePrompter.PromptConfigContextMock = func(cfgContexts []*config.Context) (*config.Context, error) {
					promptCfgCtxCallCount++
					return nil, fmt.Errorf("failed to prompt")
				}
				c.prmptr = fakePrompter
				cfg.Config.CurrentContext = ""
				err := resolveConfigContext(c, ctx)
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				c := createNoOpAnchorCollaborators()
				promptCfgCtxCallCount := 0
				fakePrompter := prompter.CreateFakePrompter()
				fakePrompter.PromptConfigContextMock = func(cfgContexts []*config.Context) (*config.Context, error) {
					promptCfgCtxCallCount++
					return &config.Context{
						Name: prompter.CancelActionName,
					}, nil
				}
				c.prmptr = fakePrompter
				cfg.Config.CurrentContext = ""
				err := resolveConfigContext(c, ctx)
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				c := createNoOpAnchorCollaborators()
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
				c.prmptr = fakePrompter
				c.s = fakeShell
				cfg.Config.CurrentContext = ""
				err := resolveConfigContext(c, ctx)
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				c := createNoOpAnchorCollaborators()
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
				c.prmptr = fakePrompter
				c.s = fakeShell
				cfg.Config.CurrentContext = ""
				err := resolveConfigContext(c, ctx)
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
				c := createNoOpAnchorCollaborators()
				cfg.Config.ActiveContext.Context.Repository = nil
				_, err := loadRepository(c, ctx)
				assert.NotNil(t, err)
			})
		})
	})
}

var FailedToLoadRepository = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				c := createNoOpAnchorCollaborators()
				cfg.Config.ActiveContext.Context.Repository.Local.Path = "/invalid/path"
				_, err := loadRepository(c, ctx)
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
				c := createNoOpAnchorCollaborators()
				cfg.Config.ActiveContext.Context.Repository.Local.Path = ctx.AnchorFilesPath()
				repoPath, err := loadRepository(c, ctx)
				assert.Nil(t, err)
				assert.Equal(t, ctx.AnchorFilesPath(), repoPath)
			})
		})
	})
}

var FailToRunPreRunSequenceCollaborators = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			collaborators := NewAnchorCollaborators()

			// Prepare registry items
			collaborators.prepareRegistryItemsFunc = func(c *AnchorCollaborators, ctx common.Context) error {
				return fmt.Errorf("failed to prepare registry items")
			}
			err := collaborators.Run(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, "failed to prepare registry items", err.Error())
			collaborators.prepareRegistryItemsFunc = func(c *AnchorCollaborators, ctx common.Context) error {
				return nil
			}

			// Resolve config context
			collaborators.resolveConfigContextFunc = func(c *AnchorCollaborators, ctx common.Context) error {
				return fmt.Errorf("failed to resolve config context")
			}
			err = collaborators.Run(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, "failed to resolve config context", err.Error())
			collaborators.resolveConfigContextFunc = func(c *AnchorCollaborators, ctx common.Context) error {
				return nil
			}

			// Load repository
			collaborators.loadRepositoryFunc = func(c *AnchorCollaborators, ctx common.Context) (string, error) {
				return "", fmt.Errorf("failed to load repository")
			}
			err = collaborators.Run(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, "failed to load repository", err.Error())
			collaborators.loadRepositoryFunc = func(c *AnchorCollaborators, ctx common.Context) (string, error) {
				return "", nil
			}

			// Scan anchorfiles repository
			collaborators.scanAnchorfilesFunc = func(c *AnchorCollaborators, ctx common.Context, repoPath string) error {
				return fmt.Errorf("failed to scan anchorfiles repo")
			}
			err = collaborators.Run(ctx)
			assert.NotNil(t, err)
			assert.Equal(t, "failed to scan anchorfiles repo", err.Error())
			collaborators.scanAnchorfilesFunc = func(c *AnchorCollaborators, ctx common.Context, repoPath string) error {
				return nil
			}

			err = collaborators.Run(ctx)
			assert.Nil(t, err)
		})
	})
}

var RunPreRunSequenceSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			collaborators := NewAnchorCollaborators()
			collaborators.prepareRegistryItemsFunc = func(c *AnchorCollaborators, ctx common.Context) error {
				return nil
			}
			collaborators.resolveConfigContextFunc = func(c *AnchorCollaborators, ctx common.Context) error {
				return nil
			}
			collaborators.loadRepositoryFunc = func(c *AnchorCollaborators, ctx common.Context) (string, error) {
				return "", nil
			}
			collaborators.scanAnchorfilesFunc = func(c *AnchorCollaborators, ctx common.Context, repoPath string) error {
				return nil
			}
			err := collaborators.Run(ctx)
			assert.Nil(t, err)
		})
	})
}

var PrepareRegistryItemsSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		c := createNoOpAnchorCollaborators()

		err := prepareRegistryItems(c, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", shell.Identifier))
		reg.Set(shell.Identifier, shell.CreateFakeShell())

		err = prepareRegistryItems(c, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", prompter.Identifier))
		reg.Set(prompter.Identifier, prompter.CreateFakePrompter())

		err = prepareRegistryItems(c, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", locator.Identifier))
		reg.Set(locator.Identifier, locator.CreateFakeLocator())

		err = prepareRegistryItems(c, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", extractor.Identifier))
		reg.Set(extractor.Identifier, extractor.CreateFakeExtractor())

		err = prepareRegistryItems(c, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", parser.Identifier))
		reg.Set(parser.Identifier, parser.CreateFakeParser())

		err = prepareRegistryItems(c, ctx)
		assert.Nil(t, err)
	})
}

func createNoOpAnchorCollaborators() *AnchorCollaborators {
	return &AnchorCollaborators{
		prepareRegistryItemsFunc: func(c *AnchorCollaborators, ctx common.Context) error {
			return nil
		},
		resolveConfigContextFunc: func(c *AnchorCollaborators, ctx common.Context) error {
			return nil
		},
		loadRepositoryFunc: func(c *AnchorCollaborators, ctx common.Context) (string, error) {
			return "", nil
		},
		scanAnchorfilesFunc: func(c *AnchorCollaborators, ctx common.Context, repoPath string) error {
			return nil
		},
	}
}
