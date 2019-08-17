package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

func RunContainer(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Running container [%v]", name))
	}

	if runCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandRun); err != nil {
		return err
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + runCmd + "\n")
		}

		if out, err := common.ShellExec.ExecuteWithOutput(runCmd); err != nil {
			logger.Info(out)
			return err
		} else {
			tailCmd := fmt.Sprintf("docker logs -f %v", out)
			if err := common.ShellExec.Execute(tailCmd); err != nil {
				return err
			}
		}
	}
	return nil
}
