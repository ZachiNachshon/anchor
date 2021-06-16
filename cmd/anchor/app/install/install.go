package install

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/prompter"
	"github.com/spf13/cobra"
)

type installCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context) *installCmd {
	var cobraCmd = &cobra.Command{
		Use:   "install",
		Short: "Install an application",
		Long:  `Install an application`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if orchestrator, err := prompter.NewOrchestrator(ctx); err != nil {
				logger.Fatal(err.Error())
			} else {
				if selection, err := orchestrator.OrchestrateAppInstructionSelection(); err != nil {
					logger.Fatal(err.Error())
				} else {
					logger.Infof("Selected: %v", selection.Id)
				}
			}
		},
	}

	return &installCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *installCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *installCmd) InitFlags() {
}

func (cmd *installCmd) InitSubCommands() {
}
