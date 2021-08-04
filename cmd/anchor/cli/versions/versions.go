package versions

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cli"
	"github.com/spf13/cobra"
)

type versionsCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, cliActions *cli.CliActions) *versionsCmd {
	var cobraCmd = &cobra.Command{
		Use:   "versions",
		Short: "Print versions of all CLI application",
		Long:  `Print versions of all CLI application`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cliActions.Versions(ctx)
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

func (cmd *versionsCmd) InitFlags() {
}

func (cmd *versionsCmd) InitSubCommands() {
}
