package app

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	_select "github.com/ZachiNachshon/anchor/internal/cmd/anchor/app/select"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/app/status"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type appCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"select", "status"}

func NewCommand(ctx common.Context, loadRepoOrFail cmd.LoadRepoOrFailFunc) *appCmd {
	var cobraCmd = &cobra.Command{
		Use:       "app",
		Short:     "Application commands",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return loadRepoOrFail(ctx)
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
	cmd.cobraCmd.AddCommand(_select.NewCommand(cmd.ctx, _select.AppSelect).GetCobraCmd())
	cmd.cobraCmd.AddCommand(status.NewCommand(cmd.ctx, status.AppStatus).GetCobraCmd())
}
