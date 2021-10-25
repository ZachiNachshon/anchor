package set_context_entry

import (
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/use_context"
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
	cfgCtxName         string
	cfgManager         config.ConfigManager
	setAsCurrCfgCtx    bool
	changes            map[string]string
	useCtxOrchestrator *use_context.UseContextOrchestrator

	runFunc                     func(o *setContextEntryOrchestrator, ctx common.Context) error
	setCurrentConfigContextFunc func(ctx common.Context, useCtxOrchestrator *use_context.UseContextOrchestrator) error
}

func NewOrchestrator(cfgManager config.ConfigManager, cfgCtxName string, setAsCurrCfgCtx bool, changes map[string]string) *setContextEntryOrchestrator {
	return &setContextEntryOrchestrator{
		cfgManager:         cfgManager,
		cfgCtxName:         cfgCtxName,
		setAsCurrCfgCtx:    setAsCurrCfgCtx,
		changes:            changes,
		useCtxOrchestrator: use_context.NewOrchestrator(cfgManager, cfgCtxName),

		runFunc:                     run,
		setCurrentConfigContextFunc: setCurrentConfigContext,
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
		if o.setAsCurrCfgCtx {
			if err = o.setCurrentConfigContextFunc(ctx, o.useCtxOrchestrator); err != nil {
				return err
			}
		}
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

func setCurrentConfigContext(ctx common.Context, useCtxOrchestrator *use_context.UseContextOrchestrator) error {
	return use_context.ConfigUseContext(ctx, useCtxOrchestrator)
}
