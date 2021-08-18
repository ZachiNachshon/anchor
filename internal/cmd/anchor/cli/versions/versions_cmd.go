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

type NewCommandFunc func(ctx common.Context, cliVersions CliVersionsFunc) *versionsCmd

func NewCommand(ctx common.Context, cliVersions CliVersionsFunc) *versionsCmd {
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
	}
}

func (c *versionsCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *versionsCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), CliVersions)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
