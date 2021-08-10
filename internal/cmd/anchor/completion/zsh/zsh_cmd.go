package zsh

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/spf13/cobra"
	"os"
)

type zshCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
}

func NewCommand(root *cobra.Command) (*zshCmd, error) {
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
	}, nil
}

func (cmd *zshCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *zshCmd) InitFlags() error {
	return nil
}

func (cmd *zshCmd) InitSubCommands() error {
	return nil
}
