package kubernetes

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/input"
	"github.com/spf13/cobra"
)

type destroyCmd struct {
	cobraCmd *cobra.Command
	opts     DestroyCmdOptions
}

type DestroyCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewDestroyCmd(opts *common.CmdRootOptions) *destroyCmd {
	var cobraCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy local Kubernetes cluster",
		Long:  `Destroy local Kubernetes cluster`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Destroy Kubernetes Cluster")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {
				_ = loadKubeConfig()

				if err := destroyKubernetesCluster(name); err != nil {
					logger.Fatal(err.Error())
				}
			}

			logger.PrintCompletion()
		},
	}

	var destroyCmd = new(destroyCmd)
	destroyCmd.cobraCmd = cobraCmd
	destroyCmd.opts.CmdRootOptions = opts

	if err := destroyCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return destroyCmd
}

func (cmd *destroyCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *destroyCmd) initFlags() error {
	return nil
}

func destroyKubernetesCluster(name string) error {
	in := input.NewYesNoInput()
	destroyInputFormat := fmt.Sprintf("Are you sure you want to destroy Kubernetes cluster [%v]?", name)

	if result, err := in.WaitForInput(destroyInputFormat); err != nil || !result {
		logger.Info("skipping.")
	} else {

		// Kill possible running kubectl proxy
		_ = killKubectlProxy()

		// Remove docker registry and clean /etc/host entry
		_ = uninstallRegistry()

		destroyCmd := fmt.Sprintf("kind delete cluster --name %v", name)
		if err := common.ShellExec.Execute(destroyCmd); err != nil {
			return err
		}
	}
	return nil
}
