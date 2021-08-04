package status

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/app"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, appActions *app.ApplicationActions) *statusCmd {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: "Check status validity of supported applications",
		Long:  `Check status validity of supported applications`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return appActions.Status(ctx)
		},
	}

	return &statusCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *statusCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *statusCmd) InitFlags() {
}

func (cmd *statusCmd) InitSubCommands() {
}
