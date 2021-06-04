package controller

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/spf13/cobra"
)

type controllerCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{""}

func NewCommand(ctx common.Context) *controllerCmd {
	var cobraCmd = &cobra.Command{
		Use:       "controller",
		Short:     "Kubernetes controllers commands",
		Aliases:   []string{"kc"},
		ValidArgs: validArgs,
	}

	var cmd = &controllerCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	if err := cmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return cmd
}

func (cmd *controllerCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *controllerCmd) initSubCommands() error {
	// TODO
	//cmd.cobraCmd.AddCommand(install.NewCommand(cmd.ctx).GetCobraCmd())
	return nil
}
