package status

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context

	initFlagsFunc func(o *statusCmd) error
}

var validStatusOnlyFlagName = "valid-only"
var validStatusOnlyFlagValue = false

var invalidStatusOnlyFlagName = "invalid-only"
var invalidStatusOnlyFlagValue = false

type NewCommandFunc func(ctx common.Context, statusFunc AppStatusFunc) *statusCmd

func NewCommand(ctx common.Context, statusFunc AppStatusFunc) *statusCmd {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: "Check status validity of supported applications",
		Long:  `Check status validity of supported applications`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			o := NewOrchestrator()
			if validStatusOnlyFlagValue && invalidStatusOnlyFlagValue {
				return fmt.Errorf("--%s and --%s flags are mutual exclusive flags", validStatusOnlyFlagName, invalidStatusOnlyFlagName)
			}
			o.validStatusOnly = validStatusOnlyFlagValue
			o.invalidStatusOnly = invalidStatusOnlyFlagValue
			return statusFunc(ctx, o)
		},
	}

	return &statusCmd{
		cobraCmd:      cobraCmd,
		ctx:           ctx,
		initFlagsFunc: initFlags,
	}
}

func (c *statusCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *statusCmd) GetContext() common.Context {
	return c.ctx
}

func initFlags(c *statusCmd) error {
	c.cobraCmd.Flags().BoolVar(
		&validStatusOnlyFlagValue,
		validStatusOnlyFlagName,
		validStatusOnlyFlagValue,
		fmt.Sprintf("anchor app status --%s", validStatusOnlyFlagName))

	c.cobraCmd.Flags().BoolVar(
		&invalidStatusOnlyFlagValue,
		invalidStatusOnlyFlagName,
		invalidStatusOnlyFlagValue,
		fmt.Sprintf("anchor app status --%s", invalidStatusOnlyFlagName))

	c.cobraCmd.PersistentFlags().SortFlags = false
	return nil
}

func AddCommand(parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), AppStatus)
	err := newCmd.initFlagsFunc(newCmd)
	if err != nil {
		return err
	}
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
