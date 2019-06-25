package kubernetes

import (
	"os"
	"strings"

	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
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
		Long:  `Create a local Kubernetes cluster based on Kind.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Creating Kubernetes Cluster")
			if err := createKubernetesCluster(); err != nil {
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
	cmd.cobraCmd.Flags().StringVarP(
		&common.GlobalOptions.KindClusterName,
		"Kind cluster name",
		"n",
		common.GlobalOptions.KindClusterName,
		"docker image DOCKER_IMAGE_TAG")
	return nil
}

func createKubernetesCluster() error {
	if exists, err := checkForActiveCluster(); err != nil {
		return err
	} else if !exists {
		createCmd := "kind create cluster --name ${CLUSTER_NAME}"
		if err := common.ShellExec.Execute(createCmd); err != nil {
			return err
		}

		// TODO: install dashboard

	} else {
		clusterName := os.Getenv("CLUSTER_NAME")
		logger.Infof("Cluster %v already exists, skipping creation", clusterName)
	}
	return nil
}

func checkForActiveCluster() (bool, error) {
	createCmd := "kind get clusters"
	if out, err := common.ShellExec.ExecuteWithOutput(createCmd); err != nil {
		return false, err
	} else {
		clusterName := os.Getenv("CLUSTER_NAME")
		contains := strings.Contains(out, clusterName)
		return contains, nil
	}
}
