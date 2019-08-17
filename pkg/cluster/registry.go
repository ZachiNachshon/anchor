package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

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
		nodes := strings.Split(out, "\n")
		return nodes, nil
	}
}

func overrideContainerdConfig(nodeName string) error {
	logger.Infof("\nOverwriting containerd config.toml on %v...", nodeName)
	var tempDir = os.TempDir()
	if file, err := ioutil.TempFile(tempDir, "anchor-containerd-config-template"); err != nil {
		return err
	} else {
		// Remove after finished
		defer os.Remove(file.Name())

		// Temporary do not allow overriding the docker registry name via ENV var on the config.toml string content
		if _, err := file.WriteString(config.RegistryContainerdConfigTemplate); err != nil {
			return err
		} else {
			replaceConfigCmd := fmt.Sprintf("envsubst < %v > %v/config.toml; docker cp %v/config.toml %v:/etc/containerd/config.toml",
				file.Name(), tempDir, tempDir, nodeName)

			if common.GlobalOptions.Verbose {
				logger.Info("\n" + replaceConfigCmd + "\n")
			}

			return common.ShellExec.Execute(replaceConfigCmd)
		}
	}
}

func restartContainerd(nodeName string) error {
	logger.Infof("Restarting containerd on %v...", nodeName)
	containerdRestartCmd := fmt.Sprintf("docker exec %v systemctl restart containerd", nodeName)
	return common.ShellExec.Execute(containerdRestartCmd)
}

func restartKubeletService(nodeName string) error {
	logger.Infof("Restarting kubelet service on %v...", nodeName)
	kubeletRestartCmd := fmt.Sprintf("docker exec %v systemctl restart kubelet.service", nodeName)
	return common.ShellExec.Execute(kubeletRestartCmd)
}

func installDockerRegistryPod() error {
	logger.PrintCommandHeader("Deploying Registry")
	if file, err := ioutil.TempFile(os.TempDir(), "anchor-registry-manifest.yaml"); err != nil {
		return err
	} else {
		// Remove after finished
		defer os.Remove(file.Name())

		if _, err := file.WriteString(config.KubernetesRegistryManifest); err != nil {
			return err
		} else {
			createDashboardCmd := fmt.Sprintf("cat %v | kubectl apply -f -",
				file.Name())

			if err := common.ShellExec.Execute(createDashboardCmd); err != nil {
				return err
			}
		}
	}
	return nil
}

func waitForDockerRegistryPod() error {
	logger.Info("Waiting for docker registry pod to become ready (5m timeout)...")
	waitContainerCmd := fmt.Sprintf("kubectl wait -n container-registry -l app=registry --timeout=5m --for=condition=Ready pod")
	return common.ShellExec.Execute(waitContainerCmd)
}

func forwardDockerRegistryPort() error {
	logger.Info("Port forwarding container registry 32001 --> 5000...")
	portFwdCmd := fmt.Sprintf("kubectl port-forward -n container-registry service/registry 32001:5000")
	return common.ShellExec.ExecuteInBackground(portFwdCmd)
}

func killRegistryPortForwarding() error {
	return common.ShellExec.Execute(`ps -ef | grep "kubectl port-forward -n container-registry" | grep -v grep | awk '{print $2}' | xargs kill -9`)
}

func printCatalogContent(catalog string) error {
	c := RegistryCatalog{}
	if err := json.Unmarshal([]byte(catalog), &c); err != nil {
		return err
	}

	getTagsFormat := `docker exec -t anchor-control-plane /bin/sh -c "curl -X GET http://registry.anchor:32001/v2/%v/tags/list"`

	result := ExportedCatalog{}
	for _, image := range c.Repositories {

		info := ImageInfo{}
		getTagsCmd := fmt.Sprintf(getTagsFormat, image)
		if tags, err := common.ShellExec.ExecuteWithOutput(getTagsCmd); err != nil {
			logger.Infof("Could not identify tags list for %v", image)
		} else {
			if err := json.Unmarshal([]byte(tags), &info); err == nil {
				result.Images = append(result.Images, info)
			}
		}
	}

	if catalogJson, err := json.Marshal(result); err == nil {
		logger.Info("\nCatalog:")
		src := []byte(catalogJson)
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, src, "", "  "); err != nil {
			// Do noting
			logger.Info(catalog)
			return err
		} else {
			logger.Info(prettyJSON.String())
		}
	}
	return nil
}

func createHostRegistryDns() error {
	logger.PrintCommandHeader(fmt.Sprintf("Checking /etc/hosts for %v", common.GlobalOptions.DockerRegistryDns))
	if err := validateHostsFile(); err != nil {
		return err
	}

	verifyDnsCmd := fmt.Sprintf("hostess has %v", common.GlobalOptions.DockerRegistryDns)
	if err := common.ShellExec.Execute(verifyDnsCmd); err == nil {
		logger.Infof("Found %v on /etc/hosts, no need to add", common.GlobalOptions.DockerRegistryDns)
		return nil
	}

	logger.PrintCommandHeader(fmt.Sprintf("Adding %v to /etc/hosts", common.GlobalOptions.DockerRegistryDns))
	addDnsCmd := fmt.Sprintf("sudo hostess add %v 127.0.0.1", common.GlobalOptions.DockerRegistryDns)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + addDnsCmd + "\n")
	}

	if err := common.ShellExec.Execute(addDnsCmd); err != nil {
		logger.Fatalf("Failed to create DNS record %v on /etc/hosts, cannot start private docker registry", common.GlobalOptions.DockerRegistryDns)
		return err
	}
	return nil
}

