package kubernetes

import (
	"fmt"
	"github.com/anchor/cmd/docker"
	"github.com/anchor/config"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/locator"
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

				if _, err := removeManifest(args[0]); err != nil {
					logger.Fatal(err.Error())
				}
				if err := disablePortForwarding(args[0]); err != nil {
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

func removeManifest(identifier string) (string, error) {
	if manifestFilePath, err := locator.DirLocator.Manifest(identifier); err != nil {
		return "", err
	} else {
		// Load .env file
		config.LoadEnvVars(identifier)

		if common.GlobalOptions.Verbose {
			logManifestCmd := fmt.Sprintf("cat %v | envsubst", manifestFilePath)
			_ = common.ShellExec.Execute(logManifestCmd)
		}

		removeCmd := fmt.Sprintf("envsubst < %v | kubectl delete -f -", manifestFilePath)
		if err := common.ShellExec.Execute(removeCmd); err != nil {
			// Do noting
		}
		return manifestFilePath, nil
	}
}

func disablePortForwarding(dirname string) error {
	identifier := docker.ComposeDockerContainerIdentifierNoTag(dirname)
	killPortFwdCmd := fmt.Sprintf(`ps -ef | grep "%v" | grep -v grep | awk '{print $2}' | xargs kill -9`, identifier)
	if common.GlobalOptions.Verbose {
		logger.Info("\n" + killPortFwdCmd + "\n")
	}
	if err := common.ShellExec.Execute(killPortFwdCmd); err != nil {
		// Do nothing
	}

	return nil
}
