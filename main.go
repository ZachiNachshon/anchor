package main

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func injectComponents(ctx common.Context) {
	locator.ToRegistry(ctx.Registry(), locator.New(ctx.AnchorFilesPath()))
	shell.ToRegistry(ctx.Registry(), shell.New())
	extractor.ToRegistry(ctx.Registry(), extractor.New())

	//registry.Initialize().CmdExtractor = extractor.New()
	//registry.Initialize().Clipboard = clipboard.New(registry.Initialize().Shell)
}

func main() {
	ctx := common.EmptyAnchorContext()

	if err := logger.LogrusLoggerLoader(false); err != nil {
		logger.Fatalf("Failed to setup logger. error: %s", err.Error())
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
	anchor.Main(ctx)
}
