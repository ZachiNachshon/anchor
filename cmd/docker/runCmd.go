package docker

import (
	"fmt"

	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
	"github.com/spf13/cobra"
)

type runCmd struct {
	cobraCmd *cobra.Command
	opts     RunCmdOptions
}

type RunCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewRunCmd(opts *common.CmdRootOptions) *runCmd {
	var cobraCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a docker container based on Dockerfile instructions",
		Long:  `Run a Dockerfile from the list of available docker image directories (ex. DIR=alpine).`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Deploying Docker Container on LOCAL Machine")
			if err := stopContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := removeContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := runContainer(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var runCmd = new(runCmd)
	runCmd.cobraCmd = cobraCmd
	runCmd.opts.CmdRootOptions = opts

	if err := runCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return runCmd
}

func runContainer(dirname string) error {
	if dockerfilePath, err := getDockerfileContextPath(dirname); err != nil {
		return err
	} else {
		if runCmd, err := extractDockerCmd(dockerfilePath, DockerCommandRun); err != nil {
			return err
		} else {
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + runCmd)
			}

			if containerId, err := common.ShellExec.ExecuteWithOutput(runCmd); err != nil {
				return err
			} else {
				tailCmd := fmt.Sprintf("docker logs -f %v", containerId)
				if err := common.ShellExec.Execute(tailCmd); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (cmd *runCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *runCmd) initFlags() error {
	return nil
}
