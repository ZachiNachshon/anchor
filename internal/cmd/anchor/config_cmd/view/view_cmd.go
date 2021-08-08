package view

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

type viewCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, viewFunc ConfigViewFunc) *viewCmd {
	var cobraCmd = &cobra.Command{
		Use:   "view",
		Short: "Display configuration file settings",
		Long:  `Display configuration file settings`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewFunc(ctx, config.GetConfigFilePath)
		},
	}

	return &viewCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *viewCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *viewCmd) InitFlags() {
}

func (cmd *viewCmd) InitSubCommands() {
}
