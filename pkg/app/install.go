package app

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/utils/banner"
	"github.com/manifoldco/promptui"
)

func StartApplicationInstallFlow(ctx common.Context) error {
	if o, err := orchestrator.FromRegistry(ctx.Registry()); err != nil {
		return err
	} else {
		banner.Print()
		if selection, err := o.OrchestrateAppInstructionSelection(); err != nil {
			if err == promptui.ErrInterrupt {
				logger.Debug("exit due to keyboard interrupt")
				return nil
			} else {
				logger.Error(err.Error())
				return err
			}
		} else {
			logger.Debugf("Selected: %v", selection.Id)
		}
	}
	return nil
}
