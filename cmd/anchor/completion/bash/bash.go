package bash

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/spf13/cobra"
	"os"
)

type bashCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
}

func NewCommand(root *cobra.Command) *bashCmd {
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
	}
}

func (cmd *bashCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *bashCmd) InitFlags() {
}

func (cmd *bashCmd) InitSubCommands() {
}
