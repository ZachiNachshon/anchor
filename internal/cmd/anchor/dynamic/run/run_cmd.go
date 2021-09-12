package run

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type runCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

type NewCommandFunc func(ctx common.Context, parentFolderName string, runFunc DynamicRunFunc) *runCmd

func NewCommand(ctx common.Context, parentFolderName string, runFunc DynamicRunFunc) *runCmd {
	var cobraCmd = &cobra.Command{
		Use:   "run",
		Short: "Run an action without selection",
		Long:  `Run an action without selection`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFunc(ctx, NewOrchestrator(parentFolderName))
		},
	}

	return &runCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (c *runCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *runCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, parentFolderName string, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), parentFolderName, DynamicRun)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
