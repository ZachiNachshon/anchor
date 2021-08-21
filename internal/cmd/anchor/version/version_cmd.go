package version

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

type NewCommandFunc func(ctx common.Context, versionFunc VersionVersionFunc) *versionCmd

func NewCommand(ctx common.Context, versionFunc VersionVersionFunc) *versionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "version",
		Short: "Print anchor CLI version",
		Long:  `Print anchor CLI version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return versionFunc(ctx, NewOrchestrator())
		},
	}

	return &versionCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (c *versionCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *versionCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), VersionVersion)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
