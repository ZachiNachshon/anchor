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

type NewCommandFunc func(ctx common.Context, statusFunc AppStatusFunc) *statusCmd

func NewCommand(ctx common.Context, statusFunc AppStatusFunc) *statusCmd {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: "Check status validity of supported applications",
		Long:  `Check status validity of supported applications`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return statusFunc(ctx, NewOrchestrator())
		},
	}

	return &statusCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (c *statusCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *statusCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), AppStatus)
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
