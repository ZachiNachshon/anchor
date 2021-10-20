package set_context_entry

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"strconv"
)

type ConfigSetContextEntryFunc func(ctx common.Context, o *setContextEntryOrchestrator) error

var ConfigSetContextEntry = func(ctx common.Context, o *setContextEntryOrchestrator) error {
	return o.runFunc(o, ctx)
}

type setContextEntryOrchestrator struct {
	cfgCtxName string
	cfgManager config.ConfigManager
	changes    map[string]string

	runFunc func(o *setContextEntryOrchestrator, ctx common.Context) error
}

func NewOrchestrator(cfgManager config.ConfigManager, cfgCtxName string, changes map[string]string) *setContextEntryOrchestrator {
	return &setContextEntryOrchestrator{
		cfgManager: cfgManager,
		cfgCtxName: cfgCtxName,
		changes:    changes,
		runFunc:    run,
	}
}

func run(o *setContextEntryOrchestrator, ctx common.Context) error {
	cfg := config.FromContext(ctx)
	cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, o.cfgCtxName)
	if cfgCtx == nil {
		cfgCtx = config.AppendEmptyConfigContext(cfg, o.cfgCtxName)
		if err := o.cfgManager.SetDefaultsPostCreation(cfg); err != nil {
			return err
		}
	}
	if err := populateConfigContextChanges(cfgCtx, o.changes); err != nil {
		return err
	}
	if err := o.cfgManager.OverrideConfig(cfg); err != nil {
		return err
	} else {
		logger.Infof("Updated config context entries successfully. context: %s", o.cfgCtxName)
		return nil
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
				if parsedBool, err := strconv.ParseBool(element); err != nil {
					return err
				} else {
					cfgCtx.Context.Repository.Remote.AutoUpdate = parsedBool
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
