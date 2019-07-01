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
			logger.PrintHeadline("Remove Container Manifest")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {
				_ = loadKubeConfig()

				if err := removeManifest(args[0]); err != nil {
					logger.Fatal(err.Error())
				}
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

func removeManifest(dirname string) error {
	if manifestPath, err := getContainerManifestsDir(dirname); err != nil {
		return err
	} else {
		removeCmd := fmt.Sprintf("envsubst < %v | kubectl delete -f -", manifestPath)
		if err := common.ShellExec.Execute(removeCmd); err != nil {
			return err
		}
	}
	return nil
}
