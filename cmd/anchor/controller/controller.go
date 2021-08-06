package controller

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/controller/install"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/controller"
	"github.com/ZachiNachshon/anchor/pkg/root"
	"github.com/spf13/cobra"
)

type controllerCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{""}

func NewCommand(ctx common.Context, rootActions *root.RootCommandActions) *controllerCmd {
	var cobraCmd = &cobra.Command{
		Use:       "controller",
		Short:     "Kubernetes controllers commands",
		Aliases:   []string{"kc"},
		ValidArgs: validArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			rootActions.LoadRepoOrFail(ctx)
		},
	}

	var cmd = &controllerCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	cmd.InitSubCommands()
	return cmd
}

func (cmd *controllerCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *controllerCmd) InitFlags() {
}

func (cmd *controllerCmd) InitSubCommands() {
	actions := controller.DefineControllerActions()
	cmd.cobraCmd.AddCommand(install.NewCommand(cmd.ctx, actions).GetCobraCmd())
}
