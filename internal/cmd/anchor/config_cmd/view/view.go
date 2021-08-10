package view

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

type ConfigViewFunc func(ctx common.Context, cfgManager config.ConfigManager) error

var ConfigView = func(ctx common.Context, cfgManager config.ConfigManager) error {
	cfg := config.FromContext(ctx)
	if cfgText, err := config.ConfigObjToYaml(cfg); err != nil {
		logger.Error(err.Error())
		return err
	} else {
		cfgFilePath, _ := cfgManager.GetConfigFilePath()
		return printConfiguration(ctx, cfgFilePath, cfgText)
	}
}

func printConfiguration(ctx common.Context, cfgFilePath string, cfgText string) error {
	var p printer.Printer
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		p = resolved.(printer.Printer)
		p.PrintConfiguration(cfgFilePath, cfgText)
	}
	return nil
}
