package install

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type installCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, installFunc ControllerInstallFunc) (*installCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:   "install",
		Short: "Install a controller",
		Long:  `Install a controller`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return installFunc(ctx)
		},
	}

	return &installCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}, nil
}

func (cmd *installCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *installCmd) InitFlags() error {
	return nil
}

func (cmd *installCmd) InitSubCommands() error {
	return nil
}
