package docker

import (
	"fmt"

	"github.com/anchor/pkg/common"
	"github.com/spf13/cobra"
)

type pushCmd struct {
	cobraCmd *cobra.Command
	opts     PushCmdOptions
}

type PushCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewPushCmd(opts *common.CmdRootOptions) *pushCmd {
	var cobraCmd = &cobra.Command{
		Use:   "push",
		Short: "Push a docker image to remote/local repository",
		Long:  `Push a docker image to remote/local repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Clean docker containers and images...")
		},
	}

	var pushCmd = new(pushCmd)
	pushCmd.cobraCmd = cobraCmd
	pushCmd.opts.CmdRootOptions = opts

	if err := pushCmd.initFlags(); err != nil {
		// TODO: log error
	}

	return pushCmd
}

func (cmd *pushCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *pushCmd) initFlags() error {
	return nil
}
