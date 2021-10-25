package use_context

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
)

type ConfigUseContextFunc func(ctx common.Context, o *UseContextOrchestrator) error

var ConfigUseContext = func(ctx common.Context, o *UseContextOrchestrator) error {
	return o.RunFunc(o, ctx)
}

type UseContextOrchestrator struct {
	cfgCtxName string
	cfgManager config.ConfigManager

	RunFunc func(o *UseContextOrchestrator, ctx common.Context) error
}

func NewOrchestrator(cfgManager config.ConfigManager, cfgCtxName string) *UseContextOrchestrator {
	return &UseContextOrchestrator{
		cfgManager: cfgManager,
		cfgCtxName: cfgCtxName,
		RunFunc:    run,
	}
}

func run(o *UseContextOrchestrator, ctx common.Context) error {
	cfg := config.FromContext(ctx)
	if cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, o.cfgCtxName); cfgCtx == nil {
		return fmt.Errorf("could not identify config context. name: %s", o.cfgCtxName)
	} else {
		err := o.cfgManager.OverrideConfigEntry("config.currentContext", o.cfgCtxName)
		if err != nil {
			return err
		}
		logger.Infof("Current config context set successfully. name: %s", o.cfgCtxName)
	}
	return nil
}
