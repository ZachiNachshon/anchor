package use_context

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

type setContextCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

type NewCommandFunc func(
	ctx common.Context,
	cfgManager config.ConfigManager,
	useContextFunc ConfigUseContextFunc) *setContextCmd

func NewCommand(
	ctx common.Context,
	cfgManager config.ConfigManager,
	useContextFunc ConfigUseContextFunc) *setContextCmd {

	var cobraCmd = &cobra.Command{
		Use:           "use-context",
		Short:         "Sets the current context in the anchor configuration file",
		Long:          `Sets the current context in the anchor configuration file`,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true, // Fatal errors are being logged by parent anchor.go
		SilenceErrors: true, // Fatal errors are being logged by parent anchor.go
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgCtxName := args[0]
			return useContextFunc(ctx, NewOrchestrator(cfgManager, cfgCtxName))
		},
	}

	return &setContextCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (c *setContextCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *setContextCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(
	parent cmd.AnchorCommand,
	cfgManager config.ConfigManager,
	createCmd NewCommandFunc) error {

	newCmd := createCmd(parent.GetContext(), cfgManager, ConfigUseContext)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
