package kubernetes

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type removeCmd struct {
	cobraCmd *cobra.Command
	opts     RemoveCmdOptions
}

type RemoveCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewRemoveCmd(opts *common.CmdRootOptions) *removeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "remove",
		Short: "Removed a previously deployed container manifest",
		Long:  `Removed a previously deployed container manifest`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := common.GlobalOptions.KindClusterName
			logger.PrintHeadline("Remove Container Manifest")
			_ = loadKubeConfig()
			if err := RemoveManifest(name, args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var removeCmd = new(removeCmd)
	removeCmd.cobraCmd = cobraCmd
	removeCmd.opts.CmdRootOptions = opts

	if err := removeCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return removeCmd
}

func (cmd *removeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *removeCmd) initFlags() error {
	return nil
}

func RemoveManifest(clusterName string, dirname string) error {
	if exists, err := checkForActiveCluster(clusterName); err != nil {
		return err
	} else if exists {
		if manifestPath, err := getContainerManifestsDir(dirname); err != nil {
			return err
		} else {
			removeCmd := fmt.Sprintf("envsubst < %v | kubectl delete -f -", manifestPath)
			if err := common.ShellExec.Execute(removeCmd); err != nil {
				return err
			}
		}
	}
	return nil
}
