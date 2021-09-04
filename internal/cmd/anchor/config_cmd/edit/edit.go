package edit

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

type ConfigEditFunc func(ctx common.Context, o *editOrchestrator) error

var ConfigEdit = func(ctx common.Context, o *editOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	return o.runFunc(o, ctx)
}

type editOrchestrator struct {
	s          shell.Shell
	cfgManager config.ConfigManager

	prepareFunc func(o *editOrchestrator, ctx common.Context) error
	runFunc     func(o *editOrchestrator, ctx common.Context) error
}

func NewOrchestrator(cfgManager config.ConfigManager) *editOrchestrator {
	return &editOrchestrator{
		cfgManager:  cfgManager,
		prepareFunc: prepare,
		runFunc:     run,
	}
}

func prepare(o *editOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(shell.Identifier); err != nil {
		return err
	} else {
		o.s = resolved.(shell.Shell)
	}
	return nil
}

func run(o *editOrchestrator, ctx common.Context) error {
	cfgFilePath, _ := o.cfgManager.GetConfigFilePath()
	editScript := fmt.Sprintf("vi %s/config.yaml", cfgFilePath)
	return o.s.ExecuteTTY(editScript)
}
