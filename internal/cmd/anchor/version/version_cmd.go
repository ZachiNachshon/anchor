package version

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, versionFunc VersionVersionFunc) (*versionCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:   "version",
		Short: "Print anchor CLI version",
		Long:  `Print anchor CLI version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return versionFunc(ctx)
		},
	}

	return &versionCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}, nil
}

func (cmd *versionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *versionCmd) InitFlags() error {
	return nil
}

func (cmd *versionCmd) InitSubCommands() error {
	return nil
}
