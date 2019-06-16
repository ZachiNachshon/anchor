package docker

import (
	"fmt"
	"github.com/kit/cmd/types"

	"github.com/spf13/cobra"
)

type PushCmd struct {
	cobraCmd *cobra.Command
	opts     PushCmdOptions
}

type PushCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewPushCmd(opts *common.CmdRootOptions) *PushCmd {
	var cobraCmd = &cobra.Command{
		Use:   "push",
		Short: "Push a docker image to remote/local repository",
		Long:  `Push a docker image to remote/local repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Clean docker containers and images...")
		},
	}

	var pushCmd = new(PushCmd)
	pushCmd.cobraCmd = cobraCmd
	pushCmd.opts.CmdRootOptions = opts

	if err := pushCmd.initFlags(); err != nil {
		// TODO: log error
	}

	return pushCmd
}

func (cmd *PushCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *PushCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}
