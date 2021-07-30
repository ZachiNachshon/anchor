package cfg

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func StartConfigEditFlow(ctx common.Context) error {
	cfgFilePath, _ := config.GetConfigFilePath()
	if s, err := shell.FromRegistry(ctx.Registry()); err != nil {
		return err
	} else {
		editScript := fmt.Sprintf("vi %s", cfgFilePath)
		err := s.ExecuteTTY(editScript)
		if err != nil {
			return err
		}
	}
	return nil
}
