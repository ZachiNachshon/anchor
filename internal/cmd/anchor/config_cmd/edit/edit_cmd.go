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

func NewCommand(ctx common.Context, editFunc ConfigEditFunc) *editCmd {
	var cobraCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		Long:  `Edit configuration file`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editFunc(ctx, config.GetConfigFilePath)
		},
	}

	return &editCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *editCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *editCmd) InitFlags() {
}

func (cmd *editCmd) InitSubCommands() {
}
