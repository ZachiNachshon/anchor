package main

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MainShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "set valid collaborators",
			Func: SetValidCollaborators,
		},
		{
			Name: "fail logger collaborator",
			Func: FailLoggerCollaborator,
		},
		{
			Name: "fail configuration collaborator",
			Func: FailConfigurationCollaborator,
		},
		{
			Name: "run registry collaborator successfully",
			Func: RunRegistryCollaboratorSuccessfully,
		},
		{
			Name: "run start cli commands collaborator successfully",
			Func: RunStartCliCommandsCollaboratorSuccessfully,
		},
		//{
		//	Name: "exit application successfully",
		//	Func: ExitApplicationSuccessfully,
		//},
		{
			Name: "start cli commands successfully",
			Func: StartCliCommandsSuccessfully,
		},
		{
			Name: "start main entry point successfully",
			Func: StartMainEntryPointSuccessfully,
		},
		{
			Name: "run collaborators in a specific order",
			Func: RunCollaboratorsInASpecificOrder,
		},
		{
			Name: "fail to run collaborators in sequence",
			Func: FailToRunCollaboratorsInSequence,
		},
		{
			Name: "initialize logger successfully",
			Func: InitializeLoggerSuccessfully,
		},
		{
			Name: "fail to append stdout logger",
			Func: FailToAppendStdoutLogger,
		},
		{
			Name: "fail to append file logger",
			Func: FailToAppendFileLogger,
		},
		{
			Name: "fail to set active logger",
			Func: FailToSetActiveLogger,
		},
		{
			Name: "fail to load log file path",
			Func: FailToAppendStdoutLogger,
		},
		{
			Name: "initialize configuration successfully",
			Func: InitializeConfigurationSuccessfully,
		},
		{
			Name: "fail to set up config file loader",
			Func: FailToSetupConfigFileLoader,
		},
		{
			Name: "fail to create config object",
			Func: FailToCreateConfigObject,
		},
		{
			Name: "initialize registry successfully",
			Func: InitializeRegistrySuccessfully,
		},
		{
			Name: "run pre run sequence: do nothing when no args",
			Func: PreRunSequenceDoNothingWhenNoArgs,
		},
		{
			Name: "run pre run sequence: run for root command",
			Func: PreRunSequenceRunForRootCommand,
		},
		{
			Name: "run pre run sequence: do not run for excluded command",
			Func: PreRunSequenceDoNotRunForExcludedCommand,
		},
		{
			Name: "run pre run sequence: run for non excluded command",
			Func: PreRunSequenceRunForNonExcludedCommand,
		},
	}
	harness.RunTests(t, tests)
}

var SetValidCollaborators = func(t *testing.T) {
	col := GetCollaborators()
	assert.NotNil(t, col.Logger, "expected collaborator not to be empty. name: Logger")
	assert.NotNil(t, col.Configuration, "expected collaborator not to be empty. name: Configuration")
	assert.NotNil(t, col.Registry, "expected collaborator not to be empty. name: Registry")
	assert.NotNil(t, col.StartCliCommands, "expected collaborator not to be empty. name: StartCliCommands")
}

var FailLoggerCollaborator = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeLoggerManager := logger.CreateFakeLoggerManager()
		fakeLoggerManager.CreateEmptyLoggerMock = func() (logger.Logger, error) {
			return nil, fmt.Errorf("fail to create logger")
		}
		err := collaborators.Logger(ctx, fakeLoggerManager)
		assert.NotNil(t, err)
		assert.Equal(t, "fail to create logger", err.Error())
	})
}

var FailConfigurationCollaborator = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeConfigManager := config.CreateFakeConfigManager()
		fakeConfigManager.SetupConfigFileLoaderMock = func() error {
			return fmt.Errorf("fail to create config file loader")
		}
		err := collaborators.Configuration(ctx, fakeConfigManager)
		assert.NotNil(t, err)
		assert.Equal(t, "fail to create config file loader", err.Error())
	})
}

var RunRegistryCollaboratorSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		err := collaborators.Registry(ctx)
		assert.Nil(t, err)
	})
}

var RunStartCliCommandsCollaboratorSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				// To pass local repo path validation
				cfg.Config.ActiveContext.Context.Repository.Local.Path = ctx.AnchorFilesPath()

				// Do not scan actual repo, use mocks
				fakeLocator := locator.CreateFakeLocator()
				fakeLocator.ScanMock = func(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
					return nil
				}
				fakeLocator.AnchorFoldersMock = func() []*models.AnchorFolderInfo {
					return stubs.GenerateAnchorFolderInfoTestData()
				}
				reg := ctx.Registry()
				reg.Set(locator.Identifier, fakeLocator)
				reg.Set(prompter.Identifier, prompter.CreateFakePrompter())
				reg.Set(shell.Identifier, shell.CreateFakeShell())
				reg.Set(extractor.Identifier, extractor.CreateFakeExtractor())
				reg.Set(parser.Identifier, parser.CreateFakeParser())
				shouldStartPreRunSeq := false
				err := collaborators.StartCliCommands(ctx, shouldStartPreRunSeq)
				assert.Nil(t, err)
			})
		})
	})
}

// testing os.exit based on:
// https://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go
//var ExitApplicationSuccessfully = func(t *testing.T) {
//	if os.Getenv("EXIT_APPLICATION") == "1" {
//		exitApplication(1, "an error occurred")
//		return
//	}
//	cmd := exec.Command(os.Args[0], "-test.run=ExitApplicationSuccessfully")
//	cmd.Env = append(os.Environ(), "EXIT_APPLICATION=1")
//	err := cmd.Run()
//	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
//		return
//	}
//	t.Fatalf("process ran with err %v, want exit status 1", err)
//}

var StartCliCommandsSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				// To pass local repo path validation
				cfg.Config.ActiveContext.Context.Repository.Local.Path = ctx.AnchorFilesPath()

				// Do not scan actual repo, use mocks
				fakeLocator := locator.CreateFakeLocator()
				fakeLocator.ScanMock = func(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
					return nil
				}
				fakeLocator.AnchorFoldersMock = func() []*models.AnchorFolderInfo {
					return stubs.GenerateAnchorFolderInfoTestData()
				}
				reg := ctx.Registry()
				reg.Set(locator.Identifier, fakeLocator)
				reg.Set(prompter.Identifier, prompter.CreateFakePrompter())
				reg.Set(shell.Identifier, shell.CreateFakeShell())
				reg.Set(extractor.Identifier, extractor.CreateFakeExtractor())
				reg.Set(parser.Identifier, parser.CreateFakeParser())
				shouldStartPreRunSeq := false
				err := startCliCommands(ctx, shouldStartPreRunSeq)
				assert.Nil(t, err)
			})
		})
	})
}

var StartMainEntryPointSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				GetCollaborators = func() *MainCollaborators {
					return &MainCollaborators{
						Logger: func(ctx common.Context, loggerManager logger.LoggerManager) error {
							return nil
						},
						Configuration: func(ctx common.Context, configManager config.ConfigManager) error {
							return nil
						},
						Registry: func(ctx common.Context) error {
							return nil
						},
						StartCliCommands: func(ctx common.Context, shouldStartPreRunSeq bool) error {
							return nil
						},
					}
				}
				main()
				// set with previous collaborators if a future testing method should be created for the main entry point
				GetCollaborators = func() *MainCollaborators {
					return collaborators
				}
			})
		})
	})
}

