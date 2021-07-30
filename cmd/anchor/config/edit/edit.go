package edit

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
	"github.com/spf13/cobra"
)

type editCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, configActions *cfg.ConfigurationActions) *editCmd {
	var cobraCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		Long:  `Edit configuration file`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return configActions.Edit(ctx)
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
