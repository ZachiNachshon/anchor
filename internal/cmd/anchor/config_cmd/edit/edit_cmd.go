package edit

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

type editCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

type NewCommandFunc func(
	ctx common.Context,
	cfgManager config.ConfigManager,
	editFunc ConfigEditFunc) *editCmd

func NewCommand(
	ctx common.Context,
	cfgManager config.ConfigManager,
	editFunc ConfigEditFunc) *editCmd {

	var cobraCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		Long:  `Edit configuration file`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editFunc(ctx, cfgManager)
		},
	}

	return &editCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (c *editCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *editCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(
	parent cmd.AnchorCommand,
	cfgManager config.ConfigManager,
	createCmd NewCommandFunc) error {

	newCmd := createCmd(parent.GetContext(), cfgManager, ConfigEdit)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
