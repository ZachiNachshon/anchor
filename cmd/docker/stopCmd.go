package docker

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/locator"

	"github.com/spf13/cobra"
)

type stopCmd struct {
	cobraCmd *cobra.Command
	opts     StopCmdOptions
}

type StopCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewStopCmd(opts *common.CmdRootOptions) *stopCmd {
	var cobraCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop a docker container",
		Long:  `Stop a docker container`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Stop: Docker Container")

			if err := stopContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := removeContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var stopCmd = new(stopCmd)
	stopCmd.cobraCmd = cobraCmd
	stopCmd.opts.CmdRootOptions = opts

	if err := stopCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return stopCmd
}

func stopContainers(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Stopping container [%v]", name))
	}

	var dirname = ""
	var err error
	if dirname, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	}

	stopContainerFmt := "docker stop %v"
	imageIdentifier := ComposeDockerContainerIdentifierNoTag(dirname)

	runningContainerCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if runningContainer, err := common.ShellExec.ExecuteWithOutput(runningContainerCmd); err != nil {
		return err
	} else if len(runningContainer) > 0 {
		logger.Infof("Stopping docker container for name: %v", imageIdentifier)
		stopRunningContainersCmd := fmt.Sprintf(stopContainerFmt, runningContainer)
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + stopRunningContainersCmd + "\n")
		}
		_ = common.ShellExec.Execute(stopRunningContainersCmd)
	} else {
		logger.Infof("No containers are running for name: %v", imageIdentifier)
	}

	return nil
}

func removeContainers(identifier string) error {
	logger.PrintCommandHeader(fmt.Sprintf("Removing container %v", identifier))
	var dirname = ""
	var err error
	if dirname, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	}

	removeContainerFmt := "docker rm -f %v"
	imageIdentifier := ComposeDockerContainerIdentifierNoTag(dirname)

	runningContainerCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if existingContainer, err := common.ShellExec.ExecuteWithOutput(runningContainerCmd); err != nil {
		return err
	} else if len(existingContainer) > 0 {
		logger.Infof("Removing existing container for name: %v", imageIdentifier)
		removeExistingContainersCmd := fmt.Sprintf(removeContainerFmt, existingContainer)
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + removeExistingContainersCmd + "\n")
		}
		_ = common.ShellExec.Execute(removeExistingContainersCmd)
	} else {
		logger.Infof("No existing containers identified for name: %v", imageIdentifier)
	}

	pastContainerCmd := fmt.Sprintf("docker ps -a | grep '%v' | awk {'print $1'}", imageIdentifier)
	if pastContainer, err := common.ShellExec.ExecuteWithOutput(pastContainerCmd); err != nil {
		return err
	} else if len(pastContainer) > 0 {
		logger.Infof("Removing past container for name: %v", imageIdentifier)
		removePastContainersCmd := fmt.Sprintf(removeContainerFmt, pastContainer)
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + removePastContainersCmd + "\n")
		}
		_ = common.ShellExec.Execute(removePastContainersCmd)
	} else {
		logger.Infof("No past containers identified for name: %v", imageIdentifier)
	}

	return nil
}

func (cmd *stopCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *stopCmd) initFlags() error {
	return nil
}
