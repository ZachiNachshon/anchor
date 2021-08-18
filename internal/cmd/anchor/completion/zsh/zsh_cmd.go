package zsh

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
	"os"
)

type zshCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
}

type NewCommandFunc func(rootCmd cmd.AnchorCommand) *zshCmd

func NewCommand(rootCmd cmd.AnchorCommand) *zshCmd {
	var cobraCmd = &cobra.Command{
		Use:   "zsh",
		Short: "Generate auto completion script for zsh",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GetCobraCmd().GenZshCompletion(os.Stdout)
		},
	}

	return &zshCmd{
		cobraCmd: cobraCmd,
	}
}

func (c *zshCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *zshCmd) GetContext() common.Context {
	return nil
}

func AddCommand(root cmd.AnchorCommand, parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(root)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
