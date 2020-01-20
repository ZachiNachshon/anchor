package deploy

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var namespace = common.GlobalOptions.DockerImageNamespace

type deployCmd struct {
	cobraCmd *cobra.Command
	opts     DeployOptions
}

type DeployOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *deployCmd {
	var cobraCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a fully managed Kubernetes resource",
		Long:  `Deploy a fully managed Kubernetes resource`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cluster.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.ClusterHeadline, "Deploy")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if err := cluster.Deploy(args[0], namespace); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var deployCmd = new(deployCmd)
	deployCmd.cobraCmd = cobraCmd
	deployCmd.opts.CmdRootOptions = opts

	if err := deployCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return deployCmd
}

func (cmd *deployCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *deployCmd) initFlags() error {
	cmd.cobraCmd.PersistentFlags().StringVarP(
		&namespace,
		"namespace",
		"n",
		namespace,
		"anchor cluster deploy <name> -n <namespace>")
	return nil
}
