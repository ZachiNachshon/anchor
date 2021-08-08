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
	"github.com/spf13/cobra"
)

func RunCliRootCommand(ctx common.Context) error {
	c := NewCommand(ctx, SetLoggerVerbosity).Initialize()
	return c.GetCobraCmd().Execute()
}

type anchorCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"app", "cli", "completion", "config", "controller", "version"}

const verboseFlagName = "verbose"

var verboseFlagValue = false

func NewCommand(ctx common.Context, setLoggerVerbosity cmd.SetLoggerVerbosityFunc) *anchorCmd {
	var rootCmd = &cobra.Command{
		Use:       "anchor",
		Short:     "Anchor your Ops environment into a version controlled repository",
		Long:      `Anchor your Ops environment into a version controlled repository`,
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return setLoggerVerbosity(ctx.Logger(), verboseFlagValue)
		},
	}

	return &anchorCmd{
		cobraCmd: rootCmd,
		ctx:      ctx,
	}
}

func (c *anchorCmd) InitFlags() {
	c.cobraCmd.PersistentFlags().BoolVarP(
		&verboseFlagValue,
		verboseFlagName,
		"v",
		verboseFlagValue,
		"anchor <command> -v")

	c.cobraCmd.PersistentFlags().SortFlags = false
}

func (c *anchorCmd) InitSubCommands() {

	//cobra.EnableCommandSorting = false

	// Apps Commands
	c.cobraCmd.AddCommand(app.NewCommand(c.ctx, LoadRepoOrFail).GetCobraCmd())

	// CLI Commands
	c.cobraCmd.AddCommand(cli.NewCommand(c.ctx, LoadRepoOrFail).GetCobraCmd())

	// Controller Commands
	c.cobraCmd.AddCommand(controller.NewCommand(c.ctx, LoadRepoOrFail).GetCobraCmd())

	// Config Commands
	c.cobraCmd.AddCommand(config_cmd.NewCommand(c.ctx).GetCobraCmd())

	// Version
	c.cobraCmd.AddCommand(version.NewCommand(c.ctx, version.VersionVersion).GetCobraCmd())

	// Auto completion
	c.cobraCmd.AddCommand(completion.NewCommand(c.cobraCmd, c.ctx).GetCobraCmd())
}

func (c *anchorCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *anchorCmd) Initialize() *anchorCmd {
	// Cannot run on the command Run() method itself since we must initialize the logger
	// logger must be available at the PersistentPreRun() stage
	c.InitFlags()
	c.InitSubCommands()
	return c
}
