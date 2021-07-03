package main

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/cmd/anchor"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"os"
)

func injectComponents(ctx common.Context) {
	l := locator.New()
	locator.ToRegistry(ctx.Registry(), l)

	s := shell.New()
	shell.ToRegistry(ctx.Registry(), s)

	e := extractor.New()
	extractor.ToRegistry(ctx.Registry(), e)

	pa := parser.New()
	parser.ToRegistry(ctx.Registry(), pa)

	pr := prompter.New()
	prompter.ToRegistry(ctx.Registry(), pr)

	prntr := printer.New()
	printer.ToRegistry(ctx.Registry(), prntr)

	o := orchestrator.New(pr, l, e, pa)
	orchestrator.ToRegistry(ctx.Registry(), o)

	//registry.Initialize().Clipboard = clipboard.New(registry.Initialize().shell)
}

func main() {
	ctx := common.EmptyAnchorContext(registry.Initialize())

	logFilePath, err := config.GetDefaultLoggerLogFilePath()
	if err != nil {
		fmt.Printf("failed to resolve logger file path. error: %s", err)
		os.Exit(1)
	}

	if err = logger.LogrusLoggerLoader(false, logFilePath); err != nil {
		fmt.Printf("Failed to initialize logger. error: %s", err.Error())
	}

	cfg, err := config.ViperConfigFileLoader()
	if err != nil {
		logger.Fatalf("Failed to load configuration. error: %s", err.Error())
		return
	}
	ctx.(common.ConfigSetter).SetConfig(*cfg)

	injectComponents(ctx)
	anchor.Main(ctx)
}
