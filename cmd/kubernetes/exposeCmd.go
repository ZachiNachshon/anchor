package kubernetes

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/extractor"
	"github.com/spf13/cobra"
)

type exposeCmd struct {
	cobraCmd *cobra.Command
	opts     ExposeCmdOptions
}

type ExposeCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewExposeCmd(opts *common.CmdRootOptions) *exposeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "expose",
		Short: "Expose a container port to the host instance",
		Long:  `Expose a container port to the host instance`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Expose Container Port")
			name := common.GlobalOptions.KindClusterName

			if exists, err := checkForActiveCluster(name); err != nil {
				logger.Fatal(err.Error())
			} else if !exists {
				logger.Info("No active cluster.")
			} else {

				_ = loadKubeConfig()

				if err := enablePortForwarding(args[0]); err != nil {
					logger.Fatal(err.Error())
				}
			}

			logger.PrintCompletion()
		},
	}

	var exposeCmd = new(exposeCmd)
	exposeCmd.cobraCmd = cobraCmd
	exposeCmd.opts.CmdRootOptions = opts

	if err := exposeCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}
	return exposeCmd
}

func (cmd *exposeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *exposeCmd) initFlags() error {
	return nil
}

func enablePortForwarding(identifier string) error {
	if exposeCmd, err := extractor.CmdExtractor.ManifestCmd(identifier, extractor.ManifestCommandPortForward); err != nil {
		return err
	} else {
		logger.Infof("==> Enabling port forwarding for resource %v...", identifier)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + exposeCmd + "\n")
		}

		if err := common.ShellExec.ExecuteInBackground(exposeCmd); err != nil {
			return err
		}
	}
	return nil
}
