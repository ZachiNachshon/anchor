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

type NewCommandFunc func(ctx common.Context, installFunc ControllerInstallFunc) *installCmd

func NewCommand(ctx common.Context, installFunc ControllerInstallFunc) *installCmd {
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
	}
}

func (c *installCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *installCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), ControllerInstall)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
