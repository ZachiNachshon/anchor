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

func NewCommand(ctx common.Context, preRunSequence cmd.PreRunSequence) (*controllerCmd, error) {
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

	err := cmd.InitSubCommands()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (cmd *controllerCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *controllerCmd) InitFlags() error {
	return nil
}

func (cmd *controllerCmd) InitSubCommands() error {
	if installCmd, err := install.NewCommand(cmd.ctx, install.ControllerInstall); err != nil {
		return err
	} else {
		cmd.cobraCmd.AddCommand(installCmd.GetCobraCmd())
	}
	return nil
}
