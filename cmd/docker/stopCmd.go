package docker

import (
	"fmt"

	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"

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
		Short: "Stop containers",
		Long:  `Stop containers`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Cleanup: Containers")
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

func stopContainers(dirname string) error {
	stopContainerFmt := "docker stop %v"
	imageIdentifier := composeDockerImageIdentifierNoTag(dirname)

	runningContainerCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if runningContainer, err := common.ShellExec.ExecuteWithOutput(runningContainerCmd); err != nil {
		return err
	} else if len(runningContainer) > 0 {
		logger.Infof("Stopping docker container for name: %v", imageIdentifier)
		containerIds := fmt.Sprintf(stopContainerFmt, runningContainer)
		_ = common.ShellExec.Execute(containerIds)
	} else {
		logger.Infof("No containers are running for name: %v", imageIdentifier)
	}

	return nil
}

func removeContainers(dirname string) error {
	removeContainerFmt := "docker rm -f %v"
	imageIdentifier := composeDockerImageIdentifierNoTag(dirname)

	runningContainerCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if existingContainer, err := common.ShellExec.ExecuteWithOutput(runningContainerCmd); err != nil {
		return err
	} else if len(existingContainer) > 0 {
		logger.Infof("Removing existing container for name: %v", imageIdentifier)
		containerIds := fmt.Sprintf(removeContainerFmt, existingContainer)
		_ = common.ShellExec.Execute(containerIds)
	} else {
		logger.Infof("No existing containers identified for name: %v", imageIdentifier)
	}

	pastContainerCmd := fmt.Sprintf("docker ps -a | grep '%v' | awk {'print $1'}", imageIdentifier)
	if pastContainer, err := common.ShellExec.ExecuteWithOutput(pastContainerCmd); err != nil {
		return err
	} else if len(pastContainer) > 0 {
		logger.Infof("Removing past container for name: %v", imageIdentifier)
		containerIds := fmt.Sprintf(removeContainerFmt, pastContainer)
		_ = common.ShellExec.Execute(containerIds)
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
