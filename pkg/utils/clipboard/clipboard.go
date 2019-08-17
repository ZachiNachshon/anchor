package clipboard

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"runtime"
)

func Load(cmd string) error {
	switch runtime.GOOS {
	case "darwin":
		{
			if err := common.ShellExec.Execute(fmt.Sprintf("echo %v | pbcopy", cmd)); err != nil {
				logger.Info("Failed setting value to clipboard using 'pbcopy ...'")
				return err
			}
			break
		}
	case "linux":
		{
			if err := common.ShellExec.Execute(fmt.Sprintf("xclip -selection %v", cmd)); err != nil {
				logger.Info("Failed setting value to clipboard using 'xclip -selection ...'")
				return err
			}
			break
		}
	}
	return nil
}
