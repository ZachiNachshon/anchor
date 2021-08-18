package view

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

type viewCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

type NewCommandFunc func(
	ctx common.Context,
	cfgManager config.ConfigManager,
	editFunc ConfigViewFunc) *viewCmd

func NewCommand(
	ctx common.Context,
	cfgManager config.ConfigManager,
	viewFunc ConfigViewFunc) *viewCmd {

	var cobraCmd = &cobra.Command{
		Use:   "view",
		Short: "Display configuration file settings",
		Long:  `Display configuration file settings`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewFunc(ctx, cfgManager)
		},
	}

	return &viewCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (c *viewCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *viewCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(
	parent cmd.AnchorCommand,
	cfgManager config.ConfigManager,
	createCmd NewCommandFunc) error {

	newCmd := createCmd(parent.GetContext(), cfgManager, ConfigView)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
