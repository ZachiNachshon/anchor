package main

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/cmd/anchor"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/resolver"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func injectComponents(ctx common.Context) {
	l := locator.New(ctx.AnchorFilesPath())
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

	//registry.Initialize().Clipboard = clipboard.New(registry.Initialize().Shell)
}

func scanAnchorfilesRepositoryTree(ctx common.Context) {
	l, _ := locator.FromRegistry(ctx.Registry())
	err := l.Scan()
	if err != nil {
		logger.Fatalf("Failed to locate and scan anchorfiles repository content")
	}
}

func main() {
	ctx := common.EmptyAnchorContext(registry.Initialize())

	if err := logger.LogrusLoggerLoader(false); err != nil {
		fmt.Printf("Failed to initialize logger. error: %s", err.Error())
	}

	cfg, err := config.ViperConfigFileLoader()
	if err != nil {
		logger.Fatalf("Failed to load configuration. error: %s", err.Error())
		return
	}
	ctx.(common.ConfigSetter).SetConfig(*cfg)

	rslvr, err := resolver.GetResolverBasedOnConfig(cfg.Config.Repository)
	if err != nil {
		logger.Fatalf(err.Error())
		return
	}

	repoPath, err := rslvr.ResolveRepository(ctx)
	if err != nil {
		logger.Fatalf(err.Error())
		return
	}

	ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(repoPath)

	injectComponents(ctx)
	scanAnchorfilesRepositoryTree(ctx)
	anchor.Main(ctx)
}
