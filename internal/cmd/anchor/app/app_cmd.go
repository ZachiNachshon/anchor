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
}

var validArgs = []string{"select", "status"}

func NewCommand(ctx common.Context, preRunSequence cmd.PreRunSequence) (*appCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:       "app",
		Short:     "Application commands",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunSequence(ctx)
		},
	}

	var command = &appCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	err := command.InitSubCommands()
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (cmd *appCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *appCmd) InitFlags() error {
	return nil
}

func (cmd *appCmd) InitSubCommands() error {
	if selectCmd, err := _select.NewCommand(cmd.ctx, _select.AppSelect); err != nil {
		return err
	} else {
		cmd.cobraCmd.AddCommand(selectCmd.GetCobraCmd())
	}

	if statusCmd, err := status.NewCommand(cmd.ctx, status.AppStatus); err != nil {
		return err
	} else {
		cmd.cobraCmd.AddCommand(statusCmd.GetCobraCmd())
	}
	return nil
}
