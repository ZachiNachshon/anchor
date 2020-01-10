package run

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type runCmd struct {
	cobraCmd *cobra.Command
	opts     RunOptions
}

type RunOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *runCmd {
	var cobraCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a docker container",
		Long:  `Run a docker container`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.DockerHeadline, "Run")

			if err := docker.StopContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := docker.RemoveContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := docker.RunContainer(args[0]); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var runCmd = new(runCmd)
	runCmd.cobraCmd = cobraCmd
	runCmd.opts.CmdRootOptions = opts

	if err := runCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return runCmd
}

func (cmd *runCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *runCmd) initFlags() error {
	cmd.cobraCmd.PersistentFlags().BoolVar(
		&common.GlobalOptions.DockerRunAutoLog,
		"auto-log",
		common.GlobalOptions.DockerRunAutoLog,
		"anchor docker run <image> --auto-log=false")

	return nil
}
