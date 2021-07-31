package cfg

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/resolver"
	"github.com/ZachiNachshon/anchor/logger"
	"strconv"
)

func StartConfigSetContextEntryFlow(ctx common.Context, cfgCtxName string, changes map[string]string) error {
	cfg := ctx.Config().(config.AnchorConfig)
	if cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, cfgCtxName); cfgCtx == nil {
		return fmt.Errorf("could not identify config context. name: %s", cfgCtxName)
	} else {
		if err := populateConfigContextChanges(cfgCtx, changes); err != nil {
			return err
		}
		if err := config.OverrideConfig(cfg); err != nil {
			return err
		} else {
			logger.Infof("Updated config context entries successfully. context: %s", cfgCtxName)
			return nil
		}
	}
}

func populateConfigContextChanges(cfgCtx *config.Context, changes map[string]string) error {
	for key, element := range changes {
		switch key {
		case resolver.RemoteUrlFlagName:
			{
				cfgCtx.Context.Repository.Remote.Url = element
			}
		case resolver.RemoteBranchFlagName:
			{
				cfgCtx.Context.Repository.Remote.Branch = element
			}
		case resolver.RemoteRevisionFlagName:
			{
				cfgCtx.Context.Repository.Remote.Revision = element
			}
		case resolver.RemoteClonePathFlagName:
			{
				cfgCtx.Context.Repository.Remote.ClonePath = element
			}
		case resolver.RemoteAutoUpdateFlagName:
			{
				if parseBool, err := strconv.ParseBool(element); err != nil {
					return err
				} else {
					cfgCtx.Context.Repository.Remote.AutoUpdate = parseBool
				}
			}
		case resolver.LocalPathFlagName:
			{
				cfgCtx.Context.Repository.Local.Path = element
			}
		}
	}
	return nil
}
