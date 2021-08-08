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
	Logger           func(ctx common.Context)
	Configuration    func(ctx common.Context)
	Registry         func(ctx common.Context)
	StartCliCommands func(ctx common.Context)
}

var collaborators = &MainCollaborators{
	Logger: func(ctx common.Context) {
		initLogger(ctx, logger.GetDefaultLoggerLogFilePath, logger.LogrusLoggerLoader)
	},
	Configuration: func(ctx common.Context) {
		initConfiguration(ctx, config.ViperConfigFileLoader, config.ListenOnConfigFileChanges)
	},
	Registry: func(ctx common.Context) {
		initRegistry(ctx)
	},
	StartCliCommands: func(ctx common.Context) {
		startCliCommands(ctx)
	},
}

var exitApplication = func(code int, message string) {
	fmt.Printf(message)
	os.Exit(code)
}

func initLogger(
	ctx common.Context,
	logFileResolver func() (string, error),
	loggerCreator func(verbose bool, logFilePath string) (logger.Logger, error)) {

	logFilePath, err := logFileResolver()
	if err != nil {
		exitApplication(1, fmt.Sprintf("failed to resolve logger file path. error: %s", err))
	}

	if lgr, err := loggerCreator(false, logFilePath); err != nil {
		exitApplication(1, fmt.Sprintf("Failed to initialize logger. error: %s", err.Error()))
	} else {
		ctx.(common.LoggerSetter).SetLogger(lgr)
	}
}

func initConfiguration(
	ctx common.Context,
	configLoader func() (*config.AnchorConfig, error),
	configChangesListener func(ctx common.Context)) {

	cfg, err := configLoader()
	if err != nil {
		exitApplication(1, fmt.Sprintf("failed to load configuration. error: %s", err.Error()))
	} else {
		configChangesListener(ctx)
		ctx.(common.ConfigSetter).SetConfig(*cfg)
	}
}

func initRegistry(ctx common.Context) {
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

	//registry.Initialize().Clipboard = clipboard.New(registry.Initialize().shell)
}

func startCliCommands(ctx common.Context) {
	if err := anchor.RunCliRootCommand(ctx); err != nil {
		logger.Fatal(err.Error())
	}
}

func runCollaboratorsInSequence(ctx common.Context, collaborators *MainCollaborators) {
	collaborators.Logger(ctx)
	collaborators.Configuration(ctx)
	collaborators.Registry(ctx)
	collaborators.StartCliCommands(ctx)
}

func main() {
	ctx := common.EmptyAnchorContext(registry.Initialize())
	runCollaboratorsInSequence(ctx, collaborators)
}
