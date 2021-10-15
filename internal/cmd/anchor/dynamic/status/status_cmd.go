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

type NewCommandFunc func(ctx common.Context, commandFolderName string, statusFunc DynamicStatusFunc) *statusCmd

func NewCommand(ctx common.Context, commandFolderName string, statusFunc DynamicStatusFunc) *statusCmd {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: "Check status validity of supported dynamic commands",
		Long:  `Check status validity of supported dynamic commands`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			o := NewOrchestrator(commandFolderName)
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
		fmt.Sprintf("anchor <anchor-folder-item> status --%s", validStatusOnlyFlagName))

	c.cobraCmd.Flags().BoolVar(
		&invalidStatusOnlyFlagValue,
		invalidStatusOnlyFlagName,
		invalidStatusOnlyFlagValue,
		fmt.Sprintf("anchor <anchor-folder-item> status --%s", invalidStatusOnlyFlagName))

	c.cobraCmd.PersistentFlags().SortFlags = false
	return nil
}

func AddCommand(parent cmd.AnchorCommand, commandFolderName string, createCmd NewCommandFunc) error {
	newCmd := createCmd(parent.GetContext(), commandFolderName, DynamicStatus)
	err := newCmd.initFlagsFunc(newCmd)
	if err != nil {
		return err
	}
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
