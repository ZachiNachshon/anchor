package kubernetes

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
	"runtime"
)

type tokenCmd struct {
	cobraCmd *cobra.Command
	opts     ConnectCmdOptions
}

type TokenCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewTokenCmd(opts *common.CmdRootOptions) *tokenCmd {
	var cobraCmd = &cobra.Command{
		Use:   "token",
		Short: "Generate export KUBECONFIG command and load to clipboard",
		Long:  `Generate export KUBECONFIG command and load to clipboard`,
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Generate export KUBECONFIG command")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {
				_ = loadKubeConfig()

				exportCmd := fmt.Sprintf(`export KUBECONFIG="$(kind get kubeconfig-path --name=%s)"`, common.GlobalOptions.KindClusterName)
				_ = loadToClipboard(exportCmd)

				kubeConfig := fmt.Sprintf(`
Paste from clipboard the following script to apply kube config on existing session:

 %v`, exportCmd)

				logger.Info(kubeConfig)
			}

			logger.PrintCompletion()
		},
	}

	var tokenCmd = new(tokenCmd)
	tokenCmd.cobraCmd = cobraCmd
	tokenCmd.opts.CmdRootOptions = opts

	if err := tokenCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return tokenCmd
}

func (cmd *tokenCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *tokenCmd) initFlags() error {
	return nil
}

func loadToClipboard(cmd string) error {
	switch runtime.GOOS {
	case "darwin":
		{
			if err := common.ShellExec.Execute(fmt.Sprintf("echo %v | pbcopy", cmd)); err != nil {
				logger.Info("Failed setting value to clipboard using 'pbcopy ...'")
				return err
			}
			break
		}
	case "linux":
		{
			if err := common.ShellExec.Execute(fmt.Sprintf("xclip -selection %v", cmd)); err != nil {
				logger.Info("Failed setting value to clipboard using 'xclip -selection ...'")
				return err
			}
			break
		}
	}
	return nil
}
