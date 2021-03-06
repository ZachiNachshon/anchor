package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

func Delete(identifier string, namespace string) (string, error) {
	var name string
	var err error

	if name, err = locator.DirLocator.Name(identifier); err != nil {
		return "", err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Deleting %v", name))
	}

	if manifestFilePath, err := locator.DirLocator.Manifest(name); err != nil {
		return "", err
	} else {

		// Load .env file
		config.LoadEnvVars(name)

		if common.GlobalOptions.Verbose {
			logManifestCmd := fmt.Sprintf("cat %v | envsubst", manifestFilePath)
			_ = common.ShellExec.Execute(logManifestCmd)
		}

		removeCmd := fmt.Sprintf("envsubst < %v | kubectl delete -f -", manifestFilePath)
		if err := common.ShellExec.Execute(removeCmd); err != nil {
			// Do noting
		}

		// Check if volume should be unmounted via hostPath on manifest.yaml
		if stateful, err := extractor.CmdExtractor.ManifestContent(name, extractor.ManifestCommandStateful); err == nil && stateful {
			if err := unMountHostPath(name, namespace, false); err != nil {
				return "", err
			}
		}

		logger.PrintCommandHeader(fmt.Sprintf("Stopping kubectl process for %v", name))
		if err := KillRunningKubectl(name); err != nil {
			msg := fmt.Sprintf("Failed stopping kubectl process for %v, please do so manually", name)
			logger.Info(msg)
		}

		return manifestFilePath, nil
	}
}

func DisablePortForwarding(dirname string) error {
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
