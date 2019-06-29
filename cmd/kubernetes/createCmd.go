package kubernetes

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
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
			name := common.GlobalOptions.KindClusterName
			logger.PrintHeadline("Creating Kubernetes Cluster")
			if err := createKubernetesCluster(name); err != nil {
				logger.Fatal(err.Error())
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
	if exists, err := checkForActiveCluster(name); err != nil {
		return err
	} else if !exists {
		createCmd := "kind create cluster --name " + name
		if err := common.ShellExec.Execute(createCmd); err != nil {
			return err
		}

		_ = loadKubeConfig()

		//_ = deployKubernetesDashboard()

		//_ = deployDockerRegistry()

		_ = printClusterStatus(name)

	} else {
		logger.Infof("Cluster %v already exists, skipping creation", name)
	}
	return nil
}
