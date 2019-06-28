package kubernetes

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type deleteCmd struct {
	cobraCmd *cobra.Command
	opts     DeleteCmdOptions
}

type DeleteCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewDeleteCmd(opts *common.CmdRootOptions) *deleteCmd {
	var cobraCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete local Kubernetes cluster",
		Long:  `Delete local Kubernetes cluster`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Delete Kubernetes Cluster")

			_ = loadKubeConfig()

			// Kill possible running kubectl proxy
			_ = killKubectlProxy()

			name := common.GlobalOptions.KindClusterName
			if err := deleteKubernetesCluster(name); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var deleteCmd = new(deleteCmd)
	deleteCmd.cobraCmd = cobraCmd
	deleteCmd.opts.CmdRootOptions = opts

	if err := deleteCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return deleteCmd
}

func (cmd *deleteCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *deleteCmd) initFlags() error {
	return nil
}

func deleteKubernetesCluster(name string) error {
	if exists, err := checkForActiveCluster(name); err != nil {
		return err
	} else if exists {
		deleteCmd := "kind delete cluster --name " + name
		if err := common.ShellExec.Execute(deleteCmd); err != nil {
			return err
		}
	} else {
		logger.Infof("Cluster %v does not exist, skipping deletion", name)
	}
	return nil
}
