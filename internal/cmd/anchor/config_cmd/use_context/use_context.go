package use_context

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
)

type ConfigUseContextFunc func(
	ctx common.Context,
	cfgCtxName string,
	cfgManager config.ConfigManager) error

var ConfigUseContext = func(ctx common.Context, cfgCtxName string, cfgManager config.ConfigManager) error {
	cfg := config.FromContext(ctx)
	if cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, cfgCtxName); cfgCtx == nil {
		return fmt.Errorf("could not identify config context. name: %s", cfgCtxName)
	} else {
		err := cfgManager.OverrideConfigEntry("config.currentContext", cfgCtxName)
		if err != nil {
			return err
		}
		logger.Infof("Current config context set successfully. name: %s", cfgCtxName)
	}
	return nil
}
