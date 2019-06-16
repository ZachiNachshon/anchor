package cmd

import (
	"fmt"
	"github.com/kit/cmd/types"
	"github.com/spf13/cobra"
)

type VersionCmd struct {
	cobraCmd *cobra.Command
	opts     VersionCmdOptions
}

type VersionCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewVersionCmd(opts *common.CmdRootOptions) *VersionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  `Print the version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print kit version...")
		},
	}

	var versionCmd = new(VersionCmd)
	versionCmd.cobraCmd = cobraCmd
	versionCmd.opts.CmdRootOptions = opts

	if err := versionCmd.initFlags(); err != nil {
		// TODO: log error
	}

	return versionCmd
}

func (cmd *VersionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *VersionCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}
