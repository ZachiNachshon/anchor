package docker

import (
	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
	"github.com/spf13/cobra"
)

type buildCmd struct {
	cobraCmd *cobra.Command
	opts     BuildCmdOptions
}

type BuildCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewBuildCmd(opts *common.CmdRootOptions) *buildCmd {
	var cobraCmd = &cobra.Command{
		Use:   "build",
		Short: "Builds a Dockerfile",
		Long:  `Builds a docker image from the DOCKER_FILES repository.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Building Docker Image")
			if err := buildDockerfile(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var buildCmd = new(buildCmd)
	buildCmd.cobraCmd = cobraCmd
	buildCmd.opts.CmdRootOptions = opts

	if err := buildCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return buildCmd
}

func buildDockerfile(dirname string) error {
	if dockerfilePath, err := getDockerfileContextPath(dirname); err != nil {
		return err
	} else {
		if buildCmd, err := extractDockerCmd(dockerfilePath, DockerCommandBuild); err != nil {
			return err
		} else {
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + buildCmd)
			}

			if err = shellExec.ExecShell(buildCmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cmd *buildCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *buildCmd) initFlags() error {
	cmd.cobraCmd.Flags().StringVarP(&DOCKER_IMAGE_TAG, "DOCKER_IMAGE_TAG", "s", "latest", "docker image DOCKER_IMAGE_TAG")
	return nil
}
