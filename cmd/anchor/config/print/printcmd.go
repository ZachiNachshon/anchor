package printcmd

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
	"github.com/spf13/cobra"
)

type printCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, configActions *cfg.ConfigurationActions) *printCmd {
	var cobraCmd = &cobra.Command{
		Use:   "print",
		Short: "Print configuration content",
		Long:  `Print configuration content`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return configActions.Print(ctx)
		},
	}

	return &printCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *printCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *printCmd) InitFlags() {
}

func (cmd *printCmd) InitSubCommands() {
}
