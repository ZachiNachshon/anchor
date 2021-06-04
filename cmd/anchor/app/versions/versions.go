package versions

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/spf13/cobra"
)

type versionsCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context) *versionsCmd {
	var cobraCmd = &cobra.Command{
		Use:   "versions",
		Short: "Print versions of all application installed components",
		Long:  `Print versions of all application installed components`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return &versionsCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *versionsCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}
