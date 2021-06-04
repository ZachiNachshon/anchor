package uninstall

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/spf13/cobra"
)

type uninstallCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context) *uninstallCmd {
	var cobraCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall an application",
		Long:  `Uninstall an application`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return &uninstallCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *uninstallCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}
