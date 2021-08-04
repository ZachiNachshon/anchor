package cli

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/cli/versions"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/spf13/cobra"
)

type cliCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"versions"}

func NewCommand(ctx common.Context, loadRepoOrFail func(ctx common.Context)) *cliCmd {
	var cobraCmd = &cobra.Command{
		Use:       "cli",
		Short:     "CLI applications",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			loadRepoOrFail(ctx)
		},
	}

	var cmd = &cliCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	cmd.InitSubCommands()
	return cmd
}

func (cmd *cliCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *cliCmd) InitFlags() {
}

func (cmd *cliCmd) InitSubCommands() {
	//actions := app.DefineApplicationActions()
	cmd.cobraCmd.AddCommand(versions.NewCommand(cmd.ctx).GetCobraCmd())
}
