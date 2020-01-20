package log

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var namespace = common.GlobalOptions.DockerImageNamespace

type logCmd struct {
	cobraCmd *cobra.Command
	opts     LogOptions
}

type LogOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *logCmd {
	var cobraCmd = &cobra.Command{
		Use:   "logs",
		Short: "Log a running kubernetes pod by name",
		Long:  `Log a running kubernetes pod by name`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cluster.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.ClusterHeadline, "Log")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if err := cluster.LogRunningPod(args[0], namespace); err != nil {
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
	cmd.cobraCmd.PersistentFlags().StringVarP(
		&namespace,
		"namespace",
		"n",
		namespace,
		"anchor cluster log <name> -n <namespace>")
	return nil
}
