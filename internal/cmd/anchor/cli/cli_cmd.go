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

func NewCommand(ctx common.Context, preRunSequence cmd.PreRunSequence) *cliCmd {
	var cobraCmd = &cobra.Command{
		Use:       "cli",
		Short:     "CLI applications",
		Aliases:   []string{"a"},
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunSequence(ctx)
		},
	}

	var cmd = &cliCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	cmd.InitSubCommands()
	return cmd
}

func (cmd *cliCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *cliCmd) InitFlags() {
}

func (cmd *cliCmd) InitSubCommands() {
	cmd.cobraCmd.AddCommand(versions.NewCommand(cmd.ctx, versions.CliVersions).GetCobraCmd())
}
