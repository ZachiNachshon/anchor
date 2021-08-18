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
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"os"
)

type MainCollaborators struct {
	Logger           func(ctx common.Context, loggerManager logger.LoggerManager) error
	Configuration    func(ctx common.Context, configManager config.ConfigManager) error
	Registry         func(ctx common.Context) error
	StartCliCommands func(ctx common.Context) error
}

var GetCollaborators = func() *MainCollaborators {
	return collaborators
}

var collaborators = &MainCollaborators{
	Logger: func(ctx common.Context, loggerManager logger.LoggerManager) error {
		ctx.Registry().Set(logger.Identifier, loggerManager)
		return initLogger(ctx, loggerManager)
	},
	Configuration: func(ctx common.Context, configManager config.ConfigManager) error {
		ctx.Registry().Set(config.Identifier, configManager)
		return initConfiguration(ctx, configManager)
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

func initLogger(ctx common.Context, logManager logger.LoggerManager) error {
	lgr, err := logManager.CreateEmptyLogger()
	if err != nil {
		return err
	}

	lgr, err = logManager.AppendStdoutLogger("info")
	if err != nil {
		return err
	}

	// TODO: add retention for xx log files with log rotation to conserve disk space
	//       currently file based logger use debug level for visibility
	lgr, err = logManager.AppendFileLogger("debug")
	if err != nil {
		return err
	}

	err = logManager.SetActiveLogger(&lgr)
	if err != nil {
		return err
	}

	logger.SetInContext(ctx, &lgr)
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

	//registry.Append().Clipboard = clipboard.NewManager(registry.Append().shell)
	return nil
}

func startCliCommands(ctx common.Context) error {
	return anchor.RunCliRootCommand(ctx)
}

func runCollaboratorsInSequence(ctx common.Context, collaborators *MainCollaborators) error {
	loggerManager := logger.NewManager()
	err := collaborators.Logger(ctx, loggerManager)
	if err != nil {
		return err
	}
	configManager := config.NewManager()
	err = collaborators.Configuration(ctx, configManager)
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
	err := runCollaboratorsInSequence(ctx, GetCollaborators())
	if err != nil {
		exitApplication(1, err.Error())
	}
}
