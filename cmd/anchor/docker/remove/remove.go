package remove

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type removeCmd struct {
	cobraCmd *cobra.Command
	opts     RemoveOptions
}

type RemoveOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *removeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove docker containers and images",
		Long:  `Remove docker containers and images`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.DockerHeadline, "Remove")

			if err := docker.StopContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := docker.RemoveContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := docker.RemoveImages(args[0]); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var removeCmd = new(removeCmd)
	removeCmd.cobraCmd = cobraCmd
	removeCmd.opts.CmdRootOptions = opts

	if err := removeCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return removeCmd
}

func (cmd *removeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *removeCmd) initFlags() error {
	return nil
}
