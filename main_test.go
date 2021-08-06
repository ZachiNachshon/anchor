package main

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
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
			Name: "fail to load configuration",
			Func: FailToLoadConfiguration,
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

var RunCollaboratorsInASpecificOrder = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var callOrder []string
		loggerCallCount := 0
		configCallCount := 0
		registryCallCount := 0
		startCallCount := 0
		collaborators := &MainCollaborators{
			Logger: func() {
				callOrder = append(callOrder, "logger")
				loggerCallCount++
			},
			Configuration: func(ctx common.Context) {
				callOrder = append(callOrder, "configuration")
				configCallCount++
			},
			Registry: func(ctx common.Context) {
				callOrder = append(callOrder, "registry")
				registryCallCount++
			},
			StartCliCommands: func(ctx common.Context) {
				callOrder = append(callOrder, "start")
				startCallCount++
			},
		}
		runCollaboratorsInSequence(ctx, collaborators)
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
	logFileResolverCallCount := 0
	logFileResolver := func() (string, error) {
		logFileResolverCallCount++
		return "/path/to/log/anchor.log", nil
	}
	loggerCreatorCallCount := 0
	loggerCreator := func(verbose bool, logFilePath string) error {
		loggerCreatorCallCount++
		return nil
	}

	initLogger(logFileResolver, loggerCreator)
	assert.Equal(t, 1, logFileResolverCallCount, "expected func to be called exactly once")
	assert.Equal(t, 1, loggerCreatorCallCount, "expected func to be called exactly once")
}

var FailToResolveLogFilePath = func(t *testing.T) {
	logFileResolverCallCount := 0
	logFileResolver := func() (string, error) {
		logFileResolverCallCount++
		return "", fmt.Errorf("failed to resolve")
	}
	loggerCreator := func(verbose bool, logFilePath string) error {
		return nil
	}
	exitCallCount := 0
	exitApplication = func(code int, message string) {
		exitCallCount++
	}
	initLogger(logFileResolver, loggerCreator)
	assert.Equal(t, 1, logFileResolverCallCount, "expected func to be called exactly once")
	assert.Equal(t, 1, exitCallCount, "expected exit to to be called exactly once")
}

var FailToCreateLogger = func(t *testing.T) {
	logFileResolverCallCount := 0
	logFileResolver := func() (string, error) {
		logFileResolverCallCount++
		return "/path/to/log/anchor.log", nil
	}
	loggerCreatorCallCount := 0
	loggerCreator := func(verbose bool, logFilePath string) error {
		loggerCreatorCallCount++
		return fmt.Errorf("failed to create")
	}
	exitCallCount := 0
	exitApplication = func(code int, message string) {
		exitCallCount++
	}
	initLogger(logFileResolver, loggerCreator)
	assert.Equal(t, 1, logFileResolverCallCount, "expected func to be called exactly once")
	assert.Equal(t, 1, loggerCreatorCallCount, "expected func to be called exactly once")
	assert.Equal(t, 1, exitCallCount, "expected exit to to be called exactly once")
}

var InitializeConfigurationSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		configInUse := config.AnchorConfig{}
		configLoaderCallCount := 0
		configLoader := func() (*config.AnchorConfig, error) {
			configLoaderCallCount++
			return &configInUse, nil
		}
		configChangesListenerCallCount := 0
		configChangesListener := func(ctx common.Context) {
			configChangesListenerCallCount++
		}

		initConfiguration(ctx, configLoader, configChangesListener)
		assert.Equal(t, 1, configLoaderCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, configChangesListenerCallCount, "expected func to be called exactly once")
		assert.Equal(t, configInUse, ctx.Config().(config.AnchorConfig))
	})
}

var FailToLoadConfiguration = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		configLoaderCallCount := 0
		configLoader := func() (*config.AnchorConfig, error) {
			configLoaderCallCount++
			return nil, fmt.Errorf("failed to load config")
		}
		configChangesListenerCallCount := 0
		configChangesListener := func(ctx common.Context) {
			configChangesListenerCallCount++
		}
		exitCallCount := 0
		exitApplication = func(code int, message string) {
			exitCallCount++
		}
		initConfiguration(ctx, configLoader, configChangesListener)
		assert.Equal(t, 1, configLoaderCallCount, "expected func to be called exactly once")
		assert.Equal(t, 0, configChangesListenerCallCount, "expected func not to be called")
		assert.Nil(t, ctx.Config(), "expected context not to have config set")
		assert.Equal(t, 1, exitCallCount, "expected exit to to be called exactly once")
	})
}

var InitializeRegistrySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		initRegistry(ctx)

		r, _ := locator.FromRegistry(ctx.Registry())
		assert.NotNil(t, r, "expected item from registry to exist. name: locator")

		s, _ := shell.FromRegistry(ctx.Registry())
		assert.NotNil(t, s, "expected item from registry to exist. name: shell")

		e, _ := extractor.FromRegistry(ctx.Registry())
		assert.NotNil(t, e, "expected item from registry to exist. name: extractor")

		pa, _ := parser.FromRegistry(ctx.Registry())
		assert.NotNil(t, pa, "expected item from registry to exist. name: parser")

		pr, _ := prompter.FromRegistry(ctx.Registry())
		assert.NotNil(t, pr, "expected item from registry to exist. name: prompter")

		prntr, _ := printer.FromRegistry(ctx.Registry())
		assert.NotNil(t, prntr, "expected item from registry to exist. name: printer")

		in, _ := input.FromRegistry(ctx.Registry())
		assert.NotNil(t, in, "expected item from registry to exist. name: input")

		o, _ := orchestrator.FromRegistry(ctx.Registry())
		assert.NotNil(t, o, "expected item from registry to exist. name: orchestrator")
	})
}
