package controller

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/controller/install"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type controllerCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{""}

func NewCommand(ctx common.Context, preRunSequence cmd.PreRunSequence) *controllerCmd {
	var cobraCmd = &cobra.Command{
		Use:       "controller",
		Short:     "Kubernetes controllers commands",
		Aliases:   []string{"kc"},
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunSequence(ctx)
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
	cmd.cobraCmd.AddCommand(install.NewCommand(cmd.ctx, install.ControllerInstall).GetCobraCmd())
}
