package status

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context) *statusCmd {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: "Read status for installed application",
		Long:  `Read status for installed application`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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
