package anchor

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/app"
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion"
	"github.com/ZachiNachshon/anchor/cmd/anchor/controller"
	"github.com/ZachiNachshon/anchor/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/spf13/cobra"
)

type anchorCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"app", "controller", "completion", "list", "version"}

func NewCommand(ctx common.Context) *anchorCmd {
	var rootCmd = &cobra.Command{
		Use:       "anchor",
		Short:     "Anchor your Ops environment into a version controlled repository",
		Long:      `Anchor your Ops environment into a version controlled repository`,
		ValidArgs: validArgs,
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		level := "info"
		if common.GlobalOptions.Verbose {
			level = "debug"
		}
		if err := logger.SetVerbosityLevel(level); err != nil {
			return err
		}
		return nil
	}

	return &anchorCmd{
		cobraCmd: rootCmd,
		ctx:      ctx,
	}
}

func (cmd *anchorCmd) InitFlags() {
	cmd.cobraCmd.PersistentFlags().BoolVarP(
		&common.GlobalOptions.Verbose,
		"verbose",
		"v",
		common.GlobalOptions.Verbose,
		"anchor <command> -v")

	cmd.cobraCmd.PersistentFlags().SortFlags = false
}

func (cmd *anchorCmd) InitSubCommands() {

	//cobra.EnableCommandSorting = false

	// Docker Commands
	cmd.cobraCmd.AddCommand(app.NewCommand(cmd.ctx).GetCobraCmd())

	// Kubernetes Commands
	cmd.cobraCmd.AddCommand(controller.NewCommand(cmd.ctx).GetCobraCmd())

	// Admin
	cmd.cobraCmd.AddCommand(version.NewCommand(cmd.ctx).GetCobraCmd())

	// Auto completion
	cmd.cobraCmd.AddCommand(completion.NewCommand(cmd.cobraCmd, cmd.ctx).GetCobraCmd())
}

func (cmd *anchorCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *anchorCmd) Execute() {
	cmd.InitFlags()
	cmd.InitSubCommands()

	if err := cmd.cobraCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}

func Main(ctx common.Context) {
	NewCommand(ctx).Execute()
}
