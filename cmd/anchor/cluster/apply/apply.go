package apply

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var namespace = common.GlobalOptions.DockerImageNamespace

type applyCmd struct {
	cobraCmd *cobra.Command
	opts     ApplyOptions
}

type ApplyOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *applyCmd {
	var cobraCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply a Kubernetes manifest resource",
		Long:  `Apply a Kubernetes manifest resource`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.ClusterHeadline, "Apply")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if _, err := cluster.Apply(args[0], namespace); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var deployCmd = new(applyCmd)
	deployCmd.cobraCmd = cobraCmd
	deployCmd.opts.CmdRootOptions = opts

	if err := deployCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return deployCmd
}

func (cmd *applyCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *applyCmd) initFlags() error {
	cmd.cobraCmd.PersistentFlags().StringVarP(
		&namespace,
		"namespace",
		"n",
		namespace,
		"anchor cluster apply <name> -n <namespace>")
	return nil
}
