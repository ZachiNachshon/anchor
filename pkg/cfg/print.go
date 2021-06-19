package cfg

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
)

func StartConfigPrintFlow(ctx common.Context) error {
	cfg := ctx.Config().(config.AnchorConfig)
	if out, err := converters.ConfigObjToYaml(cfg); err != nil {
		logger.Error(err.Error())
		return err
	} else {
		fmt.Println(out)
	}
	return nil
}
