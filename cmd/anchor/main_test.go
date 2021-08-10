package main

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MainShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "set valid initializations to collaborators",
			Func: SetValidInitializationsToCollaborators,
		},
		{
			Name: "init logger collaborator successfully",
			Func: InitLoggerCollaboratorSuccessfully,
		},
		{
			Name: "init configuration collaborator successfully",
			Func: InitConfigurationCollaboratorSuccessfully,
		},
		{
			Name: "init registry collaborator successfully",
			Func: InitRegistryCollaboratorSuccessfully,
		},
		{
			Name: "init start cli commands collaborator successfully",
			Func: InitStartCliCommandsCollaboratorSuccessfully,
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
			Name: "initialize logger successfully",
			Func: InitializeLoggerSuccessfully,
		},
		{
			Name: "fail to resolve log file path",
			Func: FailToResolveLogFilePath,
		},
		{
			Name: "fail to create logger",
			Func: FailToCreateLogger,
		},
		{
			Name: "fail to load log file path",
			Func: FailToResolveLogFilePath,
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
	}
	harness.RunTests(t, tests)
}

var SetValidInitializationsToCollaborators = func(t *testing.T) {
	assert.NotNil(t, collaborators.Logger, "expected collaborator not to be empty. name: Logger")
	assert.NotNil(t, collaborators.Configuration, "expected collaborator not to be empty. name: Configuration")
	assert.NotNil(t, collaborators.Registry, "expected collaborator not to be empty. name: Registry")
	assert.NotNil(t, collaborators.StartCliCommands, "expected collaborator not to be empty. name: StartCliCommands")
}

var InitLoggerCollaboratorSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		err := collaborators.Logger(ctx)
		assert.Nil(t, err)
	})
}

var InitConfigurationCollaboratorSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		err := collaborators.Configuration(ctx)
		assert.Nil(t, err)
	})
}

var InitRegistryCollaboratorSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		err := collaborators.Registry(ctx)
		assert.Nil(t, err)
	})
}

var InitStartCliCommandsCollaboratorSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		err := collaborators.StartCliCommands(ctx)
		assert.Nil(t, err)
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
		err := startCliCommands(ctx)
		assert.Nil(t, err)
	})
}

var StartMainEntryPointSuccessfully = func(t *testing.T) {
	main()
}

var RunCollaboratorsInASpecificOrder = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var callOrder []string
		loggerCallCount := 0
		configCallCount := 0
		registryCallCount := 0
		startCallCount := 0
		collaborators := &MainCollaborators{
			Logger: func(ctx common.Context) error {
				callOrder = append(callOrder, "logger")
				loggerCallCount++
				return nil
			},
			Configuration: func(ctx common.Context) error {
				callOrder = append(callOrder, "configuration")
				configCallCount++
				return nil
			},
			Registry: func(ctx common.Context) error {
				callOrder = append(callOrder, "registry")
				registryCallCount++
				return nil
			},
			StartCliCommands: func(ctx common.Context) error {
				callOrder = append(callOrder, "start")
				startCallCount++
				return nil
			},
		}
		err := runCollaboratorsInSequence(ctx, collaborators)
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

var InitializeLoggerSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		logFileResolverCallCount := 0
		logFileResolver := func() (string, error) {
			logFileResolverCallCount++
			return "/path/to/log/anchor.log", nil
		}
		loggerCreatorCallCount := 0
		loggerCreator := func(verbose bool, logFilePath string) (logger.Logger, error) {
			loggerCreatorCallCount++
			return nil, nil
		}

		err := initLogger(ctx, logFileResolver, loggerCreator)
		assert.Nil(t, err)
		assert.Equal(t, 1, logFileResolverCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, loggerCreatorCallCount, "expected func to be called exactly once")
	})
}

var FailToResolveLogFilePath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		logFileResolverCallCount := 0
		logFileResolver := func() (string, error) {
			logFileResolverCallCount++
			return "", fmt.Errorf("failed to resolve")
		}
		loggerCreator := func(verbose bool, logFilePath string) (logger.Logger, error) {
			return nil, nil
		}
		err := initLogger(ctx, logFileResolver, loggerCreator)
		assert.NotNil(t, err)
		assert.Equal(t, 1, logFileResolverCallCount, "expected func to be called exactly once")
	})
}

var FailToCreateLogger = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		logFileResolverCallCount := 0
		logFileResolver := func() (string, error) {
			logFileResolverCallCount++
			return "/path/to/log/anchor.log", nil
		}
		loggerCreatorCallCount := 0
		loggerCreator := func(verbose bool, logFilePath string) (logger.Logger, error) {
			loggerCreatorCallCount++
			return nil, fmt.Errorf("failed to create")
		}
		err := initLogger(ctx, logFileResolver, loggerCreator)
		assert.NotNil(t, err)
		assert.Equal(t, 1, logFileResolverCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, loggerCreatorCallCount, "expected func to be called exactly once")
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

		o, _ := reg.SafeGet(orchestrator.Identifier)
		assert.NotNil(t, o, "expected item from registry to exist. name: orchestrator")
	})
}
