package pod

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var namespace = common.GlobalOptions.DockerImageNamespace

type podCmd struct {
	cobraCmd *cobra.Command
}

func NewCommand() *podCmd {
	var cobraCmd = &cobra.Command{
		Use:   "pod",
		Short: "Connect to a kubernetes pod (if ^M appear as Enter, run - stty sane)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cluster.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.ClusterHeadline, "Connect Pod")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if err := cluster.ConnectToPod(args[0], namespace); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var podCmd = new(podCmd)
	podCmd.cobraCmd = cobraCmd

	if err := podCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return podCmd
}

func (cmd *podCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *podCmd) initFlags() error {
	cmd.cobraCmd.PersistentFlags().StringVarP(
		&namespace,
		"namespace",
		"n",
		namespace,
		"anchor cluster connect pod <name> -n <namespace>")
	return nil
}
