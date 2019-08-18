package cluster

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

func Apply(identifier string, namespace string) (string, error) {
	var name = ""
	var err error
	if name, err = locator.DirLocator.Name(identifier); err != nil {
		return "", err
	}

	if manifestFilePath, err := locator.DirLocator.Manifest(name); err != nil {
		return "", err
	} else {

		// Load .env file
		config.LoadEnvVars(identifier)

		// Check if volume should be mounted via hostPath on manifest.yaml
		if stateful, err := extractor.CmdExtractor.ManifestContent(name, extractor.ManifestCommandStateful); err == nil && stateful {
			logger.PrintCommandHeader(fmt.Sprintf("Applying %v (Stateful)", name))
			if err := mountHostPath(name, namespace); err != nil {
				return "", err
			}
		} else {
			logger.PrintCommandHeader(fmt.Sprintf("Applying %v", name))
		}

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
