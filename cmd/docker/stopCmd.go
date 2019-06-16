package docker

import (
	"fmt"
	"github.com/kit/cmd/logger"
	"github.com/kit/cmd/types"
	"github.com/kit/cmd/utils"

	"github.com/spf13/cobra"
)

type StopCmd struct {
	cobraCmd *cobra.Command
	opts     StopCmdOptions
}

type StopCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewStopCmd(opts *common.CmdRootOptions) *StopCmd {
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

	var stopCmd = new(StopCmd)
	stopCmd.cobraCmd = cobraCmd
	stopCmd.opts.CmdRootOptions = opts

	if err := stopCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return stopCmd
}

func stopContainers(dirname string) error {
	stopContainerFmt := "docker stop %v"
	imageIdentifier := composeDockerImageIdentifierNoSuite(dirname)

	runningContainerCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if runningContainer, err := utils.ExecShellWithOutput(runningContainerCmd); err != nil {
		return err
	} else if len(runningContainer) > 0 {
		logger.Infof("Stopping docker container for name: %v", imageIdentifier)
		containerIds := fmt.Sprintf(stopContainerFmt, runningContainer)
		utils.ExecShell(containerIds)
	} else {
		logger.Infof("No containers are running for name: %v", imageIdentifier)
	}

	return nil
}

func removeContainers(dirname string) error {
	removeContainerFmt := "docker rm -f %v"
	imageIdentifier := composeDockerImageIdentifierNoSuite(dirname)

	runningContainerCmd := fmt.Sprintf("docker ps | grep '%v' | awk {'print $1'}", imageIdentifier)
	if existingContainer, err := utils.ExecShellWithOutput(runningContainerCmd); err != nil {
		return err
	} else if len(existingContainer) > 0 {
		logger.Infof("Removing existing container for name: %v", imageIdentifier)
		containerIds := fmt.Sprintf(removeContainerFmt, existingContainer)
		utils.ExecShell(containerIds)
	} else {
		logger.Infof("No existing containers identified for name: %v", imageIdentifier)
	}

	pastContainerCmd := fmt.Sprintf("docker ps -a | grep '%v' | awk {'print $1'}", imageIdentifier)
	if pastContainer, err := utils.ExecShellWithOutput(pastContainerCmd); err != nil {
		return err
	} else if len(pastContainer) > 0 {
		logger.Infof("Removing past container for name: %v", imageIdentifier)
		containerIds := fmt.Sprintf(removeContainerFmt, pastContainer)
		utils.ExecShell(containerIds)
	} else {
		logger.Infof("No past containers identified for name: %v", imageIdentifier)
	}

	return nil
}

func (cmd *StopCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *StopCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}
