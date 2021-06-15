package list

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/printer"
	"github.com/spf13/cobra"
)

type listCmd struct {
	models.AnchorCommand
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
			l, err := locator.FromRegistry(ctx.Registry())
			if err != nil {
				logger.Fatal(err.Error())
			}

			p, err := printer.FromRegistry(ctx.Registry())
			if err != nil {
				logger.Fatal(err.Error())
			}

			printSupportedApplications(l, p)
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

func printSupportedApplications(l locator.Locator, p printer.Printer) {
	if err := l.Scan(); err != nil {
		logger.Fatalf("Scanning for available application failed.")
	} else {
		p.PrintApplications(l.Applications())
	}
}
