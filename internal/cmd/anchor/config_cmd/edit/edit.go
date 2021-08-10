package edit

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

type ConfigEditFunc func(ctx common.Context, cfgManager config.ConfigManager) error

var ConfigEdit = func(ctx common.Context, cfgManager config.ConfigManager) error {
	cfgFilePath, _ := cfgManager.GetConfigFilePath()
	var s shell.Shell
	if resolved, err := ctx.Registry().SafeGet(shell.Identifier); err != nil {
		return err
	} else {
		s = resolved.(shell.Shell)
		editScript := fmt.Sprintf("vi %s", cfgFilePath)
		err := s.ExecuteTTY(editScript)
		if err != nil {
			return err
		}
	}
	return nil
}
