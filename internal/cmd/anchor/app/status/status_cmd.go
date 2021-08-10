package status

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, statusFunc AppStatusFunc) (*statusCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: "Check status validity of supported applications",
		Long:  `Check status validity of supported applications`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return statusFunc(ctx)
		},
	}

	return &statusCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}, nil
}

func (cmd *statusCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *statusCmd) InitFlags() error {
	return nil
}

func (cmd *statusCmd) InitSubCommands() error {
	return nil
}
