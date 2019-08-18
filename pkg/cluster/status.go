package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"strings"
)

func printClusterInfo() error {
	return common.ShellExec.Execute("kubectl cluster-info")
}

func printConfiguration() error {
	name := common.GlobalOptions.KindClusterName
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
Usage...: export KUBECONFIG="$(kind get kubeconfig-path --name=%s)"`, strings.Trim(out, "\n"), name))
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
	getServicesCmd := "kubectl get services --all-namespaces=true"
	return common.ShellExec.Execute(getServicesCmd)
}

func printDeployments() error {
	logger.Info(`
Deployments:
------------`)
	getDeploymentsCmd := "kubectl get deployments --all-namespaces=true"
	return common.ShellExec.Execute(getDeploymentsCmd)
}

func printNodes() error {
	logger.Info(`
Nodes:
------`)
	getNodesCmd := "kubectl get nodes --all-namespaces=true"
	return common.ShellExec.Execute(getNodesCmd)
}

func printPods() error {
	logger.Info(`
Pods:
-----`)
	//getPodsCmd := "kubectl get pods -o wide"
	getPodsCmd := "kubectl get pods --all-namespaces=true"
	return common.ShellExec.Execute(getPodsCmd)
}

func Status() error {
	name := common.GlobalOptions.KindClusterName
	logger.PrintCommandHeader(fmt.Sprintf("Print %v cluster status", name))

	// Double check since other command might call directly to this method
	if exists, err := CheckForActiveCluster(name); err != nil {
		return err
	} else if !exists {
		logger.Info("\nNo active cluster.")
	} else {
		logger.Infof("Found active %v cluster !\n", name)

		_ = printClusterInfo()
		_ = PrintDashboardInfo()
		_ = PrintRegistryInfo()
		_ = printConfiguration()
		_ = printNamespaces()
		_ = printIngress()
		_ = printServices()
		_ = printDeployments()
		_ = printNodes()
		_ = printPods()

	}
	return nil
}
