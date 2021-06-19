package app

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/app/install"
	"github.com/ZachiNachshon/anchor/cmd/anchor/app/status"
	"github.com/ZachiNachshon/anchor/cmd/anchor/app/uninstall"
	"github.com/ZachiNachshon/anchor/cmd/anchor/app/versions"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/app"
	"github.com/spf13/cobra"
)

type appCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"install", "uninstall", "status", "versions"}

func NewCommand(ctx common.Context) *appCmd {
	var cobraCmd = &cobra.Command{
		Use:       "app",
		Short:     "Application commands",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
	}

	var cmd = &appCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	if err := cmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return cmd
}

func (cmd *appCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *appCmd) initSubCommands() error {
	actions := app.DefineApplicationActions()
	cmd.cobraCmd.AddCommand(install.NewCommand(cmd.ctx, actions).GetCobraCmd())
	cmd.cobraCmd.AddCommand(uninstall.NewCommand(cmd.ctx).GetCobraCmd())
	cmd.cobraCmd.AddCommand(status.NewCommand(cmd.ctx).GetCobraCmd())
	cmd.cobraCmd.AddCommand(versions.NewCommand(cmd.ctx).GetCobraCmd())
	return nil
}
