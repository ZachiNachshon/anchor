package docker

import (
	"fmt"
	"github.com/anchor/pkg/utils/locator"
	"strings"

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

func removeImages(identifier string) error {
	var dirname = ""
	var err error
	if dirname, err = locator.DirLocator.Name(identifier); err != nil {
		return err
	}

	removeImagesFmt := "docker rmi -f %v"

	unknownImagesCmd := "docker images | grep '<none>' | awk {'print $3'}"
	if unknownImages, err := common.ShellExec.ExecuteWithOutput(unknownImagesCmd); err != nil {
		return err
	} else if len(unknownImages) > 0 {
		logger.Info("Removing docker images for name: <none>")
		removeUnknownCmd := fmt.Sprintf(removeImagesFmt, unknownImages)
		removeUnknownCmd = strings.Replace(removeUnknownCmd, "\n", " ", -1)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + removeUnknownCmd + "\n")
		}
		_ = common.ShellExec.Execute(removeUnknownCmd)

	} else {
		logger.Info("No images can be found for name: <none>")
	}

	imageIdentifier := ComposeDockerImageIdentifierNoTag(dirname)
	containerImagesCmd := fmt.Sprintf("docker images | grep '%v' | awk {'print $3'}", imageIdentifier)
	if containerImages, err := common.ShellExec.ExecuteWithOutput(containerImagesCmd); err != nil {
		return err
	} else if len(containerImages) > 0 {
		logger.Infof("Removing docker images for name: %v", imageIdentifier)
		removeImageCmd := fmt.Sprintf(removeImagesFmt, containerImages)
		removeImageCmd = strings.Replace(removeImageCmd, "\n", " ", -1)

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + removeImageCmd + "\n")
		}
		_ = common.ShellExec.Execute(removeImageCmd)

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
