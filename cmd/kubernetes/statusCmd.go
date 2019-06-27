package kubernetes

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	cobraCmd *cobra.Command
	opts     CreateCmdOptions
}

type StatusCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewStatusCmd(opts *common.CmdRootOptions) *statusCmd {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: "Get active cluster status",
		Long:  `Get active cluster status`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			name := common.GlobalOptions.KindClusterName
			logger.PrintHeadline("Retrieve Cluster Status")
			_ = loadKubeConfig()
			if err := printClusterStatus(name); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var statusCmd = new(statusCmd)
	statusCmd.cobraCmd = cobraCmd
	statusCmd.opts.CmdRootOptions = opts

	if err := statusCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return statusCmd
}

func (cmd *statusCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *statusCmd) initFlags() error {
	return nil
}

func printClusterStatus(name string) error {
	if exists, err := checkForActiveCluster(name); err != nil {
		return err
	} else if !exists {
		logger.Info("No active cluster.")
	} else {
		logger.Infof("Found active %v cluster !", common.GlobalOptions.KindClusterName)

		getConfigCmd := "kind get kubeconfig-path"
		//logger.Info("\nConfiguration:\n")
		logger.Info(`
Configuration:
--------------`)
		_ = common.ShellExec.Execute(getConfigCmd)

		getNodesCmd := fmt.Sprintf("kind get nodes --name %s", common.GlobalOptions.KindClusterName)
		//logger.Info("\nNodes:\n")
		logger.Info(`
Nodes:
------`)
		_ = common.ShellExec.Execute(getNodesCmd)

		getPodsCmd := "kubectl get pods -o wide"
		//logger.Info("\nPods:\n")
		logger.Info(`
Pods:
-----`)
		_ = common.ShellExec.Execute(getPodsCmd)

		getServicesCmd := "kubectl get services"
		//logger.Info("\nServices:\n")
		logger.Info(`
Services:
---------`)
		_ = common.ShellExec.Execute(getServicesCmd)
	}
	return nil
}
