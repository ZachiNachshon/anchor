package build

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/docker"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type buildCmd struct {
	cobraCmd *cobra.Command
	opts     BuildOptions
}

type BuildOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *buildCmd {
	var cobraCmd = &cobra.Command{
		Use:   "build",
		Short: "Builds a docker image",
		Long:  `Builds a docker image`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := docker.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.DockerHeadline, "Build")

			if err := docker.Build(args[0]); err != nil {
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
