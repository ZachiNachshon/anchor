package _select

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type selectCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, selectFunc AppSelectFunc) (*selectCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:   "select",
		Short: "Select an application",
		Long:  `Select an application`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return selectFunc(ctx)
		},
	}

	return &selectCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}, nil
}

func (cmd *selectCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *selectCmd) InitFlags() error {
	return nil
}

func (cmd *selectCmd) InitSubCommands() error {
	return nil
}
