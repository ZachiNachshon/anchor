package docker

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
)

func PurgeAll() error {
	in := input.NewYesNoInput()
	if result, err := in.WaitForInput("Purge ALL docker images and containers?"); err != nil || !result {
		logger.Info("skipping.")
	} else {
		if err := common.ShellExec.Execute("docker system prune --all --force"); err != nil {
			return err
		}
	}
	return nil
}
