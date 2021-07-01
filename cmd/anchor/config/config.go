package config

import (
	printcmd "github.com/ZachiNachshon/anchor/cmd/anchor/config/print"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
	"github.com/spf13/cobra"
)

type configCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"print"}

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
	actions := cfg.DefineConfigurationActions()
	cmd.cobraCmd.AddCommand(printcmd.NewCommand(cmd.ctx, actions).GetCobraCmd())
}
