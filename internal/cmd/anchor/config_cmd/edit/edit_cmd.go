package edit

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

type editCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(
	ctx common.Context,
	cfgManager config.ConfigManager,
	editFunc ConfigEditFunc) (*editCmd, error) {

	var cobraCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		Long:  `Edit configuration file`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editFunc(ctx, cfgManager)
		},
	}

	return &editCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}, nil
}

func (cmd *editCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *editCmd) InitFlags() error {
	return nil
}

func (cmd *editCmd) InitSubCommands() error {
	return nil
}
