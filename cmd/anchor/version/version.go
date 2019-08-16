package version

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	cobraCmd *cobra.Command
	opts     VersionOptions
}

type VersionOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

// TODO: take from config
const Version = "v0.3.0"

func NewCommand(opts *common.CmdRootOptions) *versionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "version",
		Short: "Print anchor CLI version",
		Long:  `Print anchor CLI version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
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
	return nil
}
