package kubernetes

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	cobraCmd *cobra.Command
	opts     CreateCmdOptions
}

type DeployCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewDeployCmd(opts *common.CmdRootOptions) *deployCmd {
	var cobraCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy container resource",
		Long:  `Deploy container resource from a directory`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := common.GlobalOptions.KindClusterName
			logger.PrintHeadline("Deploy Container Resource")
			_ = loadKubeConfig()
			if err := deployContainerResource(name, args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var deployCmd = new(deployCmd)
	deployCmd.cobraCmd = cobraCmd
	deployCmd.opts.CmdRootOptions = opts

	if err := deployCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return deployCmd
}

func (cmd *deployCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *deployCmd) initFlags() error {
	return nil
}

func deployContainerResource(clusterName string, dirname string) error {
	if exists, err := checkForActiveCluster(clusterName); err != nil {
		return err
	} else if exists {
		if resfilePath, err := getContainerResourceDir(dirname); err != nil {
			return err
		} else {
			deployCmd := "kubectl apply -f " + resfilePath
			if err := common.ShellExec.Execute(deployCmd); err != nil {
				return err
			}
		}
	}
	return nil
}
