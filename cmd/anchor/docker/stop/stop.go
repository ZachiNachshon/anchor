package stop

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type stopCmd struct {
	cobraCmd *cobra.Command
	opts     StopOptions
}

type StopOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *stopCmd {
	var cobraCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop a docker container",
		Long:  `Stop a docker container`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := docker.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.DockerHeadline, "Stop")

			if err := docker.StopContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := docker.RemoveContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var stopCmd = new(stopCmd)
	stopCmd.cobraCmd = cobraCmd
	stopCmd.opts.CmdRootOptions = opts

	if err := stopCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return stopCmd
}

func (cmd *stopCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *stopCmd) initFlags() error {
	return nil
}
