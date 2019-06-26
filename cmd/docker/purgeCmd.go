package docker

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/input"
	"github.com/spf13/cobra"
)

type purgeCmd struct {
	cobraCmd *cobra.Command
	opts     CleanCmdOptions
}

type PurgeCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewPurgeCmd(opts *common.CmdRootOptions) *purgeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "purge",
		Short: "Purge all docker images and containers",
		Long:  `Purge all docker images and containers`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Purge: Containers & Images")
			if err := purgeAll(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var purgeCmd = new(purgeCmd)
	purgeCmd.cobraCmd = cobraCmd
	purgeCmd.opts.CmdRootOptions = opts

	if err := purgeCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return purgeCmd
}

func purgeAll() error {
	in := input.NewYesNoInput()
	if result, err := in.WaitForInput("Purge ALL docker images and containers?"); err != nil || !result {
		logger.Info("skipping.")
	} else {
		if err := common.ShellExec.Execute("docker system prune --all --force"); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *purgeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *purgeCmd) initFlags() error {
	return nil
}
