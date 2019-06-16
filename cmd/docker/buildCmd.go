package docker

import (
	"fmt"
	"github.com/kit/cmd/logger"
	"github.com/kit/cmd/types"
	"github.com/kit/cmd/utils"
	"github.com/spf13/cobra"
	"path/filepath"
)

var suite = "latest"

type BuildCmd struct {
	cobraCmd *cobra.Command
	opts     BuildCmdOptions
}

type BuildCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewBuildCmd(opts *common.CmdRootOptions) *BuildCmd {
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

	var buildCmd = new(BuildCmd)
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
		buildDir := filepath.Dir(dockerfilePath)
		imageIdentifier := composeDockerImageIdentifier(dirname)
		logger.Infof("Building %v...", imageIdentifier)
		buildCmd := fmt.Sprintf("docker build -f %v -t %v %v", dockerfilePath, imageIdentifier, buildDir)
		utils.ExecShell(buildCmd)
	}

	return nil
}

//func extractDockerBuildCmd(path string) {
//	if content, err := os.Open(path); err != nil {
//
//	}
//}

func (cmd *BuildCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *BuildCmd) initFlags() error {
	cmd.cobraCmd.Flags().StringVarP(&suite, "suite", "s", "latest", "docker image suite")
	return nil
}
