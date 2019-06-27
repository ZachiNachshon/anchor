package kubernetes

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
	"strings"
)

var dashboardUrl = "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy"

type dashboardCmd struct {
	cobraCmd *cobra.Command
	opts     DashboardCmdOptions
}

type DashboardCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewDashboardCmd(opts *common.CmdRootOptions) *dashboardCmd {
	var cobraCmd = &cobra.Command{
		Use:   "dashboard",
		Short: "Deploy a Kubernetes dashboard",
		Long:  `Deploy a Kubernetes dashboard`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Deploy Kubernetes Dashboard")

			_ = loadKubeConfig()

			// Kill possible running kubectl proxy
			_ = killKubectlProxy()

			if err := deployKubernetesDashboard(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var dashboardCmd = new(dashboardCmd)
	dashboardCmd.cobraCmd = cobraCmd
	dashboardCmd.opts.CmdRootOptions = opts

	if err := dashboardCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return dashboardCmd
}

func (cmd *dashboardCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *dashboardCmd) initFlags() error {
	return nil
}

func deployKubernetesDashboard() error {
	if exists, err := checkForActiveDashboard(); err != nil {
		return err
	} else if !exists {

		if err := installDashboard(); err != nil {
			return err
		}

		if err := startDashboard(); err != nil {
			return err
		}

	} else {
		logger.Infof("Dashboard already exists, skipping creation")
	}
	return nil
}

func checkForActiveDashboard() (bool, error) {
	getDashboardCmd := "kubectl get deployments kubernetes-dashboard"
	if out, err := common.ShellExec.ExecuteWithOutput(getDashboardCmd); err != nil {
		if strings.Contains(out, "NotFound") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func installDashboard() error {
	logger.Info("\n==> Installing dashboard...\n")
	deployCmd := "kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml"
	createCmd := "kubectl create serviceaccount -n kube-system kubernetes-dashboard"
	createRoleCmd := `kubectl create clusterrolebinding -n kube-system kubernetes-dashboard \
	--clusterrole cluster-admin \
	--serviceaccount kube-system:kubernetes-dashboard
`
	_ = common.ShellExec.Execute(deployCmd)
	_ = common.ShellExec.Execute(createCmd)
	_ = common.ShellExec.Execute(createRoleCmd)

	return nil
}

func startDashboard() error {
	logger.Info("\n==> Starting dashboard...\n")
	logger.Info("Copy the following token and paste as the dashboard TOKEN - \n")
	printSecretCmd := `
kubectl -n kube-system describe secret $(
	kubectl -n kube-system get secret | awk '/^kubernetes-dashboard-token-/{print $1}'
) | awk '$1=="token:"{print $2}'
echo 
`
	// Print token to be used for Dashboard authentication
	_ = common.ShellExec.Execute(printSecretCmd)

	// Start new kubectl proxy
	_ = startKubectlProxy()

	// Open browser and start kubectl proxy
	startProxyCmd := `sleep 2 && open "` + dashboardUrl + `"`
	logger.Info("==> Dashboard available at:\n\n  " + dashboardUrl + "\n")

	return common.ShellExec.Execute(startProxyCmd)
}

func killKubectlProxy() error {
	return common.ShellExec.Execute(`ps -ef | grep "kubectl proxy" | grep -v grep | awk '{print $2}' | xargs kill -9`)
}

func startKubectlProxy() error {
	return common.ShellExec.ExecuteInBackground("kubectl proxy")
}
