package cfg

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
)

func StartConfigUseContextFlow(ctx common.Context, cfgCtxName string) error {
	cfg := ctx.Config().(config.AnchorConfig)
	if cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, cfgCtxName); cfgCtx == nil {
		return fmt.Errorf("could not identify config context. name: %s", cfgCtxName)
	} else {
		err := config.OverrideConfigEntry("config.currentContext", cfgCtxName)
		if err != nil {
			return err
		}
		logger.Infof("Current config context set successfully. name: %s", cfgCtxName)
	}
	return nil
}
