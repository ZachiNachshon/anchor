package cmd

import (
	"fmt"

	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	cobraCmd *cobra.Command
	opts     VersionCmdOptions
}

type VersionCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewVersionCmd(opts *common.CmdRootOptions) *versionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "version",
		Short: "Print anchor version",
		Long:  `Print anchor version`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: take from config
			fmt.Println("v1.0.0")
		},
	}

	var versionCmd = new(versionCmd)
	versionCmd.cobraCmd = cobraCmd
	versionCmd.opts.CmdRootOptions = opts

	if err := versionCmd.initFlags(); err != nil {
		// TODO: log error
	}

	return versionCmd
}

func (cmd *versionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *versionCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}
