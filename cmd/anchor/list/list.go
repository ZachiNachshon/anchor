package list

import (
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/spf13/cobra"
)

var ListOnlyK8sManifestsFlag = false
var AffinityFilterFlag = ""

type listCmd struct {
	cobraCmd *cobra.Command
	opts     ListOptions
}

type ListOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *listCmd {
	var cobraCmd = &cobra.Command{
		Use:   "list",
		Short: "List all supported directories under DOCKER_FILES folder",
		Long:  `List all supported directories under DOCKER_FILES folder`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.DockerHeadline, "List")
			printSupportedDirs()
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
	cmd.cobraCmd.PersistentFlags().BoolVarP(
		&ListOnlyK8sManifestsFlag,
		"filter kubernetes manifests only",
		"k",
		ListOnlyK8sManifestsFlag,
		"anchor list -k")

	cmd.cobraCmd.PersistentFlags().StringVarP(
		&AffinityFilterFlag,
		"filter by affinity",
		"a",
		AffinityFilterFlag,
		"anchor list -a affinity-name")

	return nil
}

func printSupportedDirs() {
	opts := &locator.ListOpts{
		OnlyK8sManifests: ListOnlyK8sManifestsFlag,
		AffinityFilter:   AffinityFilterFlag,
	}
	locator.DirLocator.Print(opts)
}
