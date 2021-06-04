package bash

import (
	"github.com/spf13/cobra"
	"os"
)

type bashCmd struct {
	cobraCmd *cobra.Command
}

func NewCommand(root *cobra.Command) *bashCmd {
	var cobraCmd = &cobra.Command{
		Use:   "bash",
		Short: "Generate auto completion script for bash",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_ = root.GenBashCompletion(os.Stdout)
		},
	}

	return &bashCmd{
		cobraCmd: cobraCmd,
	}
}

func (cmd *bashCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}