var RunCollaboratorsInASpecificOrder = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var callOrder []string
		loggerCallCount := 0
		configCallCount := 0
		registryCallCount := 0
		startCallCount := 0
		testCollaborators := &MainCollaborators{
			Logger: func(ctx common.Context, loggerManager logger.LoggerManager) error {
				callOrder = append(callOrder, "logger")
				loggerCallCount++
				return nil
			},
			Configuration: func(ctx common.Context, configManager config.ConfigManager) error {
				callOrder = append(callOrder, "configuration")
				configCallCount++
				return nil
			},
			Registry: func(ctx common.Context) error {
				callOrder = append(callOrder, "registry")
				registryCallCount++
				return nil
			},
			StartCliCommands: func(ctx common.Context, shouldStartPreRunSeq bool) error {
				callOrder = append(callOrder, "start")
				startCallCount++
				return nil
			},
		}
		shouldStartPreRunSeq := false
		err := runCollaboratorsInSequence(ctx, testCollaborators, shouldStartPreRunSeq)
		assert.Nil(t, err)
		assert.Equal(t, 1, loggerCallCount, "expected collaborator to be called exactly once. name: logger")
		assert.Equal(t, 1, configCallCount, "expected collaborator to be called exactly once. name: configuration")
		assert.Equal(t, 1, registryCallCount, "expected collaborator to be called exactly once. name: registry")
		assert.Equal(t, 1, startCallCount, "expected collaborator to be called exactly once. name: start")
		assert.Equal(t, 4, len(callOrder), "expected x4 collaborators to run")
		assert.Equal(t, "logger", callOrder[0], "expected collaborator to be in order: logger")
		assert.Equal(t, "configuration", callOrder[1], "expected collaborator to be in order: configuration")
		assert.Equal(t, "registry", callOrder[2], "expected collaborator to be in order: registry")
		assert.Equal(t, "start", callOrder[3], "expected collaborator to be in order: start")
	})
}

var FailToRunCollaboratorsInSequence = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		testCollaborators := &MainCollaborators{
			Logger: func(ctx common.Context, loggerManager logger.LoggerManager) error {
				return fmt.Errorf("failed to init logger")
			},
		}
		shouldStartPreRunSeq := false
		err := runCollaboratorsInSequence(ctx, testCollaborators, shouldStartPreRunSeq)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to init logger", err.Error())

		testCollaborators = &MainCollaborators{
			Logger: func(ctx common.Context, loggerManager logger.LoggerManager) error {
				return nil
			},
			Configuration: func(ctx common.Context, configManager config.ConfigManager) error {
				return fmt.Errorf("failed to init configuration")
			},
		}
		err = runCollaboratorsInSequence(ctx, testCollaborators, shouldStartPreRunSeq)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to init configuration", err.Error())

		testCollaborators = &MainCollaborators{
			Logger: func(ctx common.Context, loggerManager logger.LoggerManager) error {
				return nil
			},
			Configuration: func(ctx common.Context, configManager config.ConfigManager) error {
				return nil
			},
			Registry: func(ctx common.Context) error {
				return fmt.Errorf("failed to init registry")
			},
		}
		err = runCollaboratorsInSequence(ctx, testCollaborators, shouldStartPreRunSeq)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to init registry", err.Error())

		testCollaborators = &MainCollaborators{
			Logger: func(ctx common.Context, loggerManager logger.LoggerManager) error {
				return nil
			},
			Configuration: func(ctx common.Context, configManager config.ConfigManager) error {
				return nil
			},
			Registry: func(ctx common.Context) error {
				return nil
			},
			StartCliCommands: func(ctx common.Context, shouldStartPreRunSeq bool) error {
				return fmt.Errorf("failed to start cli commands")
			},
		}
		err = runCollaboratorsInSequence(ctx, testCollaborators, shouldStartPreRunSeq)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to start cli commands", err.Error())
	})
}

var InitializeLoggerSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeLogger, _ := logger.CreateFakeTestingLogger(t, false)
		fakeLogManager := logger.CreateFakeLoggerManager()
		createEmptyLoggerCallCount := 0
		fakeLogManager.CreateEmptyLoggerMock = func() (logger.Logger, error) {
			createEmptyLoggerCallCount++
			return fakeLogger, nil
		}
		appendStdoutLoggerCallCount := 0
		fakeLogManager.AppendStdoutLoggerMock = func(level string) (logger.Logger, error) {
			appendStdoutLoggerCallCount++
			return fakeLogger, nil
		}
		appendFileLoggerCallCount := 0
		fakeLogManager.AppendFileLoggerMock = func(level string) (logger.Logger, error) {
			appendFileLoggerCallCount++
			return fakeLogger, nil
		}
		setActiveLoggerCallCount := 0
		fakeLogManager.SetActiveLoggerMock = func(log *logger.Logger) error {
			setActiveLoggerCallCount++
			return nil
		}
		err := initLogger(ctx, fakeLogManager)
		assert.Nil(t, err)
		assert.Equal(t, 1, createEmptyLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, appendStdoutLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, appendFileLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, setActiveLoggerCallCount, "expected func to be called exactly once")
	})
}

