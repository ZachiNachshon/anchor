package docker

import (
	"fmt"
	"github.com/anchor/pkg/utils/locator"
	"strings"

	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type removeCmd struct {
	cobraCmd *cobra.Command
	opts     RemoveCmdOptions
}

type RemoveCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewRemoveCmd(opts *common.CmdRootOptions) *removeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove docker containers and images",
		Long:  `Remove docker containers and images`,
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

	var removeCmd = new(removeCmd)
	removeCmd.cobraCmd = cobraCmd
	removeCmd.opts.CmdRootOptions = opts

	if err := removeCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return removeCmd
}

func removeImages(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Removing image [%v]", name))
	}

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

func (cmd *removeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *removeCmd) initFlags() error {
	return nil
}
