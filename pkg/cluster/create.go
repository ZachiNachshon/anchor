package cluster

import (
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

func createNamespace() error {
	var namespace = common.GlobalOptions.DockerImageNamespace
	logger.PrintCommandHeader(fmt.Sprintf("Creating %v namespace", namespace))
	namespaceManifest := strings.Replace(config.KubernetesNamespaceManifest, "NAMESPACE-TO-REPLACE", namespace, 1)

	if file, err := ioutil.TempFile(os.TempDir(), "anchor-namespace-manifest.yaml"); err != nil {
		return err
	} else {
		// Remove after finished
		defer os.Remove(file.Name())

		if _, err := file.WriteString(namespaceManifest); err != nil {
			return err
		} else {
			createNamespaceCmd := fmt.Sprintf("cat %v | kubectl apply -f -",
				file.Name())
			return common.ShellExec.Execute(createNamespaceCmd)
		}
	}
}

// After restart Kind container nodes are not running, they should be started again
func startKindContainerNodesIfNeeded() (bool, error) {
	kindNodesIdentifier := "kindest/node:"

	runningContainerCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", kindNodesIdentifier)
	if runningContainer, err := common.ShellExec.ExecuteWithOutput(runningContainerCmd); err != nil {
		return false, err
	} else if len(runningContainer) > 0 {
		// Do nothing, already running
		return false, nil
	}

	pastContainerCmd := fmt.Sprintf("docker ps -a | grep '%v' | awk {'print $1'}", kindNodesIdentifier)
	if pastContainers, err := common.ShellExec.ExecuteWithOutput(pastContainerCmd); err != nil {
		return false, err
	} else if len(pastContainers) > 0 {
		logger.Infof("\nFound %v stopped Kind containers, starting and sleeping for 15s...", len(pastContainers))
		pastContainersArr := strings.Split(pastContainers, "\n")
		for _, containerId := range pastContainersArr {
			if len(containerId) <= 0 {
				continue
			}
			logger.Info(fmt.Sprintf("Starting Kind node %v...", containerId))
			startNodeCmd := fmt.Sprintf("docker start %v", containerId)
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + startNodeCmd + "\n")
			}
			if err := common.ShellExec.Execute(startNodeCmd); err != nil {
				logger.Info(fmt.Sprintf("Failed command: %v. Error: %v", startNodeCmd, err.Error()))
			}
		}

		// Sleep for 15
		time.Sleep(15 * time.Second)
		return true, nil
	}

	return false, nil
}

func createKindConfigManifest(workerNodes int) string {
	var workerNodesStr = ""
	for i := 0; i < workerNodes; i++ {
		workerNodesStr += "- role: worker\n"
	}
	kindClusterManifest := fmt.Sprintf(config.KindClusterManifestFormat, "- role: control-plane", workerNodesStr)
	return kindClusterManifest
}

func Create() error {
	name := common.GlobalOptions.KindClusterName

	logger.Info("\nPlease specify how many worker nodes should be created (default 1):\n")
	in := input.NewNumericInput()
	var v int
	var err error
	if v, err = in.WaitForInputAllowDefault(); err != nil {
		return err
	} else if v == -1 {
		v = 1
	}

	if file, err := ioutil.TempFile(os.TempDir(), "anchor-kind-cluster-manifest.yaml"); err != nil {
		return err
	} else {
		// Remove after finished
		defer os.Remove(file.Name())

		manifest := createKindConfigManifest(v)
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + manifest)
		}

		if _, err := file.WriteString(manifest); err != nil {
			return err
		} else {
			logger.PrintCommandHeader(fmt.Sprintf("Creating cluster %v with %v worker node(s)", name, v))
			createCmdFormat := "kind create cluster --name %v --config %v"
			createCmd := fmt.Sprintf(createCmdFormat, name, file.Name())
			if err := common.ShellExec.Execute(createCmd); err != nil {
				return err
			}

			_ = LoadKubeConfig()
			_ = createNamespace()
			_ = Dashboard()
			_ = Registry()
			_ = Status()

			return nil
		}
	}

	return nil
}
