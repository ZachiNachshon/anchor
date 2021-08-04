package view

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
	"github.com/spf13/cobra"
)

type viewCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, configActions *cfg.ConfigurationActions) *viewCmd {
	var cobraCmd = &cobra.Command{
		Use:   "view",
		Short: "Display configuration file settings",
		Long:  `Display configuration file settings`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return configActions.View(ctx)
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
