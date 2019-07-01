package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/anchor/config"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var shouldDeleteRegistry = false

type registryCmd struct {
	cobraCmd *cobra.Command
	opts     RegistryCmdOptions
}

type RegistryCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewRegistryCmd(opts *common.CmdRootOptions) *registryCmd {
	var cobraCmd = &cobra.Command{
		Use:   "registry",
		Short: fmt.Sprintf("Create a private docker registry [%v]", common.GlobalOptions.DockerRegistryDns),
		Long:  fmt.Sprintf("Create a private docker registry [%v]", common.GlobalOptions.DockerRegistryDns),
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Deploying Docker Registry")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {
				_ = loadKubeConfig()

				if shouldDeleteRegistry {

					// Remove registry
					if err := uninstallRegistry(); err != nil {
						logger.Fatal(err.Error())
					}
				} else {

					// Deploy registry
					if err := deployDockerRegistry(); err != nil {
						logger.Fatal(err.Error())
					}
				}
			}

			logger.PrintCompletion()
		},
	}

	var registryCmd = new(registryCmd)
	registryCmd.cobraCmd = cobraCmd
	registryCmd.opts.CmdRootOptions = opts

	if err := registryCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return registryCmd
}

func (cmd *registryCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *registryCmd) initFlags() error {
	// TODO: Allow force creation by flag even if registry exists
	cmd.cobraCmd.Flags().BoolVarP(
		&shouldDeleteRegistry,
		"Delete Kubernetes docker registry as a pod",
		"d",
		shouldDeleteRegistry,
		"anchor cluster registry -d")
	return nil
}

func deployDockerRegistry() error {
	if exists, err := checkForActiveRegistry(); err != nil {
		return err
	} else if !exists {

		// Kill possible running kubectl registry port forwarding
		_ = killRegistryPortForwarding()

		nodes, _ := getAllNodes()

		for _, node := range nodes {
			node = strings.TrimPrefix(node, "node/")
			node = strings.TrimSuffix(node, "\n")

			// Overwrite config.toml with our own
			if err := overrideContainerdConfig(node); err != nil {
				return err
			}

			// Restart container runtime
			if err := restartContainerd(node); err != nil {
				return err
			}

			// Restart kubelet service
			if err := restartKubeletService(node); err != nil {
				return err
			}

			// Deploy docker registry as a pod
			if err := deployDockerRegistryPod(); err != nil {
				return err
			}

			// Wait for the registry pod to be ready with 2 minutes timeout
			if err := waitForDockerRegistryPod(); err != nil {
				return err
			}

			// Forwards registry port 32001 -> 5000
			if err := forwardDockerRegistryPort(); err != nil {
				return err
			}
		}
	} else {
		logger.Info("Docker registry already exists, skipping creation.")
		return printRegistryInfo()
	}

	return nil
}

func checkForActiveRegistry() (bool, error) {
	getDashboardCmd := "kubectl get deployments registry --namespace=container-registry"
	if out, err := common.ShellExec.ExecuteWithOutput(getDashboardCmd); err != nil {
		if strings.Contains(out, "NotFound") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getAllNodes() ([]string, error) {
	if out, err := common.ShellExec.ExecuteWithOutput("kubectl get nodes -oname"); err != nil {
		return nil, err
	} else {
		nodes := strings.Split(out, " ")
		return nodes, nil
	}
}

func overrideContainerdConfig(nodeName string) error {
	logger.Info("Overwriting control plane containerd config.toml...")
	if file, err := ioutil.TempFile(os.TempDir(), "anchor-containerd-config-template"); err != nil {
		return err
	} else {
		// Remove after finished
		defer os.Remove(file.Name())

		// Temporary do not allow overriding the docker registry name via ENV var on the config.toml string content
		if _, err := file.WriteString(config.RegistryContainerdConfigTemplate); err != nil {
			return err
		} else {
			replaceConfigCmd := fmt.Sprintf("envsubst < %v > config.toml; docker cp config.toml %v:/etc/containerd/config.toml",
				file.Name(), nodeName)
			return common.ShellExec.Execute(replaceConfigCmd)
		}
	}
}

func restartContainerd(nodeName string) error {
	logger.Info("Restarting control plane containerd...")
	containerdRestartCmd := fmt.Sprintf("docker exec %v systemctl restart containerd", nodeName)
	return common.ShellExec.Execute(containerdRestartCmd)
}

func restartKubeletService(nodeName string) error {
	logger.Info("Restarting control plane kubelet service...")
	kubeletRestartCmd := fmt.Sprintf("docker exec %v systemctl restart kubelet.service", nodeName)
	return common.ShellExec.Execute(kubeletRestartCmd)
}

func deployDockerRegistryPod() error {
	logger.Info("Deploying docker registry as a pod...")
	deployRegistryCmd := fmt.Sprintf("kubectl apply -f %v", common.GlobalOptions.RegistryManifest)
	return common.ShellExec.Execute(deployRegistryCmd)
}

func waitForDockerRegistryPod() error {
	logger.Info("Waiting for docker registry pod to be ready (2m timeout)...")
	waitContainerCmd := fmt.Sprintf("kubectl wait -n container-registry -l app=registry --timeout=2m --for=condition=Ready pod")
	return common.ShellExec.Execute(waitContainerCmd)
}

func forwardDockerRegistryPort() error {
	logger.Info("Port forwarding container registry 32001 --> 5000...")
	_ = killRegistryPortForwarding()
	portFwdCmd := fmt.Sprintf("kubectl port-forward -n container-registry service/registry 32001:5000")
	return common.ShellExec.ExecuteInBackground(portFwdCmd)
}

func killRegistryPortForwarding() error {
	return common.ShellExec.Execute(`ps -ef | grep "kubectl port-forward -n container-registry" | grep -v grep | awk '{print $2}' | xargs kill -9`)
}

func printRegistryInfo() error {
	if exists, err := checkForActiveRegistry(); err != nil {
		return err
	} else if exists {
		logger.Info(`
Registry:
---------`)
		logger.Infof("Registry is available at: %s", "registry.anchor:32001")

		getCatalogCmd := fmt.Sprintf(`docker exec -t anchor-control-plane /bin/sh -c "curl -X GET http://%v/v2/_catalog"`, common.GlobalOptions.DockerRegistryDns)
		if out, err := common.ShellExec.ExecuteWithOutput(getCatalogCmd); err != nil {
			// TODO: Change to warn/error ?
			logger.Info(out)
			return err
		} else {
			logger.Info("\nCatalog:")
			src := []byte(out)

			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, src, "", "  "); err != nil {
				// Do noting
				logger.Info(out)
			} else {
				logger.Info(prettyJSON.String())
			}
		}
	}
	return nil
}

func uninstallRegistry() error {
	if exists, err := checkForActiveRegistry(); err != nil {
		return err
	} else if exists {
		logger.Info("\n==> Uninstalling registry...\n")
		uninstallCmd := fmt.Sprintf("kubectl delete -f %v", common.GlobalOptions.RegistryManifest)
		if err := common.ShellExec.Execute(uninstallCmd); err != nil {
			return err
		}
	} else {
		logger.Info("Docker Registry does not exists, nothing to delete.")
	}
	return nil
}
