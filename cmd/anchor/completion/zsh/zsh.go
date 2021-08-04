package zsh

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/spf13/cobra"
	"os"
)

type zshCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
}

func NewCommand(root *cobra.Command) *zshCmd {
	var cobraCmd = &cobra.Command{
		Use:   "zsh",
		Short: "Generate auto completion script for zsh",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenZshCompletion(os.Stdout)
		},
	}

	return &zshCmd{
		cobraCmd: cobraCmd,
	}
}

func (cmd *zshCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *zshCmd) InitFlags() {
}

func (cmd *zshCmd) InitSubCommands() {
}
