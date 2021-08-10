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

func NewCommand(
	ctx common.Context,
	cfgManager config.ConfigManager,
	useContextFunc ConfigUseContextFunc) (*setContextCmd, error) {

	var cobraCmd = &cobra.Command{
		Use:           "use-context",
		Short:         "Sets the current context in the anchor configuration file",
		Long:          `Sets the current context in the anchor configuration file`,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true, // Fatal errors are being logged by parent anchor.go
		SilenceErrors: true, // Fatal errors are being logged by parent anchor.go
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgCtxName := args[0]
			return useContextFunc(ctx, cfgCtxName, cfgManager)
		},
	}

	return &setContextCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}, nil
}

func (cmd *setContextCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *setContextCmd) InitFlags() error {
	return nil
}

func (cmd *setContextCmd) InitSubCommands() error {
	return nil
}
