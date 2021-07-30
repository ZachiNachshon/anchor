package list

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/app"
	"github.com/spf13/cobra"
)

type listCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, appActions *app.ApplicationActions) *listCmd {
	var cobraCmd = &cobra.Command{
		Use:   "list",
		Short: "List all supported applications",
		Long:  `List all supported applications`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return appActions.List(ctx)
		},
	}

	return &listCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *listCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *listCmd) InitFlags() {
}

func (cmd *listCmd) InitSubCommands() {
}
