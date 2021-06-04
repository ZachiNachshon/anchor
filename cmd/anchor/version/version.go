package version

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context
}

// TODO: take from config
const Version = "v0.4.0"

func NewCommand(ctx common.Context) *versionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "version",
		Short: "Print anchor CLI version",
		Long:  `Print anchor CLI version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}

	return &versionCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}
}

func (cmd *versionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}
