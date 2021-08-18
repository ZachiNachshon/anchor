package bash

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
	"os"
)

type bashCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
}

type NewCommandFunc func(rootCmd cmd.AnchorCommand) *bashCmd

func NewCommand(rootCmd cmd.AnchorCommand) *bashCmd {
	var cobraCmd = &cobra.Command{
		Use:   "bash",
		Short: "Generate auto completion script for bash",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GetCobraCmd().GenBashCompletion(os.Stdout)
		},
	}

	return &bashCmd{
		cobraCmd: cobraCmd,
	}
}

func (c *bashCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *bashCmd) GetContext() common.Context {
	return nil
}

func AddCommand(root cmd.AnchorCommand, parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(root)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
