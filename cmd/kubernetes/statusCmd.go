package kubernetes

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
	"strings"
)

type statusCmd struct {
	cobraCmd *cobra.Command
	opts     StatusCmdOptions
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
	var clusterName = common.GlobalOptions.KindClusterName
	if exists, err := checkForActiveCluster(name); err != nil {
		return err
	} else if !exists {
		logger.Info("No active cluster.")
	} else {
		logger.Infof("Found active %v cluster !", clusterName)
		logger.Infof("\nDashboard available at:\n  %s", dashboardUrl)

		_ = printConfiguration(clusterName)
		_ = printNamespaces()
		_ = printIngress()
		_ = printServices()
		_ = printDeployments()
		_ = printNodes()
		_ = printPods()

	}
	return nil
}

func printConfiguration(clusterName string) error {
	logger.Info(`
Configuration:
--------------`)
	getConfigCmd := "kind get kubeconfig-path"
	if out, err := common.ShellExec.ExecuteWithOutput(getConfigCmd); err != nil {
		// TODO: Change to warn
		logger.Infof("Something went wrong, error: %v", err.Error())
		return err
	} else {
		logger.Info(fmt.Sprintf(`Path:...: %s
Usage...: export KUBECONFIG="$(kind get kubeconfig-path --name=%s)"`, strings.Trim(out, "\n"), clusterName))
	}
	return nil
}

func printNamespaces() error {
	logger.Info(`
Namespaces:
-----------`)
	getNamespacesCmd := "kubectl get namespaces"
	return common.ShellExec.Execute(getNamespacesCmd)
}

func printIngress() error {
	logger.Info(`
Ingress:
--------`)
	getIngressCmd := "kubectl get ingress"
	return common.ShellExec.Execute(getIngressCmd)
}

func printServices() error {
	logger.Info(`
Services:
---------`)
	getServicesCmd := "kubectl get services"
	return common.ShellExec.Execute(getServicesCmd)
}

func printDeployments() error {
	logger.Info(`
Services:
---------`)
	getDeploymentsCmd := "kubectl get deployments"
	return common.ShellExec.Execute(getDeploymentsCmd)
}

func printNodes() error {
	logger.Info(`
Nodes:
------`)
	getNodesCmd := "kubectl get nodes"
	return common.ShellExec.Execute(getNodesCmd)
}

func printPods() error {
	logger.Info(`
Pods:
-----`)
	//getPodsCmd := "kubectl get pods -o wide"
	getPodsCmd := "kubectl get pods"
	return common.ShellExec.Execute(getPodsCmd)
}
