package cli

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/cli/versions"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type cliCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"versions"}

func NewCommand(ctx common.Context, preRunSequence cmd.PreRunSequence) (*cliCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:       "cli",
		Short:     "CLI applications",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunSequence(ctx)
		},
	}

	var command = &cliCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	err := command.InitSubCommands()
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (cmd *cliCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *cliCmd) InitFlags() error {
	return nil
}

func (cmd *cliCmd) InitSubCommands() error {
	if versionCmd, err := versions.NewCommand(cmd.ctx, versions.CliVersions); err != nil {
		return err
	} else {
		cmd.cobraCmd.AddCommand(versionCmd.GetCobraCmd())
	}
	return nil
}
