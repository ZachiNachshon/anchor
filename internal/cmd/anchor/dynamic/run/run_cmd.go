package run

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/runner"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/spf13/cobra"
	"strconv"
)

const (
	actionNameFlagName   = "action"
	workflowNameFlagName = "workflow"
)

var actionNameFlagValue = ""
var workflowNameFlagValue = ""

type runCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context

	initFlagsFunc func(o *runCmd) error
}

type NewCommandFunc func(ctx common.Context, parentFolderName string, runFunc DynamicRunFunc) *runCmd

func NewCommand(ctx common.Context, commandFolderName string, runFunc DynamicRunFunc) *runCmd {
	var cobraCmd = &cobra.Command{
		Use: fmt.Sprintf(`run [COMMAND_ITEM_NAME] [--%s=ACTION-ID or --%s=WORKFLOW-ID]`,
			actionNameFlagName,
			workflowNameFlagName,
		),
		Short:                 "Run an action without interactive selection",
		Long:                  `Run an action without interactive selection`,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dynamicCommandName := args[0]
			if len(actionNameFlagValue) > 0 && len(workflowNameFlagValue) > 0 {
				return fmt.Errorf("--%s and --%s flags are mutual exclusive flags", actionNameFlagName, workflowNameFlagName)
			} else if len(actionNameFlagValue) == 0 && len(workflowNameFlagValue) == 0 {
				return fmt.Errorf("must use either --%s or --%s flags, missing mandatory flag(s)", actionNameFlagName, workflowNameFlagName)
			}

			rnr := runner.NewOrchestrator(commandFolderName)
			o := NewOrchestrator(rnr, commandFolderName, dynamicCommandName)

			verboseFlag := cmd.Flag(globals.VerboseFlagName)
			if verboseFlag != nil {
				if isVerbose, err := strconv.ParseBool(verboseFlag.Value.String()); err == nil {
					o.verboseFlag = isVerbose
				}
			}

			var identifier = ""
			if len(actionNameFlagValue) > 0 {
				identifier = actionNameFlagValue
				o.activeRunFunc = o.runActionFunc
			} else if len(workflowNameFlagValue) > 0 {
				identifier = workflowNameFlagValue
				o.activeRunFunc = o.runWorkflowFunc
			}

			return runFunc(ctx, o, identifier)
		},
	}

	return &runCmd{
		cobraCmd:      cobraCmd,
		ctx:           ctx,
		initFlagsFunc: initFlags,
	}
}

func (c *runCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *runCmd) GetContext() common.Context {
	return c.ctx
}

func initFlags(c *runCmd) error {
	c.cobraCmd.Flags().StringVar(
		&actionNameFlagValue,
		actionNameFlagName,
		"",
		fmt.Sprintf("--%s=hello-world", actionNameFlagName))

	c.cobraCmd.Flags().StringVar(
		&workflowNameFlagValue,
		workflowNameFlagName,
		"",
		fmt.Sprintf("--%s=do-it-bigtime", workflowNameFlagName))

	c.cobraCmd.Flags().SortFlags = false
	return nil
}

func AddCommand(parent cmd.AnchorCommand, commandFolderName string, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), commandFolderName, DynamicRun)
	err := newCmd.initFlagsFunc(newCmd)
	if err != nil {
		return err
	}
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
