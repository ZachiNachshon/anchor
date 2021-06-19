package install

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/app"
	"github.com/spf13/cobra"
)

type installCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, appActions *app.ApplicationActions) *installCmd {
	var cobraCmd = &cobra.Command{
		Use:   "install",
		Short: "Install an application",
		Long:  `Install an application`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := appActions.Install(ctx)
			if err != nil {
				logger.Fatal(err.Error())
			}
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
