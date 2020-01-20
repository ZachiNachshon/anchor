package token

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/clipboard"
	"github.com/spf13/cobra"
)

type tokenCmd struct {
	cobraCmd *cobra.Command
	opts     TokenOptions
}

type TokenOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *tokenCmd {
	var cobraCmd = &cobra.Command{
		Use:   "token",
		Short: "Generate export KUBECONFIG command and load to clipboard",
		Long:  `Generate export KUBECONFIG command and load to clipboard`,
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cluster.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.ClusterHeadline, "Token")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			exportCmd := fmt.Sprintf(`export KUBECONFIG="$(kind get kubeconfig-path --name=%s)"`, common.GlobalOptions.KindClusterName)
			_ = clipboard.Load(exportCmd)

			kubeConfig := fmt.Sprintf(`
Paste from clipboard the following script to apply kube config on existing session:

 %v`, exportCmd)

			logger.Info(kubeConfig)
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
