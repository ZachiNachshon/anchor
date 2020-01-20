package push

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type pushCmd struct {
	cobraCmd *cobra.Command
	opts     PushOptions
}

type PushOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *pushCmd {
	var cobraCmd = &cobra.Command{
		Use:   "push",
		Short: fmt.Sprintf("Push a docker image to repository [%v]", common.GlobalOptions.DockerRegistryDns),
		Long:  fmt.Sprintf("Push a docker image to repository [%v]", common.GlobalOptions.DockerRegistryDns),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := docker.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.DockerHeadline, "Push")

			// Verify Registry is up and running, else start before try to tag & push
			if err := cluster.Registry(false); err != nil {
				logger.Fatal(err.Error())
			}

			if err := docker.Tag(args[0]); err != nil {
				logger.Fatal(err.Error())
			} else {
				if err := docker.Push(args[0]); err != nil {
					logger.Fatal(err.Error())
				}
			}

			logger.PrintCompletion()
		},
	}

	var pushCmd = new(pushCmd)
	pushCmd.cobraCmd = cobraCmd
	pushCmd.opts.CmdRootOptions = opts

	if err := pushCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return pushCmd
}

func (cmd *pushCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *pushCmd) initFlags() error {
	return nil
}
