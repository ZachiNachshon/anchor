package set_context_entry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"strconv"
)

type ConfigSetContextEntryFunc func(
	ctx common.Context,
	cfgCtxName string,
	changes map[string]string,
	overrideConfigFunc func(cfgToUpdate *config.AnchorConfig) error) error

var ConfigSetContextEntry = func(
	ctx common.Context,
	cfgCtxName string,
	changes map[string]string,
	overrideConfigFunc func(cfgToUpdate *config.AnchorConfig) error) error {

	cfg := config.FromContext(ctx)
	if cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, cfgCtxName); cfgCtx == nil {
		return fmt.Errorf("could not identify config context. name: %s", cfgCtxName)
	} else {
		if err := populateConfigContextChanges(cfgCtx, changes); err != nil {
			return err
		}
		if err := overrideConfigFunc(cfg); err != nil {
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
		case remoteUrlFlagName:
			{
				cfgCtx.Context.Repository.Remote.Url = element
			}
		case remoteBranchFlagName:
			{
				cfgCtx.Context.Repository.Remote.Branch = element
			}
		case remoteRevisionFlagName:
			{
				cfgCtx.Context.Repository.Remote.Revision = element
			}
		case remoteClonePathFlagName:
			{
				cfgCtx.Context.Repository.Remote.ClonePath = element
			}
		case remoteAutoUpdateFlagName:
			{
				if parseBool, err := strconv.ParseBool(element); err != nil {
					return err
				} else {
					cfgCtx.Context.Repository.Remote.AutoUpdate = parseBool
				}
			}
		case localPathFlagName:
			{
				cfgCtx.Context.Repository.Local.Path = element
			}
		}
	}
	return nil
}
