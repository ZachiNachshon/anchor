package set_context_entry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

const (
	remoteUrlFlagName                 = "repository.remote.url"
	remoteBranchFlagName              = "repository.remote.branch"
	remoteRevisionFlagName            = "repository.remote.revision"
	remoteClonePathFlagName           = "repository.remote.clonePath"
	remoteAutoUpdateFlagName          = "repository.remote.autoUpdate"
	localPathFlagName                 = "repository.local.path"
	setAsCurrentConfigContextFlagName = "set-current-context"
)

var remoteBranchFlagValue = ""
var remoteUrlFlagValue = ""
var remoteRevisionFlagValue = ""
var remoteClonePathFlagValue = ""
var remoteAutoUpdateFlagValue = ""
var localPathFlagValue = ""
var setAsCurrentConfigContextFlagValue = false

type setContextValueCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context

	initFlagsFunc func(o *setContextValueCmd) error
}

type NewCommandFunc func(
	ctx common.Context,
	cfgManager config.ConfigManager,
	setContextEntryFunc ConfigSetContextEntryFunc) *setContextValueCmd

func NewCommand(
	ctx common.Context,
	cfgManager config.ConfigManager,
	setContextEntryFunc ConfigSetContextEntryFunc) *setContextValueCmd {

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
			return setContextEntryFunc(ctx, NewOrchestrator(cfgManager, cfgCtxName, setAsCurrentConfigContextFlagValue, flags))
		},
	}

	return &setContextValueCmd{
		cobraCmd:      cobraCmd,
		ctx:           ctx,
		initFlagsFunc: initFlags,
	}
}

func (c *setContextValueCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *setContextValueCmd) GetContext() common.Context {
	return c.ctx
}

func initFlags(c *setContextValueCmd) error {
	c.cobraCmd.Flags().StringVar(
		&remoteRevisionFlagValue,
		remoteRevisionFlagName,
		"",
		fmt.Sprintf("--%s=3x4MPl3R3v1510N", remoteRevisionFlagName))

	c.cobraCmd.Flags().StringVar(
		&remoteUrlFlagValue,
		remoteUrlFlagName,
		"",
		fmt.Sprintf("--%s=git@some-repo", remoteUrlFlagName))

	c.cobraCmd.Flags().StringVar(
		&remoteBranchFlagValue,
		remoteBranchFlagName,
		"",
		fmt.Sprintf("--%s=main", remoteBranchFlagName))

	c.cobraCmd.Flags().StringVar(
		&remoteClonePathFlagValue,
		remoteClonePathFlagName,
		"",
		fmt.Sprintf("--%s=/repo/remote/clone/path", remoteClonePathFlagName))

	c.cobraCmd.Flags().StringVar(
		&remoteAutoUpdateFlagValue,
		remoteAutoUpdateFlagName,
		"",
		fmt.Sprintf("--%s=true", remoteAutoUpdateFlagName))

	c.cobraCmd.Flags().StringVar(
		&localPathFlagValue,
		localPathFlagName,
		"",
		fmt.Sprintf("--%s=/repo/local/path", localPathFlagName))

	c.cobraCmd.Flags().BoolVar(
		&setAsCurrentConfigContextFlagValue,
		setAsCurrentConfigContextFlagName,
		false,
		fmt.Sprint("--", setAsCurrentConfigContextFlagName))

	c.cobraCmd.Flags().SortFlags = false
	return nil
}

func AddCommand(
	parent cmd.AnchorCommand,
	cfgManager config.ConfigManager,
	createCmd NewCommandFunc) error {

	newCmd := createCmd(parent.GetContext(), cfgManager, ConfigSetContextEntry)
	err := newCmd.initFlagsFunc(newCmd)
	if err != nil {
		return err
	}
	parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
