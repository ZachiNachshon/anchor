package versions

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type versionsCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, cliVersions CliVersionsFunc) (*versionsCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:   "versions",
		Short: "Print versions of all CLI application",
		Long:  `Print versions of all CLI application`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cliVersions(ctx)
		},
	}

	return &versionsCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}, nil
}

func (cmd *versionsCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *versionsCmd) InitFlags() error {
	return nil
}

func (cmd *versionsCmd) InitSubCommands() error {
	return nil
}
