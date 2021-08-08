package config_cmd

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/edit"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/set_context_entry"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/use_context"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/view"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type configCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"edit", "set-context-entry", "use-context", "view"}

func NewCommand(ctx common.Context) *configCmd {
	var cobraCmd = &cobra.Command{
		Use:       "config",
		Short:     "Configuration file management",
		Long:      `Configuration file management`,
		ValidArgs: validArgs,
	}

	var cmd = &configCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	cmd.InitSubCommands()
	return cmd
}

func (cmd *configCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *configCmd) InitFlags() {
}

func (cmd *configCmd) InitSubCommands() {
	cmd.cobraCmd.AddCommand(view.NewCommand(cmd.ctx, view.ConfigView).GetCobraCmd())
	cmd.cobraCmd.AddCommand(edit.NewCommand(cmd.ctx, edit.ConfigEdit).GetCobraCmd())
	cmd.cobraCmd.AddCommand(use_context.NewCommand(cmd.ctx, use_context.ConfigUseContext).GetCobraCmd())
	cmd.cobraCmd.AddCommand(set_context_entry.NewCommand(cmd.ctx, set_context_entry.ConfigSetContextEntry).GetCobraCmd())
}
