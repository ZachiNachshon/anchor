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

type NewCommandFunc func(ctx common.Context, parentFolderName string, selectFunc DynamicSelectFunc) *selectCmd

func NewCommand(ctx common.Context, parentFolderName string, selectFunc DynamicSelectFunc) *selectCmd {
	var cobraCmd = &cobra.Command{
		Use:   "select",
		Short: "Select an anchor folder item",
		Long:  `Select an anchor folder item`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			verboseFlag := cmd.Flag(globals.VerboseFlagName)
			orchestrator := NewOrchestrator(parentFolderName)
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

func AddCommand(parent cmd.AnchorCommand, parentFolderName string, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), parentFolderName, DynamicSelect)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
