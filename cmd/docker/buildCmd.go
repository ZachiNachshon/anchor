package docker

import (
	"github.com/anchor/pkg/utils/extractor"
	"github.com/anchor/pkg/utils/locator"
	"strings"

	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
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
		Short: "Builds a docker image",
		Long:  `Builds a docker image`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Build: Docker Image")

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

func buildDockerfile(identifier string) error {
	logger.Info("\n==> Building image...\n")
	if buildCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandBuild); err != nil {
		return err
	} else {
		contextPath, _ := locator.DirLocator.DockerContext(identifier)

		// Replace docker build "." with directory absolute path
		ctxIdx := strings.LastIndex(buildCmd, ".")
		buildCmd = buildCmd[:ctxIdx]
		buildCmd += contextPath

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + buildCmd + "\n")
		}

		if err = common.ShellExec.Execute(buildCmd); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *buildCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *buildCmd) initFlags() error {
	cmd.cobraCmd.Flags().StringVarP(
		&common.GlobalOptions.DockerImageTag,
		"tag",
		"t",
		common.GlobalOptions.DockerImageTag,
		"anchor docker build <name> -t my_tag")
	return nil
}
