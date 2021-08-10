package main

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/internal/registry"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"os"
)

type MainCollaborators struct {
	Logger           func(ctx common.Context) error
	Configuration    func(ctx common.Context) error
	Registry         func(ctx common.Context) error
	StartCliCommands func(ctx common.Context) error
}

var collaborators = &MainCollaborators{
	Logger: func(ctx common.Context) error {
		return initLogger(ctx, logger.GetDefaultLoggerLogFilePath, logger.LogrusLoggerLoader)
	},
	Configuration: func(ctx common.Context) error {
		cfgManager := config.NewManager()
		ctx.Registry().Set(config.Identifier, cfgManager)
		return initConfiguration(ctx, cfgManager)
	},
	Registry: func(ctx common.Context) error {
		return initRegistry(ctx)
	},
	StartCliCommands: func(ctx common.Context) error {
		return startCliCommands(ctx)
	},
}

var exitApplication = func(code int, message string) {
	fmt.Printf(message)
	os.Exit(code)
}

func initLogger(
	ctx common.Context,
	logFileResolver func() (string, error),
	loggerCreator func(verbose bool, logFilePath string) (logger.Logger, error)) error {

	logFilePath, err := logFileResolver()
	if err != nil {
		return fmt.Errorf("failed to resolve logger file path. error: %s", err)
	}

	if lgr, err := loggerCreator(false, logFilePath); err != nil {
		return fmt.Errorf("failed to initialize logger. error: %s", err.Error())
	} else {
		ctx.(common.LoggerSetter).SetLogger(lgr)
	}
	return nil
}

func initConfiguration(ctx common.Context, cfgManager config.ConfigManager) error {
	err := cfgManager.SetupConfigFileLoader()
	if err != nil {
		return err
	}

	cfgManager.ListenOnConfigFileChanges(ctx)

	cfg, err := cfgManager.CreateConfigObject()
	if err != nil {
		return err
	}

	config.SetInContext(ctx, cfg)
	return nil
}

func initRegistry(ctx common.Context) error {
	reg := ctx.Registry()

	l := locator.New()
	reg.Set(locator.Identifier, l)

	s := shell.New()
	reg.Set(shell.Identifier, s)

	e := extractor.New()
	reg.Set(extractor.Identifier, e)

	pa := parser.New()
	reg.Set(parser.Identifier, pa)

	pr := prompter.New()
	reg.Set(prompter.Identifier, pr)

	prntr := printer.New()
	reg.Set(printer.Identifier, prntr)

	in := input.New()
	reg.Set(input.Identifier, in)

	o := orchestrator.New(pr, l, e, pa, s, in)
	reg.Set(orchestrator.Identifier, o)

	//registry.Initialize().Clipboard = clipboard.NewManager(registry.Initialize().shell)
	return nil
}

func startCliCommands(ctx common.Context) error {
	return anchor.RunCliRootCommand(ctx)
}

func runCollaboratorsInSequence(ctx common.Context, collaborators *MainCollaborators) error {
	err := collaborators.Logger(ctx)
	if err != nil {
		return err
	}
	err = collaborators.Configuration(ctx)
	if err != nil {
		return err
	}
	err = collaborators.Registry(ctx)
	if err != nil {
		return err
	}
	err = collaborators.StartCliCommands(ctx)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := common.EmptyAnchorContext(registry.Initialize())
	err := runCollaboratorsInSequence(ctx, collaborators)
	if err != nil {
		exitApplication(1, err.Error())
	}
}
