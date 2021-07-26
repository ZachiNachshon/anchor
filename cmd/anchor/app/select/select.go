package _select

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/app"
	"github.com/spf13/cobra"
)

type selectCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, appActions *app.ApplicationActions) *selectCmd {
	var cobraCmd = &cobra.Command{
		Use:   "select",
		Short: "Select an application",
		Long:  `Select an application`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return appActions.Select(ctx)
		},
	}

	return &selectCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *selectCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *selectCmd) InitFlags() {
}

func (cmd *selectCmd) InitSubCommands() {
}
