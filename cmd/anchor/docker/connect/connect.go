package connect

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type connectCmd struct {
	cobraCmd *cobra.Command
	opts     ConnectOptions
}

type ConnectOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *connectCmd {
	var cobraCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a docker container by name",
		Long:  `Connect to a docker container by name`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := docker.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.DockerHeadline, "Connect")

			if err := docker.ConnectContainer(args[0]); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var connectCmd = new(connectCmd)
	connectCmd.cobraCmd = cobraCmd
	connectCmd.opts.CmdRootOptions = opts

	if err := connectCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return connectCmd
}

func (cmd *connectCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *connectCmd) initFlags() error {
	return nil
}
