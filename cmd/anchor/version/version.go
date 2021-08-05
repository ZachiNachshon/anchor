package version

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/version"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, verActions *version.VersionActions) *versionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "version",
		Short: "Print anchor CLI version",
		Long:  `Print anchor CLI version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return verActions.Version(ctx)
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

func (cmd *versionCmd) InitFlags() {
}

func (cmd *versionCmd) InitSubCommands() {
}
