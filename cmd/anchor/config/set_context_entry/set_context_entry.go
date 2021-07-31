package set_context_entry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config/resolver"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
	"github.com/spf13/cobra"
)

var remoteBranchFlagValue = ""
var remoteUrlFlagValue = ""
var remoteRevisionFlagValue = ""
var remoteClonePathFlagValue = ""
var remoteAutoUpdateFlagValue = ""
var localPathFlagValue = ""

type setContextValueCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(ctx common.Context, configActions *cfg.ConfigurationActions) *setContextValueCmd {
	var cobraCmd = &cobra.Command{
		Use: fmt.Sprintf(
			`set-context-entry [CURRENT_CONTEXT_NAME]
	[--%s=URL] 
	[--%s=BRANCH] 
	[--%s=REVISION] 
	[--%s=CLONE_PATH] 
	[--%s=AUTO_UPDATE]
`,
			resolver.RemoteUrlFlagName,
			resolver.RemoteBranchFlagName,
			resolver.RemoteRevisionFlagName,
			resolver.RemoteClonePathFlagName,
			resolver.RemoteAutoUpdateFlagName,
		),
		Short:                 "Update config context supported entries",
		Long:                  `Update config context supported entries`,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true, // Fatal errors are being logged by parent anchor.go
		SilenceErrors:         true, // Fatal errors are being logged by parent anchor.go
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgCtxName := args[0]
			flags := make(map[string]string)
			if len(remoteUrlFlagValue) > 0 {
				flags[resolver.RemoteUrlFlagName] = remoteUrlFlagValue
			}
			if len(remoteBranchFlagValue) > 0 {
				flags[resolver.RemoteBranchFlagName] = remoteBranchFlagValue
			}
			if len(remoteRevisionFlagValue) > 0 {
				flags[resolver.RemoteRevisionFlagName] = remoteRevisionFlagValue
			}
			if len(remoteClonePathFlagValue) > 0 {
				flags[resolver.RemoteClonePathFlagName] = remoteClonePathFlagValue
			}
			if len(remoteAutoUpdateFlagValue) > 0 {
				flags[resolver.RemoteAutoUpdateFlagName] = remoteAutoUpdateFlagValue
			}
			if len(localPathFlagValue) > 0 {
				flags[resolver.LocalPathFlagName] = localPathFlagValue
			}
			return configActions.SetContextEntry(ctx, cfgCtxName, flags)
		},
	}

	var cmd = &setContextValueCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	cmd.InitFlags()
	return cmd
}

func (cmd *setContextValueCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *setContextValueCmd) InitFlags() {
	cmd.cobraCmd.Flags().StringVar(
		&remoteRevisionFlagValue,
		resolver.RemoteRevisionFlagName,
		"",
		fmt.Sprintf("--%s=3x4MPl3R3v1510N", resolver.RemoteRevisionFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteUrlFlagValue,
		resolver.RemoteUrlFlagName,
		"",
		fmt.Sprintf("--%s=git@some-repo", resolver.RemoteUrlFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteBranchFlagValue,
		resolver.RemoteBranchFlagName,
		"",
		fmt.Sprintf("--%s=main", resolver.RemoteBranchFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteClonePathFlagValue,
		resolver.RemoteClonePathFlagName,
		"",
		fmt.Sprintf("--%s=/repo/remote/clone/path", resolver.RemoteClonePathFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteAutoUpdateFlagValue,
		resolver.RemoteAutoUpdateFlagName,
		"",
		fmt.Sprintf("--%s=true", resolver.RemoteAutoUpdateFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&localPathFlagValue,
		resolver.LocalPathFlagName,
		"",
		fmt.Sprintf("--%s=/repo/local/path", resolver.LocalPathFlagName))

	cmd.cobraCmd.Flags().SortFlags = false
}

func (cmd *setContextValueCmd) InitSubCommands() {
}
