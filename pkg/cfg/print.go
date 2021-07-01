package cfg

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
	"github.com/ZachiNachshon/anchor/pkg/utils/printer"
)

func StartConfigPrintFlow(ctx common.Context) error {
	cfg := ctx.Config().(config.AnchorConfig)
	if cfgText, err := converters.ConfigObjToYaml(cfg); err != nil {
		logger.Error(err.Error())
		return err
	} else {
		cfgFilePath := config.GetConfigFilePath()
		return printConfiguration(ctx, cfgFilePath, cfgText)
	}
}

func printConfiguration(ctx common.Context, cfgFilePath string, cfgText string) error {
	if p, err := printer.FromRegistry(ctx.Registry()); err != nil {
		return err
	} else {
		p.PrintConfiguration(ctx, cfgFilePath, cfgText)
	}
	return nil
}
