package main

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/cmd/anchor"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func injectComponents(ctx common.Context) {
	locator.ToRegistry(ctx.Registry(), locator.New(ctx.AnchorFilesPath()))
	shell.ToRegistry(ctx.Registry(), shell.New())
	extractor.ToRegistry(ctx.Registry(), extractor.New())
	parser.ToRegistry(ctx.Registry(), parser.New())
	prompter.ToRegistry(ctx.Registry(), prompter.New())

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
	ctx := common.EmptyAnchorContext()

	if err := logger.LogrusLoggerLoader(false); err != nil {
		fmt.Printf("Failed to initialize logger. error: %s", err.Error())
	}

	if cfg, err := config.ViperConfigFileLoader(); err != nil {
		logger.Fatalf("Failed to load configuration. error: %s", err.Error())
	} else {
		ctx.(common.ConfigSetter).SetConfig(*cfg)
	}

	if path, err := config.ResolveAnchorfilesPathFromConfig(config.FromContext(ctx)); err != nil {
		logger.Fatal(err.Error())
	} else {
		ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(path)
	}

	injectComponents(ctx)
	scanAnchorfilesRepositoryTree(ctx)

	anchor.Main(ctx)
}
