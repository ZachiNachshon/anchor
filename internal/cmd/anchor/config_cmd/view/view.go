package view

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/pkg/printer"
)

type ConfigViewFunc func(ctx common.Context, o *viewOrchestrator) error

var ConfigView = func(ctx common.Context, o *viewOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	return o.runFunc(o, ctx)
}

type viewOrchestrator struct {
	prntr      printer.Printer
	cfgManager config.ConfigManager

	prepareFunc func(o *viewOrchestrator, ctx common.Context) error
	runFunc     func(o *viewOrchestrator, ctx common.Context) error
}

func NewOrchestrator(cfgManager config.ConfigManager) *viewOrchestrator {
	return &viewOrchestrator{
		cfgManager:  cfgManager,
		prepareFunc: prepare,
		runFunc:     run,
	}
}

func prepare(o *viewOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}
	return nil
}

func run(o *viewOrchestrator, ctx common.Context) error {
	cfg := config.FromContext(ctx)
	if cfgText, err := config.ConfigObjToYaml(cfg); err != nil {
		return err
	} else {
		cfgFilePath, _ := o.cfgManager.GetConfigFilePath()
		o.prntr.PrintConfiguration(cfgFilePath, cfgText)
		return nil
	}
}
