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

	addVersionSubCmdFunc func(parent cmd.AnchorCommand, createCmd versions.NewCommandFunc) error
}

var validArgs = []string{"versions"}

type NewCommandFunc func(ctx common.Context, preRunSequence cmd.PreRunSequence) *cliCmd

func NewCommand(ctx common.Context, preRunSequenceFunc cmd.PreRunSequence) *cliCmd {
	var cobraCmd = &cobra.Command{
		Use:       "cli",
		Short:     "CLI applications",
		ValidArgs: validArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunSequenceFunc(ctx)
		},
	}

	c := &cliCmd{
		cobraCmd:             cobraCmd,
		ctx:                  ctx,
		addVersionSubCmdFunc: versions.AddCommand,
	}
	return c
}

func (c *cliCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *cliCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, preRunSequence *cmd.AnchorCollaborators, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), preRunSequence.Run)

	err := newCmd.addVersionSubCmdFunc(newCmd, versions.NewCommand)
	if err != nil {
		return err
	}

	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
