package set_context_entry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SetContextEntryCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start set_context_entry action successfully",
			Func: StartSetContextEntryActionSuccessfully,
		},
		{
			Name: "start set_context_entry action with all flags",
			Func: StartSetContextEntryActionWithAllFlags,
		},
		{
			Name: "fail due to missing config context name",
			Func: FailDueToMissingConfigContextName,
		},
		{
			Name: "fail set context entry action",
			Func: FailSetContextEntryAction,
		},
		{
			Name: "contain cobra command",
			Func: ContainCobraCommand,
		},
		{
			Name: "contain context",
			Func: ContainContext,
		},
		{
			Name: "add itself to parent command",
			Func: AddItselfToParentCommand,
		},
		{
			Name: "fail to initialize flags",
			Func: FailToInitializeFlags,
		},
	}
	harness.RunTests(t, tests)
}

var StartSetContextEntryActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, o *setContextEntryOrchestrator) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				_, err := drivers.CLI().RunCommand(command, "test-cfg-context")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: set-context-entry")
				assert.Nil(t, err, "expected action to have no errors")
			})
		})
	})
}

var StartSetContextEntryActionWithAllFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "test-cfg-context"
				url := "git@github.com:test/flags.git"
				branch := "test-branch"
				revision := "test-revision"
				clonePath := "/test/clone/path"
				autoUpdate := "true"
				localPath := "/test/local/path"

				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, o *setContextEntryOrchestrator) error {
					assert.Equal(t, o.cfgCtxName, configContextName)

					assert.Contains(t, o.changes, remoteUrlFlagName)
					assert.Equal(t, o.changes[remoteUrlFlagName], url)

					assert.Contains(t, o.changes, remoteBranchFlagName)
					assert.Equal(t, o.changes[remoteBranchFlagName], branch)

					assert.Contains(t, o.changes, remoteRevisionFlagName)
					assert.Equal(t, o.changes[remoteRevisionFlagName], revision)

					assert.Contains(t, o.changes, remoteClonePathFlagName)
					assert.Equal(t, o.changes[remoteClonePathFlagName], clonePath)

					assert.Contains(t, o.changes, remoteAutoUpdateFlagName)
					assert.Equal(t, o.changes[remoteAutoUpdateFlagName], autoUpdate)

					assert.Contains(t, o.changes, localPathFlagName)
					assert.Equal(t, o.changes[localPathFlagName], localPath)

					assert.True(t, o.setAsCurrCfgCtx)

					callCount++
					return nil
				}

				command := NewCommand(ctx, fakeCfgManager, fun)
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command,
					configContextName,
					fmt.Sprintf("--%s=%s", remoteUrlFlagName, url),
					fmt.Sprintf("--%s=%s", remoteBranchFlagName, branch),
					fmt.Sprintf("--%s=%s", remoteRevisionFlagName, revision),
					fmt.Sprintf("--%s=%s", remoteClonePathFlagName, clonePath),
					fmt.Sprintf("--%s=%s", remoteAutoUpdateFlagName, autoUpdate),
					fmt.Sprintf("--%s=%s", localPathFlagName, localPath),
					fmt.Sprintf("--%s", setAsCurrentConfigContextFlagName),
				)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: set-context-entry")
				assert.Nil(t, err, "expected action to have no errors")
				assert.Equal(t, 7, command.cobraCmd.Flags().NFlag())
			})
		})
	})
}

var FailDueToMissingConfigContextName = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, o *setContextEntryOrchestrator) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 0, callCount, "expected action not to be called. name: set-context-entry")
				assert.NotNil(t, err, "expected action to fail")
				assert.Contains(t, "accepts 1 arg(s), received 0", err.Error())
			})
		})
	})
}

var FailSetContextEntryAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "test-cfg-context"
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, o *setContextEntryOrchestrator) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				_, err := drivers.CLI().RunCommand(command, configContextName)
				assert.Equal(t, 1, callCount, "expected action not to be called. name: set-context-entry")
				assert.NotNil(t, err, "expected action to fail")
				assert.Contains(t, "an error occurred", err.Error())
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		var fun ConfigSetContextEntryFunc = func(ctx common.Context, o *setContextEntryOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, fakeCfgManager, fun)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		var fun ConfigSetContextEntryFunc = func(ctx common.Context, o *setContextEntryOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, fakeCfgManager, fun)
		cmdCtx := anchorCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}

var AddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		parentCmd := NewCommand(ctx, fakeCfgManager, nil)
		err := AddCommand(parentCmd, fakeCfgManager, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		cmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(cmds))
		assert.Contains(t, cmds[0].Use, "set-context-entry")
	})
}

var FailToInitializeFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		parentCmd := NewCommand(ctx, fakeCfgManager, nil)

		var newCmdFunc NewCommandFunc = func(ctx common.Context, cfgManager config.ConfigManager, setContextEntryFunc ConfigSetContextEntryFunc) *setContextValueCmd {
			c := NewCommand(ctx, config.CreateFakeConfigManager(), nil)
			c.initFlagsFunc = func(o *setContextValueCmd) error {
				return fmt.Errorf("failed to initialize flags")
			}
			return c
		}

		err := AddCommand(parentCmd, fakeCfgManager, newCmdFunc)
		assert.NotNil(t, err, "expected add command to fail")
		assert.Equal(t, "failed to initialize flags", err.Error())
	})
}
