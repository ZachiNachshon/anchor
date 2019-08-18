package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"strings"
)

func Tag(identifier string) error {
	logger.PrintCommandHeader(fmt.Sprintf("Tagging image %v", identifier))
	if tagCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandTag); err != nil {
		return err
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + tagCmd + "\n")
		}
		if err = common.ShellExec.Execute(tagCmd); err != nil {
			return err
		}
		logger.Info("Successfully tagged.")
		return nil
	}
}

func Push(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Pushing image %v", name))
	}

	if pushCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandPush); err != nil {
		return err
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + pushCmd + "\n")
		}
		if err = common.ShellExec.Execute(pushCmd); err != nil {
			return err
		}

		_ = untagDockerImage(pushCmd)

		logger.Info("Successfully pushed to registry.")
	}
	return nil
}

func untagDockerImage(pushCommand string) error {
	removeImage := strings.Replace(pushCommand, "push", "rmi -f", 1)
	if err := common.ShellExec.Execute(removeImage); err != nil {
		return err
	}
	return nil
}
