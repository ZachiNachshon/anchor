package _select

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/spf13/cobra"
	"strconv"
)

type selectCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

type NewCommandFunc func(ctx common.Context, selectFunc AppSelectFunc) *selectCmd

func NewCommand(ctx common.Context, selectFunc AppSelectFunc) *selectCmd {
	var cobraCmd = &cobra.Command{
		Use:   "select",
		Short: "Select an application",
		Long:  `Select an application`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			verboseFlag := cmd.Flag(globals.VerboseFlagName)
			orchestrator := NewOrchestrator()
			if verboseFlag != nil {
				if isVerbose, err := strconv.ParseBool(verboseFlag.Value.String()); err == nil {
					orchestrator.verbose = isVerbose
				}
			}
			return selectFunc(ctx, orchestrator)
		},
	}

	return &selectCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (c *selectCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *selectCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), AppSelect)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
