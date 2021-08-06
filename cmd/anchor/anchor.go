package anchor

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/app"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cli"
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion"
	configCmd "github.com/ZachiNachshon/anchor/cmd/anchor/config"
	"github.com/ZachiNachshon/anchor/cmd/anchor/controller"
	"github.com/ZachiNachshon/anchor/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/root"
	version_pkg "github.com/ZachiNachshon/anchor/pkg/version"
	"github.com/spf13/cobra"
)

type anchorCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"app", "cli", "completion", "config", "controller", "version"}

const verboseFlagName = "verbose"

var verboseFlagValue = false

func NewCommand(ctx common.Context) *anchorCmd {
	var rootCmd = &cobra.Command{
		Use:       "anchor",
		Short:     "Anchor your Ops environment into a version controlled repository",
		Long:      `Anchor your Ops environment into a version controlled repository`,
		ValidArgs: validArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := alignLoggerWithVerboseFlag(); err != nil {
				logger.Fatal(err.Error())
			}
		},
	}

	return &anchorCmd{
		cobraCmd: rootCmd,
		ctx:      ctx,
	}
}

func (cmd *anchorCmd) InitFlags() {
	cmd.cobraCmd.PersistentFlags().BoolVarP(
		&verboseFlagValue,
		verboseFlagName,
		"v",
		verboseFlagValue,
		"anchor <command> -v")

	cmd.cobraCmd.PersistentFlags().SortFlags = false
}

func (cmd *anchorCmd) InitSubCommands() {

	//cobra.EnableCommandSorting = false

	rootActions := root.DefineRootCommandActions()

	// Apps Commands
	cmd.cobraCmd.AddCommand(app.NewCommand(cmd.ctx, rootActions).GetCobraCmd())

	// CLI Commands
	cmd.cobraCmd.AddCommand(cli.NewCommand(cmd.ctx, rootActions).GetCobraCmd())

	// Controller Commands
	cmd.cobraCmd.AddCommand(controller.NewCommand(cmd.ctx, rootActions).GetCobraCmd())

	// Config Commands
	cmd.cobraCmd.AddCommand(configCmd.NewCommand(cmd.ctx).GetCobraCmd())

	// Version
	cmd.cobraCmd.AddCommand(version.NewCommand(cmd.ctx, version_pkg.DefineVersionActions()).GetCobraCmd())

	// Auto completion
	cmd.cobraCmd.AddCommand(completion.NewCommand(cmd.cobraCmd, cmd.ctx).GetCobraCmd())
}

func (cmd *anchorCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *anchorCmd) Initialize() *anchorCmd {
	// Cannot run on the command Run() method itself since we must initialize the logger
	// logger must be available at the PersistentPreRun() stage
	cmd.InitFlags()
	cmd.InitSubCommands()
	return cmd
}

func alignLoggerWithVerboseFlag() error {
	level := "info"
	if verboseFlagValue {
		level = "debug"
	}
	if err := logger.SetVerbosityLevel(level); err != nil {
		return err
	}
	return nil
}

func Main(ctx common.Context) {
	cmd := NewCommand(ctx).Initialize()
	if err := cmd.cobraCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}
