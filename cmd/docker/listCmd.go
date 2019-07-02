package docker

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/locator"
	"github.com/spf13/cobra"
)

type listCmd struct {
	cobraCmd *cobra.Command
	opts     ListCmdOptions
}

type ListCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewListCmd(opts *common.CmdRootOptions) *listCmd {
	var cobraCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available docker supported images from DOCKER_FILES repository",
		Long:  `List all available docker supported images from DOCKER_FILES repository`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Listing all Docker images")

			if _, err := locator.GetDirNamesNoPath(true, locator.DOCKER_FILE_IDENTIFIER); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var listCmd = new(listCmd)
	listCmd.cobraCmd = cobraCmd
	listCmd.opts.CmdRootOptions = opts

	if err := listCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return listCmd
}

func (cmd *listCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *listCmd) initFlags() error {
	return nil
}
