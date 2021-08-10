package bash

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/spf13/cobra"
	"os"
)

type bashCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
}

func NewCommand(root *cobra.Command) (*bashCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:   "bash",
		Short: "Generate auto completion script for bash",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}

	return &bashCmd{
		cobraCmd: cobraCmd,
	}, nil
}

func (cmd *bashCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *bashCmd) InitFlags() error {
	return nil
}

func (cmd *bashCmd) InitSubCommands() error {
	return nil
}
