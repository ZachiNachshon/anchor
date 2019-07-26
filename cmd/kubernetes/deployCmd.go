package kubernetes

import (
	"fmt"
	"github.com/anchor/config"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/locator"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	cobraCmd *cobra.Command
	opts     DeployCmdOptions
}

type DeployCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewDeployCmd(opts *common.CmdRootOptions) *deployCmd {
	var cobraCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a container Kubernetes manifest",
		Long:  `Deploy a container Kubernetes manifest`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Deploy Container Manifest")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {

				_ = loadKubeConfig()

				if _, err := deployManifest(args[0]); err != nil {
					logger.Fatal(err.Error())
				}
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

func deployManifest(identifier string) (string, error) {
	if manifestFilePath, err := locator.DirLocator.Manifest(identifier); err != nil {
		return "", err
	} else {
		// Load .env file
		config.LoadEnvVars(identifier)

		if common.GlobalOptions.Verbose {
			logManifestCmd := fmt.Sprintf("cat %v | envsubst", manifestFilePath)
			_ = common.ShellExec.Execute(logManifestCmd)
		}

		deployCmd := fmt.Sprintf("envsubst < %v | kubectl apply -f -", manifestFilePath)
		if err := common.ShellExec.Execute(deployCmd); err != nil {
			return "", err
		}
		return manifestFilePath, nil
	}
}
