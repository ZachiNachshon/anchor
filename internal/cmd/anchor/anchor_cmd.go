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
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/spf13/cobra"
)

func RunCliRootCommand(ctx common.Context) error {
	if result, err := ctx.Registry().SafeGet(logger.Identifier); err != nil {
		return err
	} else {
		loggerManager := result.(logger.LoggerManager)
		if c, err := NewCommand(ctx, loggerManager); err != nil {
			return err
		} else {
			initRootCmd, err := c.Initialize()
			if err != nil {
				return err
			}
			return initRootCmd.GetCobraCmd().Execute()
		}
	}
}

type anchorCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"app", "cli", "completion", "config", "controller", "version"}

const verboseFlagName = "verbose"

var verboseFlagValue = false

func NewCommand(ctx common.Context, loggerManager logger.LoggerManager) (*anchorCmd, error) {
	var rootCmd = &cobra.Command{
		Use:       "anchor",
		Short:     "Anchor your Ops environment into a version controlled repository",
		Long:      `Anchor your Ops environment into a version controlled repository`,
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			level := "info"
			if verboseFlagValue {
				level = "debug"
			}
			return loggerManager.SetVerbosityLevel(level)
		},
	}

	return &anchorCmd{
		cobraCmd: rootCmd,
		ctx:      ctx,
	}, nil
}

func (c *anchorCmd) InitFlags() error {
	c.cobraCmd.PersistentFlags().BoolVarP(
		&verboseFlagValue,
		verboseFlagName,
		"v",
		verboseFlagValue,
		"anchor <command> -v")

	c.cobraCmd.PersistentFlags().SortFlags = false
	return nil
}

func (c *anchorCmd) InitSubCommands() error {
	preRunSequence := AnchorPreRunSequence()

	//cobra.EnableCommandSorting = false

	// Apps Commands
	if appCmd, err := app.NewCommand(c.ctx, preRunSequence.Run); err != nil {
		return err
	} else {
		c.cobraCmd.AddCommand(appCmd.GetCobraCmd())
	}

	// CLI Commands
	if cliCmd, err := cli.NewCommand(c.ctx, preRunSequence.Run); err != nil {
		return err
	} else {
		c.cobraCmd.AddCommand(cliCmd.GetCobraCmd())
	}

	// Controller Commands
	if controllerCmd, err := controller.NewCommand(c.ctx, preRunSequence.Run); err != nil {
		return err
	} else {
		c.cobraCmd.AddCommand(controllerCmd.GetCobraCmd())
	}

	// Config Commands
	if cfgCmd, err := config_cmd.NewCommand(c.ctx); err != nil {
		return err
	} else {
		c.cobraCmd.AddCommand(cfgCmd.GetCobraCmd())
	}

	// Version
	if versionCmd, err := version.NewCommand(c.ctx, version.VersionVersion); err != nil {
		return err
	} else {
		c.cobraCmd.AddCommand(versionCmd.GetCobraCmd())
	}

	// Auto completion
	if compCmd, err := completion.NewCommand(c.cobraCmd, c.ctx); err != nil {
		return err
	} else {
		c.cobraCmd.AddCommand(compCmd.GetCobraCmd())
	}
	return nil
}

func (c *anchorCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *anchorCmd) Initialize() (*anchorCmd, error) {
	// Cannot run on the command Run() method itself since we must initialize the logger
	// logger must be available at the PersistentPreRun() stage
	err := c.InitFlags()
	if err != nil {
		return nil, err
	}
	err = c.InitSubCommands()
	if err != nil {
		return nil, err
	}
	return c, nil
}
