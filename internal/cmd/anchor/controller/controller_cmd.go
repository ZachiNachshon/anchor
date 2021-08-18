package controller

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/controller/install"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type controllerCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context

	addInstallSubCmdFunc func(parent cmd.AnchorCommand, createCmd install.NewCommandFunc) error
}

var validArgs = []string{""}

type NewCommandFunc func(ctx common.Context, preRunSequence cmd.PreRunSequence) *controllerCmd

func NewCommand(ctx common.Context, preRunSequence cmd.PreRunSequence) *controllerCmd {
	var cobraCmd = &cobra.Command{
		Use:       "controller",
		Short:     "Kubernetes controllers commands",
		Aliases:   []string{"kc"},
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunSequence(ctx)
		},
	}

	return &controllerCmd{
		cobraCmd:             cobraCmd,
		ctx:                  ctx,
		addInstallSubCmdFunc: install.AddCommand,
	}
}

func (c *controllerCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *controllerCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, preRunSequence *cmd.AnchorCollaborators, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), preRunSequence.Run)

	err := newCmd.addInstallSubCmdFunc(newCmd, install.NewCommand)
	if err != nil {
		return err
	}

	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
