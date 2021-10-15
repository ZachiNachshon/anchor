package _select

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/runner"
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

type NewCommandFunc func(ctx common.Context, commandFolderName string, selectFunc DynamicSelectFunc) *selectCmd

func NewCommand(ctx common.Context, commandFolderName string, selectFunc DynamicSelectFunc) *selectCmd {
	var cobraCmd = &cobra.Command{
		Use:   "select",
		Short: "Select an anchor folder item",
		Long:  `Select an anchor folder item`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			verboseFlag := cmd.Flag(globals.VerboseFlagName)
			r := runner.NewOrchestrator(commandFolderName)
			orchestrator := NewOrchestrator(r, commandFolderName)
			if verboseFlag != nil {
				if isVerbose, err := strconv.ParseBool(verboseFlag.Value.String()); err == nil {
					orchestrator.verboseFlag = isVerbose
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

func AddCommand(parent cmd.AnchorCommand, commandFolderName string, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), commandFolderName, DynamicSelect)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