func removeHostRegistryDns() error {
	logger.PrintCommandHeader(fmt.Sprintf("Checking /etc/hosts for %v", common.GlobalOptions.DockerRegistryDns))
	if err := validateHostsFile(); err != nil {
		return err
	}

	verifyDnsCmd := fmt.Sprintf("hostess has %v", common.GlobalOptions.DockerRegistryDns)
	if err := common.ShellExec.Execute(verifyDnsCmd); err != nil {
		logger.Infof("Cannot find %v on /etc/hosts, no need to remove", common.GlobalOptions.DockerRegistryDns)
		return nil
	}
	logger.PrintCommandHeader(fmt.Sprintf("Removing %v from /etc/hosts", common.GlobalOptions.DockerRegistryDns))
	removeDnsCmd := fmt.Sprintf("sudo hostess del %v", common.GlobalOptions.DockerRegistryDns)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + removeDnsCmd + "\n")
	}

	if err := common.ShellExec.Execute(removeDnsCmd); err != nil {
		logger.Infof("Failed to remove DNS record %v from /etc/hosts, consider removing manually", common.GlobalOptions.DockerRegistryDns)
		return err
	}
	return nil
}

func validateHostsFile() error {
	if err := common.ShellExec.Execute("hostess fixed"); err != nil {
		ynInput := input.NewYesNoInput()
		q := "Found issues on /etc/hosts file, attempt to fix?"
		if result, err := ynInput.WaitForInput(q); err == nil && result {
			if err := common.ShellExec.Execute("sudo hostess fix"); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

type RegistryCatalog struct {
	Repositories []string `json:"repositories,omitempty"`
}

type ExportedCatalog struct {
	Images []ImageInfo
}

type ImageInfo struct {
	Name string
	Tags []string
}

func PrintRegistryInfo() error {
	if exists, err := checkForActiveRegistry(); err != nil {
		return err
	} else if exists {
		logger.Info(`
Registry:
---------`)
		logger.Infof("Registry is available at: %s", "registry.anchor:32001")

		getCatalogCmd := fmt.Sprintf(`docker exec -t anchor-control-plane /bin/sh -c "curl -X GET http://%v/v2/_catalog"`, common.GlobalOptions.DockerRegistryDnsWithIp)
		if out, err := common.ShellExec.ExecuteWithOutput(getCatalogCmd); err != nil {
			// TODO: Change to warn/error ?
			logger.Info(out)
			return err
		} else {
			_ = printCatalogContent(out)
		}
	}
	return nil
}

func DeleteRegistry() error {
	logger.PrintCommandHeader("Deleting registry")
	if exists, err := checkForActiveRegistry(); err != nil {
		return err
	} else if exists {
		if file, err := ioutil.TempFile(os.TempDir(), "anchor-registry-manifest.yaml"); err != nil {
			return err
		} else {
			// Remove after finished
			defer os.Remove(file.Name())

			if _, err := file.WriteString(config.KubernetesRegistryManifest); err != nil {
				return err
			} else {
				uninstallCmd := fmt.Sprintf("cat %v | kubectl delete -f -",
					file.Name())

				if err := common.ShellExec.Execute(uninstallCmd); err != nil {
					return err
				}
			}
		}

		if err := removeHostRegistryDns(); err != nil {
			return err
		}
	} else {
		logger.Info("Docker Registry does not exists, nothing to delete.")
	}
	return nil
}

func Registry() error {
	logger.PrintCommandHeader("Creating registry")
	if exists, err := checkForActiveRegistry(); err != nil {
		return err
	} else if !exists {

		// Kill possible running kubectl registry port forwarding
		_ = killRegistryPortForwarding()

		nodes, _ := getAllNodes()

		// Change containerd configuration on all cluster nodes
		for _, node := range nodes {
			if len(node) == 0 || node == "" {
				continue
			}

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
		}
	} else {
		logger.Info("Docker registry already exists, skipping creation.")
		return PrintRegistryInfo()
	}

	// Deploy docker registry as a pod
	if err := installDockerRegistryPod(); err != nil {
		return err
	}

	// Sleep for 3 secs to allow registry pod deployment
	time.Sleep(3 * time.Second)

	// Wait for the registry pod to become ready with 2 minutes timeout
	if err := waitForDockerRegistryPod(); err != nil {
		return err
	}

	// Forwards registry port 32001 -> 5000
	if err := forwardDockerRegistryPort(); err != nil {
		return err
	}

	// Create /etc/hosts entry with private docker registry DNS record
	if err := createHostRegistryDns(); err != nil {
		return err
	}

	return nil
}
