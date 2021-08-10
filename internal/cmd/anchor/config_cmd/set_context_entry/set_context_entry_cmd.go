package set_context_entry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

const (
	remoteUrlFlagName        = "repository.remote.url"
	remoteBranchFlagName     = "repository.remote.branch"
	remoteRevisionFlagName   = "repository.remote.revision"
	remoteClonePathFlagName  = "repository.remote.clonePath"
	remoteAutoUpdateFlagName = "repository.remote.autoUpdate"
	localPathFlagName        = "repository.local.path"
)

var remoteBranchFlagValue = ""
var remoteUrlFlagValue = ""
var remoteRevisionFlagValue = ""
var remoteClonePathFlagValue = ""
var remoteAutoUpdateFlagValue = ""
var localPathFlagValue = ""

type setContextValueCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

func NewCommand(
	ctx common.Context,
	cfgManager config.ConfigManager,
	setContextEntryFunc ConfigSetContextEntryFunc) (*setContextValueCmd, error) {

	var cobraCmd = &cobra.Command{
		Use: fmt.Sprintf(
			`set-context-entry [CURRENT_CONTEXT_NAME]
	[--%s=URL] 
	[--%s=BRANCH] 
	[--%s=REVISION] 
	[--%s=CLONE_PATH] 
	[--%s=AUTO_UPDATE]
`,
			remoteUrlFlagName,
			remoteBranchFlagName,
			remoteRevisionFlagName,
			remoteClonePathFlagName,
			remoteAutoUpdateFlagName,
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
				flags[remoteUrlFlagName] = remoteUrlFlagValue
			}
			if len(remoteBranchFlagValue) > 0 {
				flags[remoteBranchFlagName] = remoteBranchFlagValue
			}
			if len(remoteRevisionFlagValue) > 0 {
				flags[remoteRevisionFlagName] = remoteRevisionFlagValue
			}
			if len(remoteClonePathFlagValue) > 0 {
				flags[remoteClonePathFlagName] = remoteClonePathFlagValue
			}
			if len(remoteAutoUpdateFlagValue) > 0 {
				flags[remoteAutoUpdateFlagName] = remoteAutoUpdateFlagValue
			}
			if len(localPathFlagValue) > 0 {
				flags[localPathFlagName] = localPathFlagValue
			}
			return setContextEntryFunc(ctx, cfgCtxName, flags, cfgManager)
		},
	}

	var cmd = &setContextValueCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	err := cmd.InitFlags()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (cmd *setContextValueCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *setContextValueCmd) InitFlags() error {
	cmd.cobraCmd.Flags().StringVar(
		&remoteRevisionFlagValue,
		remoteRevisionFlagName,
		"",
		fmt.Sprintf("--%s=3x4MPl3R3v1510N", remoteRevisionFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteUrlFlagValue,
		remoteUrlFlagName,
		"",
		fmt.Sprintf("--%s=git@some-repo", remoteUrlFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteBranchFlagValue,
		remoteBranchFlagName,
		"",
		fmt.Sprintf("--%s=main", remoteBranchFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteClonePathFlagValue,
		remoteClonePathFlagName,
		"",
		fmt.Sprintf("--%s=/repo/remote/clone/path", remoteClonePathFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&remoteAutoUpdateFlagValue,
		remoteAutoUpdateFlagName,
		"",
		fmt.Sprintf("--%s=true", remoteAutoUpdateFlagName))

	cmd.cobraCmd.Flags().StringVar(
		&localPathFlagValue,
		localPathFlagName,
		"",
		fmt.Sprintf("--%s=/repo/local/path", localPathFlagName))

	cmd.cobraCmd.Flags().SortFlags = false
	return nil
}

func (cmd *setContextValueCmd) InitSubCommands() error {
	return nil
}
