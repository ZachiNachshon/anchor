package anchor

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/spf13/cobra"
)

func RunCliRootCommand(ctx common.Context, shouldStartPreRunSeq bool) error {
	if result, err := ctx.Registry().SafeGet(logger.Identifier); err != nil {
		return err
	} else {
		loggerManager := result.(logger.LoggerManager)
		c := NewCommand(ctx, loggerManager)
		err = c.initialize(shouldStartPreRunSeq)
		if err != nil {
			return err
		}
		return c.GetCobraCmd().Execute()
	}
}

type anchorCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context

	initFlagsFunc func(o *anchorCmd) error

	startPreRunSequence       func(parent cmd.AnchorCommand, preRunSequence func(ctx common.Context) error) error
	addDynamicSubCommandsFunc func(parent cmd.AnchorCommand, createCmd dynamic.NewCommandsFunc) error
	addConfigSubCmdFunc       func(parent cmd.AnchorCommand, createCmd config_cmd.NewCommandFunc) error
	addCompletionSubCmdFunc   func(root cmd.AnchorCommand, createCmd completion.NewCommandFunc) error
	addVersionSubCmdFunc      func(parent cmd.AnchorCommand, createCmd version.NewCommandFunc) error
}

var validArgs = []string{"completion", "config", "version"}

var verboseFlagValue = false

func NewCommand(ctx common.Context, loggerManager logger.LoggerManager) *anchorCmd {
	var rootCmd = &cobra.Command{
		Use:           "anchor",
		Short:         "Anchor your Ops environment into a version controlled repository",
		Long:          `Anchor your Ops environment into a version controlled repository`,
		ValidArgs:     validArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			level := "info"
			if verboseFlagValue {
				level = "debug"
			}
			return loggerManager.SetVerbosityLevel(level)
		},
	}

	return &anchorCmd{
		cobraCmd:      rootCmd,
		ctx:           ctx,
		initFlagsFunc: initFlags,
		startPreRunSequence: func(parent cmd.AnchorCommand, preRunSequence func(ctx common.Context) error) error {
			return preRunSequence(parent.GetContext())
		},
		addDynamicSubCommandsFunc: dynamic.AddCommands,
		addConfigSubCmdFunc:       config_cmd.AddCommand,
		addCompletionSubCmdFunc:   completion.AddCommand,
		addVersionSubCmdFunc:      version.AddCommand,
	}
}

func initFlags(c *anchorCmd) error {
	c.cobraCmd.PersistentFlags().BoolVarP(
		&verboseFlagValue,
		globals.VerboseFlagName,
		"v",
		verboseFlagValue,
		"anchor <command> -v")

	c.cobraCmd.PersistentFlags().SortFlags = false
	return nil
}

func (c *anchorCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *anchorCmd) GetContext() common.Context {
	return c.ctx
}

func (c *anchorCmd) initialize(shouldStartPreRunSeq bool) error {
	// Cannot run on the command Run() method itself since we must initialize the logger
	// logger must be available at the PersistentPreRun() stage
	err := c.initFlagsFunc(c)
	if err != nil {
		return err
	}
	//cobra.EnableCommandSorting = false

	// Pre Run Sequence
	if shouldStartPreRunSeq {
		logger.Infof("starting pre run sequence for command...")
		err = c.startPreRunSequence(c, NewAnchorCollaborators().Run)
		if err != nil {
			return err
		}
	} else {
		logger.Infof("excluded command identified, skipping pre run sequence")
	}

	// Dynamic Commands
	err = c.addDynamicSubCommandsFunc(c, dynamic.NewCommands)
	if err != nil {
		return err
	}

	// Config Commands
	err = c.addConfigSubCmdFunc(c, config_cmd.NewCommand)
	if err != nil {
		return err
	}

	// Version
	err = c.addVersionSubCmdFunc(c, version.NewCommand)
	if err != nil {
		return err
	}

	// Auto completion
	err = c.addCompletionSubCmdFunc(c, completion.NewCommand)
	if err != nil {
		return err
	}
	return nil
}
