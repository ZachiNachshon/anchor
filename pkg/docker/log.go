package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"strings"
)

func LogContainer(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Logging container [%v]", name))
	}

	var dirname = ""
	var err error
	if dirname, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	}

	// Load .env file
	config.LoadEnvVars(identifier)

	logContainerFmt := "docker logs -f %v"

	imageIdentifier := ComposeDockerImageIdentifierNoTag(dirname)
	containerImagesCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if containerImage, err := common.ShellExec.ExecuteWithOutput(containerImagesCmd); err != nil {
		return err
	} else if len(containerImage) > 0 {
		logger.Infof("Logging docker container: %v", imageIdentifier)
		logContainerCmd := fmt.Sprintf(logContainerFmt, containerImage)
		logContainerCmd = strings.Replace(logContainerCmd, "\n", " ", -1)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + logContainerCmd + "\n")
		}

		if err := common.ShellExec.Execute(logContainerCmd); err != nil {
			return err
		}
	} else {
		logger.Infof("No running containers can be found for name: %v", imageIdentifier)
	}

	return nil
}