var FailToAppendStdoutLogger = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeLogger, _ := logger.CreateFakeTestingLogger(t, false)
		fakeLogManager := logger.CreateFakeLoggerManager()
		createEmptyLoggerCallCount := 0
		fakeLogManager.CreateEmptyLoggerMock = func() (logger.Logger, error) {
			createEmptyLoggerCallCount++
			return fakeLogger, nil
		}
		appendStdoutLoggerCallCount := 0
		fakeLogManager.AppendStdoutLoggerMock = func(level string) (logger.Logger, error) {
			appendStdoutLoggerCallCount++
			return nil, fmt.Errorf("failed to append stdout logger")
		}
		err := initLogger(ctx, fakeLogManager)
		assert.NotNil(t, err)
		assert.Equal(t, 1, createEmptyLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, appendStdoutLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, "failed to append stdout logger", err.Error())
	})
}

var FailToAppendFileLogger = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeLogger, _ := logger.CreateFakeTestingLogger(t, false)
		fakeLogManager := logger.CreateFakeLoggerManager()
		createEmptyLoggerCallCount := 0
		fakeLogManager.CreateEmptyLoggerMock = func() (logger.Logger, error) {
			createEmptyLoggerCallCount++
			return fakeLogger, nil
		}
		appendStdoutLoggerCallCount := 0
		fakeLogManager.AppendStdoutLoggerMock = func(level string) (logger.Logger, error) {
			appendStdoutLoggerCallCount++
			return fakeLogger, nil
		}
		appendFileLoggerCallCount := 0
		fakeLogManager.AppendFileLoggerMock = func(level string) (logger.Logger, error) {
			appendFileLoggerCallCount++
			return nil, fmt.Errorf("failed to append file logger")
		}
		err := initLogger(ctx, fakeLogManager)
		assert.NotNil(t, err)
		assert.Equal(t, 1, createEmptyLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, appendStdoutLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, appendFileLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, "failed to append file logger", err.Error())
	})
}

var FailToSetActiveLogger = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeLogger, _ := logger.CreateFakeTestingLogger(t, false)
		fakeLogManager := logger.CreateFakeLoggerManager()
		createEmptyLoggerCallCount := 0
		fakeLogManager.CreateEmptyLoggerMock = func() (logger.Logger, error) {
			createEmptyLoggerCallCount++
			return fakeLogger, nil
		}
		appendStdoutLoggerCallCount := 0
		fakeLogManager.AppendStdoutLoggerMock = func(level string) (logger.Logger, error) {
			appendStdoutLoggerCallCount++
			return fakeLogger, nil
		}
		appendFileLoggerCallCount := 0
		fakeLogManager.AppendFileLoggerMock = func(level string) (logger.Logger, error) {
			appendFileLoggerCallCount++
			return fakeLogger, nil
		}
		setActiveLoggerCallCount := 0
		fakeLogManager.SetActiveLoggerMock = func(log *logger.Logger) error {
			setActiveLoggerCallCount++
			return fmt.Errorf("failed to set active logger")
		}
		err := initLogger(ctx, fakeLogManager)
		assert.NotNil(t, err)
		assert.Equal(t, 1, createEmptyLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, appendStdoutLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, appendFileLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, setActiveLoggerCallCount, "expected func to be called exactly once")
		assert.Equal(t, "failed to set active logger", err.Error())
	})
}

var InitializeConfigurationSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		configInUse := &config.AnchorConfig{}
		fakeCfgMgr := config.CreateFakeConfigManager()
		configLoaderCallCount := 0
		fakeCfgMgr.SetupConfigFileLoaderMock = func() error {
			configLoaderCallCount++
			return nil
		}
		configListenChangesCallCount := 0
		fakeCfgMgr.ListenOnConfigFileChangesMock = func(ctx common.Context) {
			configListenChangesCallCount++
		}
		createConfigCallCount := 0
		fakeCfgMgr.CreateConfigObjectMock = func() (*config.AnchorConfig, error) {
			createConfigCallCount++
			return configInUse, nil
		}
		err := initConfiguration(ctx, fakeCfgMgr)
		assert.Nil(t, err)
		assert.Equal(t, 1, configLoaderCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, configListenChangesCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, createConfigCallCount, "expected func to be called exactly once")
		assert.Equal(t, configInUse, ctx.Config().(*config.AnchorConfig))
	})
}

