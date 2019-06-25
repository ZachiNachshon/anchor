package docker

import (
	"fmt"

	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type cleanCmd struct {
	cobraCmd *cobra.Command
	opts     CleanCmdOptions
}

type CleanCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCleanCmd(opts *common.CmdRootOptions) *cleanCmd {
	var cobraCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean docker containers and images",
		Long:  `Clean docker containers and images`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Cleanup: Containers & Images")
			if err := stopContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := removeContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := removeImages(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var cleanCmd = new(cleanCmd)
	cleanCmd.cobraCmd = cobraCmd
	cleanCmd.opts.CmdRootOptions = opts

	if err := cleanCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return cleanCmd
}

func removeImages(dirname string) error {
	removeImagesFmt := "docker rmi -f %v"

	unknownImagesCmd := "docker images | grep '<none>' | awk {'print $3'}"
	if unknownImages, err := common.ShellExec.ExecuteWithOutput(unknownImagesCmd); err != nil {
		return err
	} else if len(unknownImages) > 0 {
		logger.Info("Removing docker images for name: <none>")
		imageIds := fmt.Sprintf(removeImagesFmt, unknownImages)
		_ = common.ShellExec.Execute(imageIds)
	} else {
		logger.Info("No images can be found for name: <none>")
	}

	imageIdentifier := composeDockerImageIdentifierNoTag(dirname)
	containerImagesCmd := fmt.Sprintf("docker images | grep '%v' | awk {'print $3'}", imageIdentifier)
	if containerImages, err := common.ShellExec.ExecuteWithOutput(containerImagesCmd); err != nil {
		return err
	} else if len(containerImages) > 0 {
		logger.Infof("Removing docker images for name: %v", imageIdentifier)
		imageIds := fmt.Sprintf(removeImagesFmt, containerImages)
		_ = common.ShellExec.Execute(imageIds)
	} else {
		logger.Infof("No images can be found for name: %v", imageIdentifier)
	}

	return nil
}

func (cmd *cleanCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *cleanCmd) initFlags() error {
	return nil
}
