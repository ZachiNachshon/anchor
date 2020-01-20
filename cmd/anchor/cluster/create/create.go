package create

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type createCmd struct {
	cobraCmd *cobra.Command
	opts     CreateOptions
}

type CreateOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *createCmd {
	var cobraCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a local Kubernetes cluster",
		Long:  `Create a local Kubernetes cluster`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cluster.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.ClusterHeadline, "Create")
			name := common.GlobalOptions.KindClusterName
			if exists, err := cluster.CheckForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if exists {
				logger.Infof("Cluster %v already exists, skipping creation", name)
			} else {
				if err := cluster.Create(); err != nil {
					logger.Fatal(err.Error())
				}
			}

			logger.PrintCompletion()
		},
	}

	var createCmd = new(createCmd)
	createCmd.cobraCmd = cobraCmd
	createCmd.opts.CmdRootOptions = opts

	if err := createCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return createCmd
}

func (cmd *createCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *createCmd) initFlags() error {
	return nil
}
