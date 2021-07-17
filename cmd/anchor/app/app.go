package app

import (
	_select "github.com/ZachiNachshon/anchor/cmd/anchor/app/select"
	"github.com/ZachiNachshon/anchor/cmd/anchor/app/status"
	"github.com/ZachiNachshon/anchor/cmd/anchor/app/versions"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/app"
	"github.com/spf13/cobra"
)

type appCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"install", "uninstall", "status", "versions", "list"}

func NewCommand(ctx common.Context, loadRepoOrFail func(ctx common.Context)) *appCmd {
	var cobraCmd = &cobra.Command{
		Use:       "app",
		Short:     "Application commands",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			loadRepoOrFail(ctx)
		},
	}

	var cmd = &appCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	cmd.InitSubCommands()
	return cmd
}

func (cmd *appCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *appCmd) InitFlags() {
}

func (cmd *appCmd) InitSubCommands() {
	actions := app.DefineApplicationActions()
	cmd.cobraCmd.AddCommand(_select.NewCommand(cmd.ctx, actions).GetCobraCmd())
	cmd.cobraCmd.AddCommand(status.NewCommand(cmd.ctx).GetCobraCmd())
	cmd.cobraCmd.AddCommand(versions.NewCommand(cmd.ctx).GetCobraCmd())
}
