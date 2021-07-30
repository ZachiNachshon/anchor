package use_context

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
	"github.com/spf13/cobra"
)

type setContextCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, configActions *cfg.ConfigurationActions) *setContextCmd {
	var cobraCmd = &cobra.Command{
		Use:           "use-context",
		Short:         "Sets the current context in the anchor configuration file",
		Long:          `Sets the current context in the anchor configuration file`,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true, // Fatal errors are being logged by parent anchor.go
		SilenceErrors: true, // Fatal errors are being logged by parent anchor.go
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgCtxName := args[0]
			return configActions.UseContext(ctx, cfgCtxName)
		},
	}

	return &setContextCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *setContextCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *setContextCmd) InitFlags() {
}

func (cmd *setContextCmd) InitSubCommands() {
}
