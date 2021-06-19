package app

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
)

func StartApplicationInstallFlow(ctx common.Context) error {
	if o, err := orchestrator.FromRegistry(ctx.Registry()); err != nil {
		return err
	} else {
		if selection, err := o.OrchestrateAppInstructionSelection(); err != nil {
			logger.Error(err.Error())
			return err
		} else {
			logger.Infof("Selected: %v", selection.Id)
		}
	}
	return nil
}
