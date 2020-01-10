package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"strings"
)

func RemoveImages(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Removing image [%v]", name))
	}

	var dirname = ""
	var err error
	if dirname, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	}

	// Load .env file
	config.LoadEnvVars(identifier)

	removeImagesFmt := "docker rmi -f %v"

	unknownImagesCmd := "docker images | grep '<none>' | awk {'print $3'}"
	if unknownImages, err := common.ShellExec.ExecuteWithOutput(unknownImagesCmd); err != nil {
		return err
	} else if len(unknownImages) > 0 {
		logger.Info("Removing docker images for name: <none>")
		removeUnknownCmd := fmt.Sprintf(removeImagesFmt, unknownImages)
		removeUnknownCmd = strings.Replace(removeUnknownCmd, "\n", " ", -1)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + removeUnknownCmd + "\n")
		}
		_ = common.ShellExec.Execute(removeUnknownCmd)

	} else {
		logger.Info("No images can be found for name: <none>")
	}

	imageIdentifier := ComposeDockerImageIdentifierNoTag(dirname)
	containerImagesCmd := fmt.Sprintf("docker images | grep '%v' | awk {'print $3'}", imageIdentifier)
	if containerImages, err := common.ShellExec.ExecuteWithOutput(containerImagesCmd); err != nil {
		return err
	} else if len(containerImages) > 0 {
		logger.Infof("Removing docker images for name: %v", imageIdentifier)
		removeImageCmd := fmt.Sprintf(removeImagesFmt, containerImages)
		removeImageCmd = strings.Replace(removeImageCmd, "\n", " ", -1)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + removeImageCmd + "\n")
		}
		_ = common.ShellExec.Execute(removeImageCmd)

	} else {
		logger.Infof("No images can be found for name: %v", imageIdentifier)
	}

	return nil
}
