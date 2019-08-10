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
)

type createCmd struct {
	cobraCmd *cobra.Command
	opts     CreateCmdOptions
}

type CreateCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCreateCmd(opts *common.CmdRootOptions) *createCmd {
	var cobraCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a local Kubernetes cluster",
		Long:  `Create a local Kubernetes cluster`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Creating Kubernetes Cluster")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if exists {
				logger.Infof("Cluster %v already exists, skipping creation", name)
			} else {
				if err := createKubernetesCluster(name); err != nil {
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

func createKubernetesCluster(name string) error {
	createCmd := "kind create cluster --name " + name
	if err := common.ShellExec.Execute(createCmd); err != nil {
		return err
	}

	_ = loadKubeConfig()
	_ = createNamespace()
	_ = deployKubernetesDashboard()
	_ = deployDockerRegistry()
	_ = printClusterStatus(name)

	return nil
}

func createNamespace() error {
	var namespace = common.GlobalOptions.DockerImageNamespace
	logger.Infof("\nCreating %v namespace...", namespace)

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
