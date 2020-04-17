package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func ConnectContainer(identifier string) error {
	var dirname = ""
	var err error
	if dirname, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	}

	// Load .env file
	config.LoadEnvVars(identifier)

	connectContainerFmt := "docker exec -it"
	imageIdentifier := ComposeDockerImageIdentifierNoTag(dirname)

	containerImagesCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if containerImage, err := common.ShellExec.ExecuteWithOutput(containerImagesCmd); err != nil {
		return err
	} else if len(containerImage) > 0 {
		execContainerBashCmd := fmt.Sprintf("%v %v %v", connectContainerFmt, containerImage, shell.BASH)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + execContainerBashCmd + "\n")
		}

		logger.PrintCommandHeader(fmt.Sprintf("Connecting to %v (%v) ", dirname, shell.BASH))
		if err := common.ShellExec.ExecuteTTY(execContainerBashCmd); err != nil {

			// Fallback to /bin/sh if /bin/bash is not available
			execContainerShCmd := fmt.Sprintf("%v %v %v", connectContainerFmt, containerImage, shell.SH)
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + execContainerShCmd + "\n")
			}

			logger.PrintCommandHeader(fmt.Sprintf("Connecting to %v (%v)", dirname, shell.SH))
			if err := common.ShellExec.ExecuteTTY(execContainerShCmd); err != nil {
				return err
			}
		}

	} else {
		logger.Infof("No running containers can be found for name: %v", imageIdentifier)
	}

	return nil
}