var FailToSetupConfigFileLoader = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		configLoaderCallCount := 0
		fakeCfgMgr := config.CreateFakeConfigManager()
		fakeCfgMgr.SetupConfigFileLoaderMock = func() error {
			configLoaderCallCount++
			return fmt.Errorf("failed to load config")
		}
		err := initConfiguration(ctx, fakeCfgMgr)
		assert.NotNil(t, err)
		assert.Equal(t, 1, configLoaderCallCount, "expected func to be called exactly once")
		assert.Equal(t, "failed to load config", err.Error())
		assert.Nil(t, ctx.Config(), "expected context not to have config set")
	})
}

var FailToCreateConfigObject = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		configLoaderCallCount := 0
		fakeCfgMgr := config.CreateFakeConfigManager()
		fakeCfgMgr.SetupConfigFileLoaderMock = func() error {
			configLoaderCallCount++
			return nil
		}
		configListenChangesCallCount := 0
		fakeCfgMgr.ListenOnConfigFileChangesMock = func(ctx common.Context) {
			configListenChangesCallCount++
		}
		createConfigCallCount := 0
		fakeCfgMgr.CreateConfigObjectMock = func() (*config.AnchorConfig, error) {
			createConfigCallCount++
			return nil, fmt.Errorf("failed to create config object")
		}
		err := initConfiguration(ctx, fakeCfgMgr)
		assert.NotNil(t, err)
		assert.Equal(t, 1, configLoaderCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, configListenChangesCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, createConfigCallCount, "expected func to be called exactly once")
		assert.Equal(t, "failed to create config object", err.Error())
		assert.Nil(t, ctx.Config(), "expected context not to have config set")
	})
}

var InitializeRegistrySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		err := initRegistry(ctx)
		assert.Nil(t, err)

		reg := ctx.Registry()

		r, _ := reg.SafeGet(locator.Identifier)
		assert.NotNil(t, r, "expected item from registry to exist. name: locator")

		s, _ := reg.SafeGet(shell.Identifier)
		assert.NotNil(t, s, "expected item from registry to exist. name: shell")

		e, _ := reg.SafeGet(extractor.Identifier)
		assert.NotNil(t, e, "expected item from registry to exist. name: extractor")

		pa, _ := reg.SafeGet(parser.Identifier)
		assert.NotNil(t, pa, "expected item from registry to exist. name: parser")

		pr, _ := reg.SafeGet(prompter.Identifier)
		assert.NotNil(t, pr, "expected item from registry to exist. name: prompter")

		prntr, _ := reg.SafeGet(printer.Identifier)
		assert.NotNil(t, prntr, "expected item from registry to exist. name: printer")

		in, _ := reg.SafeGet(input.Identifier)
		assert.NotNil(t, in, "expected item from registry to exist. name: input")
	})
}

var PreRunSequenceDoNothingWhenNoArgs = func(t *testing.T) {
	shouldStartPreRunSeq := isPreRunSequenceExcludedCommand(nil)
	assert.False(t, shouldStartPreRunSeq)
}

var PreRunSequenceRunForRootCommand = func(t *testing.T) {
	shouldStartPreRunSeq := isPreRunSequenceExcludedCommand([]string{"anchor"})
	assert.False(t, shouldStartPreRunSeq)
}

var PreRunSequenceDoNotRunForExcludedCommand = func(t *testing.T) {
	shouldStartPreRunSeq := isPreRunSequenceExcludedCommand([]string{"anchor", "config"})
	assert.True(t, shouldStartPreRunSeq)
}

var PreRunSequenceRunForNonExcludedCommand = func(t *testing.T) {
	shouldStartPreRunSeq := isPreRunSequenceExcludedCommand([]string{"anchor", "some_command"})
	assert.False(t, shouldStartPreRunSeq)
}
