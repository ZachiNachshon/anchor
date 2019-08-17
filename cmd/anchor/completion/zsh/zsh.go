package zsh

import (
	"github.com/spf13/cobra"
	"os"
)

type zshCmd struct {
	cobraCmd *cobra.Command
}

func NewCommand(root *cobra.Command) *zshCmd {
	var cobraCmd = &cobra.Command{
		Use:   "zsh",
		Short: "Generate auto completion script for zsh",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_ = root.GenZshCompletion(os.Stdout)
		},
	}

	var zshCmd = new(zshCmd)
	zshCmd.cobraCmd = cobraCmd

	return zshCmd
}

func (cmd *zshCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}
