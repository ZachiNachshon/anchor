package cli

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/cli/versions"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cli"
	"github.com/ZachiNachshon/anchor/pkg/root"
	"github.com/spf13/cobra"
)

type cliCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"versions"}

func NewCommand(ctx common.Context, rootActions *root.RootCommandActions) *cliCmd {
	var cobraCmd = &cobra.Command{
		Use:       "cli",
		Short:     "CLI applications",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			rootActions.LoadRepoOrFail(ctx)
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
	actions := cli.DefineCliActions()
	cmd.cobraCmd.AddCommand(versions.NewCommand(cmd.ctx, actions).GetCobraCmd())
}
