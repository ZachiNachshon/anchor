package install

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/controller"
	"github.com/spf13/cobra"
)

type installCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, ctlActions *controller.ControllerActions) *installCmd {
	var cobraCmd = &cobra.Command{
		Use:   "install",
		Short: "Install a controller",
		Long:  `Install a controller`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ctlActions.Install(ctx)
		},
	}

	return &installCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *installCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *installCmd) InitFlags() {
}

func (cmd *installCmd) InitSubCommands() {
}
