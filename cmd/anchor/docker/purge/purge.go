package purge

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type purgeCmd struct {
	cobraCmd *cobra.Command
	opts     PurgeOptions
}

type PurgeOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *purgeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "purge",
		Short: "Purge all docker images and containers",
		Long:  `Purge all docker images and containers`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.DockerHeadline, "Purge")

			if err := docker.PurgeAll(); err != nil {
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

func (cmd *purgeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *purgeCmd) initFlags() error {
	return nil
}
