package anchor

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/app"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/cli"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/controller"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/spf13/cobra"
)

func RunCliRootCommand(ctx common.Context) error {
	if result, err := ctx.Registry().SafeGet(logger.Identifier); err != nil {
		return err
	} else {
		loggerManager := result.(logger.LoggerManager)
		c := NewCommand(ctx, loggerManager)
		err = c.initialize()
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

	initFlagsFunc    func(o *anchorCmd) error
	addAppSubCmdFunc func(
		parent cmd.AnchorCommand,
		preRunSequence *cmd.AnchorCollaborators,
		createCmd app.NewCommandFunc) error

	addCliSubCmdFunc func(
		parent cmd.AnchorCommand,
		preRunSequence *cmd.AnchorCollaborators,
		createCmd cli.NewCommandFunc) error

	addControllerSubCmdFunc func(
		parent cmd.AnchorCommand,
		preRunSequence *cmd.AnchorCollaborators,
		createCmd controller.NewCommandFunc) error

	addConfigSubCmdFunc     func(parent cmd.AnchorCommand, createCmd config_cmd.NewCommandFunc) error
	addCompletionSubCmdFunc func(root cmd.AnchorCommand, createCmd completion.NewCommandFunc) error
	addVersionSubCmdFunc    func(parent cmd.AnchorCommand, createCmd version.NewCommandFunc) error
}

var validArgs = []string{"app", "cli", "completion", "config", "controller", "version"}

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
		cobraCmd:                rootCmd,
		ctx:                     ctx,
		initFlagsFunc:           initFlags,
		addAppSubCmdFunc:        app.AddCommand,
		addCliSubCmdFunc:        cli.AddCommand,
		addControllerSubCmdFunc: controller.AddCommand,
		addConfigSubCmdFunc:     config_cmd.AddCommand,
		addCompletionSubCmdFunc: completion.AddCommand,
		addVersionSubCmdFunc:    version.AddCommand,
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

func (c *anchorCmd) initialize() error {
	// Cannot run on the command Run() method itself since we must initialize the logger
	// logger must be available at the PersistentPreRun() stage
	err := c.initFlagsFunc(c)
	if err != nil {
		return err
	}
	//cobra.EnableCommandSorting = false

	preRunSequence := AnchorPreRunSequence()

	// Apps Commands
	err = c.addAppSubCmdFunc(c, preRunSequence, app.NewCommand)
	if err != nil {
		return err
	}

	// CLI Commands
	err = c.addCliSubCmdFunc(c, preRunSequence, cli.NewCommand)
	if err != nil {
		return err
	}

	// Controller Commands
	err = c.addControllerSubCmdFunc(c, preRunSequence, controller.NewCommand)
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
