package app

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	_select "github.com/ZachiNachshon/anchor/internal/cmd/anchor/app/select"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/app/status"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type appCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context

	addSelectSubCmdFunc func(parent cmd.AnchorCommand, createCmd _select.NewCommandFunc) error
	addStatusSubCmdFunc func(parent cmd.AnchorCommand, createCmd status.NewCommandFunc) error
}

var validArgs = []string{"select", "status"}

type NewCommandFunc func(ctx common.Context, preRunSequence cmd.PreRunSequence) *appCmd

func NewCommand(ctx common.Context, preRunSequence cmd.PreRunSequence) *appCmd {
	var cobraCmd = &cobra.Command{
		Use:       "app",
		Short:     "Application commands",
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunSequence(ctx)
		},
	}

	return &appCmd{
		cobraCmd:            cobraCmd,
		ctx:                 ctx,
		addSelectSubCmdFunc: _select.AddCommand,
		addStatusSubCmdFunc: status.AddCommand,
	}
}

func (c *appCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *appCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, preRunSequence *cmd.AnchorCollaborators, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), preRunSequence.Run)

	err := newCmd.addSelectSubCmdFunc(newCmd, _select.NewCommand)
	if err != nil {
		return err
	}

	err = newCmd.addStatusSubCmdFunc(newCmd, status.NewCommand)
	if err != nil {
		return err
	}

	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
