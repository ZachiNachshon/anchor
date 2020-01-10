package log

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type logCmd struct {
	cobraCmd *cobra.Command
	opts     LogOptions
}

type LogOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *logCmd {
	var cobraCmd = &cobra.Command{
		Use:   "log",
		Short: "Log a running docker container",
		Long:  `Run a running docker container`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.DockerHeadline, "Log")

			if err := docker.LogContainer(args[0]); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var logCmd = new(logCmd)
	logCmd.cobraCmd = cobraCmd
	logCmd.opts.CmdRootOptions = opts

	if err := logCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return logCmd
}

func (cmd *logCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *logCmd) initFlags() error {
	return nil
}
