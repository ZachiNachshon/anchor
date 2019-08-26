package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

func Expose(identifier string) error {
	var name string
	var err error

	if name, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Port forwarding %v", name))
	}

	if exposeCmd, err := extractor.CmdExtractor.ManifestCmd(name, extractor.ManifestCommandPortForward); err != nil {
		return err
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + exposeCmd + "\n")
		}

		label := fmt.Sprintf("%v=%v", "app", name)
		if ready, err := waitForPodReadiness(label, common.GlobalOptions.DockerImageNamespace); err != nil {
			// TODO: handle error gracefully
			logger.Info(err.Error())
		} else if !ready {
			logger.Infof("Cannot identify ready pods with label %v, skipping port forwarding, try 'expose' cluster command once pod is ready", label)
		} else {
			if err := common.ShellExec.ExecuteInBackground(exposeCmd); err != nil {
				return err
			}
		}
	}
	return nil
}
