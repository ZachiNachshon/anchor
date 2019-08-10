package kubernetes

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var dashboardUrl = "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy"

var shouldDeleteDashboard = false

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
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {
				_ = loadKubeConfig()

				if shouldDeleteDashboard {
					if err := uninstallDashboard(); err != nil {
						logger.Fatal(err.Error())
					}
				} else {
					if err := deployKubernetesDashboard(); err != nil {
						logger.Fatal(err.Error())
					}
				}
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
	// TODO: Allow force creation by flag even if dashboard exists ?
	cmd.cobraCmd.Flags().BoolVarP(
		&shouldDeleteDashboard,
		"Delete Kubernetes dashboard",
		"d",
		shouldDeleteDashboard,
		"anchor cluster dashboard -d")
	return nil
}

func deployKubernetesDashboard() error {
	if exists, err := checkForActiveDashboard(); err != nil {
		return err
	} else if !exists {

		// Kill possible running kubectl proxy
		_ = killKubectlProxy()

		if err := installDashboard(); err != nil {
			return err
		}

		if err := startDashboard(); err != nil {
			return err
		}

	} else {
		logger.Info("Dashboard already exists, skipping creation.")
		return printDashboardInfo()
	}
	return nil
}

func checkForActiveDashboard() (bool, error) {
	getDashboardCmd := "kubectl get deployments kubernetes-dashboard --namespace=kube-system"
	if out, err := common.ShellExec.ExecuteWithOutput(getDashboardCmd); err != nil {
		if strings.Contains(out, "NotFound") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func uninstallDashboard() error {
	if exists, err := checkForActiveDashboard(); err != nil {
		return err
	} else if exists {
		logger.Info("\n==> Uninstalling dashboard...\n")

		// Kill possible running kubectl proxy
		_ = killKubectlProxy()

		if file, err := ioutil.TempFile(os.TempDir(), "anchor-dashboard-manifest.yaml"); err != nil {
			return err
		} else {
			// Remove after finished
			defer os.Remove(file.Name())

			if _, err := file.WriteString(config.KubernetesDashboardManifest); err != nil {
				return err
			} else {
				removeDashboardCmd := fmt.Sprintf("cat %v | kubectl delete -f -",
					file.Name())

				if err := common.ShellExec.Execute(removeDashboardCmd); err != nil {
					return err
				}
			}
		}
	} else {
		logger.Info("Dashboard does not exists, nothing to delete.")
	}
	return nil
}

func installDashboard() error {
	logger.Info("\n==> Installing dashboard...\n")
	if file, err := ioutil.TempFile(os.TempDir(), "anchor-dashboard-manifest.yaml"); err != nil {
		return err
	} else {
		// Remove after finished
		defer os.Remove(file.Name())

		if _, err := file.WriteString(config.KubernetesDashboardManifest); err != nil {
			return err
		} else {
			createDashboardCmd := fmt.Sprintf("cat %v | kubectl apply -f -",
				file.Name())

			if err := common.ShellExec.Execute(createDashboardCmd); err != nil {
				return err
			}

			createCmd := "kubectl create serviceaccount -n kube-system kubernetes-dashboard"
			createRoleCmd := `kubectl create clusterrolebinding -n kube-system kubernetes-dashboard \
			--clusterrole cluster-admin \
			--serviceaccount kube-system:kubernetes-dashboard
		`
			_ = common.ShellExec.Execute(createCmd)
			_ = common.ShellExec.Execute(createRoleCmd)
		}
	}
	return nil
}

func startDashboard() error {
	logger.Info("\n==> Starting dashboard...\n")

	// Start new kubectl proxy
	_ = startKubectlProxy()

	// Sleep for 3 secs to allow kubectl proxy to start
	time.Sleep(5 * time.Second)

	// Wait until dashboard pod is ready
	if err := waitForDashboardPod(); err != nil {
		// TODO: handle error gracefully
		//return err
	}

	_ = printDashboardInfo()

	// Open browser and start kubectl proxy
	startProxyCmd := fmt.Sprintf(`open "%v"`, dashboardUrl)

	return common.ShellExec.Execute(startProxyCmd)
}

func getSecret() error {
	// TODO: Prevent overload of secrets if there is no dashboard secret
	printSecretCmd := `
kubectl -n kube-system describe secret $(
	kubectl -n kube-system get secret | awk '/^kubernetes-dashboard-token-/{print $1}'
) | awk '$1=="token:"{print $2}'
echo 
`
	// Print token to be used for Dashboard authentication
	return common.ShellExec.Execute(printSecretCmd)
}

func waitForDashboardPod() error {
	logger.Info("Waiting for dashboard pod to be ready (2m timeout)...")
	waitContainerCmd := fmt.Sprintf("kubectl wait -n kube-system -l k8s-app=kubernetes-dashboard --timeout=2m --for=condition=Ready pod")
	return common.ShellExec.Execute(waitContainerCmd)
}

func killKubectlProxy() error {
	return common.ShellExec.Execute(`ps -ef | grep "kubectl proxy" | grep -v grep | awk '{print $2}' | xargs kill -9`)
}

// Creates a proxy server or application-level gateway between localhost and the Kubernetes API Server.
func startKubectlProxy() error {
	return common.ShellExec.ExecuteInBackground("kubectl proxy")
}

func printDashboardInfo() error {
	if exists, err := checkForActiveDashboard(); err != nil {
		return err
	} else if exists {
		logger.Info(`
Dashboard:
----------`)
		logger.Info("Copy the following and paste as the dashboard TOKEN - \n")
		_ = getSecret()
		logger.Info("==> Dashboard available at:\n\n  " + dashboardUrl + "\n")
	}
	return nil
}
