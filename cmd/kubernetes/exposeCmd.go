package kubernetes

import (
	"fmt"
	"github.com/anchor/config"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/locator"
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

func enablePortForwarding(dirname string) error {
	l := locator.NewLocator()
	if manifestFilePath, err := l.Manifest(dirname); err != nil {
		return err
	} else {
		manDir := l.GetRootFromManifestFile(manifestFilePath)
		config.LoadEnvVars(manDir)

		if exposeCmd, err := extractManifestCmd(manifestFilePath, ManifestCommandPortForward); err != nil {
			return err
		} else if len(exposeCmd) > 0 {
			logger.Infof("==> Enabling port forwarding for resource %v...", dirname)
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + exposeCmd + "\n")
			}
			if err := common.ShellExec.ExecuteInBackground(exposeCmd); err != nil {
				return err
			}
		} else {
			warnMsg := fmt.Sprintf("Cannot find [%v] declaration, resource won't get exposed on host.\n  ==>  %v ", ManifestCommandPortForward, manifestFilePath)
			logger.Info(warnMsg)
		}

		return nil
	}
}

//func waitForDeployedManifest(manifestFilePath string) error {
//	if waitCmd, err := extractManifestCmd(manifestFilePath, ManifestCommandWait); err != nil {
//		return err
//	} else if len(waitCmd) > 0 {
//		logger.Info("==> Waiting for resource to become ready...")
//		if common.GlobalOptions.Verbose {
//			logger.Info("\n" + waitCmd + "\n")
//		}
//		if err := common.ShellExec.Execute(waitCmd); err != nil {
//			return err
//		}
//	} else {
//		warnMsg := fmt.Sprintf("Cannot find [%v] declaration, won't wait for resource to become ready.\n  ==>  %v ", ManifestCommandWait, manifestFilePath)
//		logger.Info(warnMsg)
//	}
//
//	return nil
//}
