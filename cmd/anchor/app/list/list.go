package list

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/spf13/cobra"
)

type listCmd struct {
	common.CliCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context) *listCmd {
	var cobraCmd = &cobra.Command{
		Use:   "list",
		Short: "List all supported applications",
		Long:  `List all supported applications`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if registry, err := locator.FromRegistry(ctx.Registry()); err != nil {
				logger.Fatal(err.Error())
			} else {
				printSupportedApplications(registry)
			}
		},
	}

	return &listCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *listCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *listCmd) InitFlags() {
}

func (cmd *listCmd) InitSubCommands() {
}

func printSupportedApplications(locate locator.Locator) {

	// TODO: Move print from locator to another utility struct

	//locate.Print()
}
