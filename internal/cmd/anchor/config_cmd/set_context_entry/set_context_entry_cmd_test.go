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
	}
	harness.RunTests(t, tests)
}

var StartSetContextEntryActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, cfgCtxName string, changes map[string]string, cfgManager config.ConfigManager) error {
					callCount++
					return nil
				}
				command, err := NewCommand(ctx, fakeCfgManager, fun)
				_, err = drivers.CLI().RunCommand(command, "test-cfg-context")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: set-context-entry")
				assert.Nil(t, err, "expected cli action to have no errors")
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
				fun := func(ctx common.Context, cfgCtxName string, changes map[string]string, cfgManager config.ConfigManager) error {
					assert.Equal(t, cfgCtxName, configContextName)

					assert.Contains(t, changes, remoteUrlFlagName)
					assert.Equal(t, changes[remoteUrlFlagName], url)

					assert.Contains(t, changes, remoteBranchFlagName)
					assert.Equal(t, changes[remoteBranchFlagName], branch)

					assert.Contains(t, changes, remoteRevisionFlagName)
					assert.Equal(t, changes[remoteRevisionFlagName], revision)

					assert.Contains(t, changes, remoteClonePathFlagName)
					assert.Equal(t, changes[remoteClonePathFlagName], clonePath)

					assert.Contains(t, changes, remoteAutoUpdateFlagName)
					assert.Equal(t, changes[remoteAutoUpdateFlagName], autoUpdate)

					assert.Contains(t, changes, localPathFlagName)
					assert.Equal(t, changes[localPathFlagName], localPath)

					callCount++
					return nil
				}

				command, err := NewCommand(ctx, fakeCfgManager, fun)
				_, err = drivers.CLI().RunCommand(command,
					configContextName,
					fmt.Sprintf("--%s=%s", remoteUrlFlagName, url),
					fmt.Sprintf("--%s=%s", remoteBranchFlagName, branch),
					fmt.Sprintf("--%s=%s", remoteRevisionFlagName, revision),
					fmt.Sprintf("--%s=%s", remoteClonePathFlagName, clonePath),
					fmt.Sprintf("--%s=%s", remoteAutoUpdateFlagName, autoUpdate),
					fmt.Sprintf("--%s=%s", localPathFlagName, localPath),
				)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: set-context-entry")
				assert.Nil(t, err, "expected cli action to have no errors")
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
				fun := func(ctx common.Context, cfgCtxName string, changes map[string]string, cfgManager config.ConfigManager) error {
					callCount++
					return nil
				}
				command, err := NewCommand(ctx, fakeCfgManager, fun)
				_, err = drivers.CLI().RunCommand(command)
				assert.Equal(t, 0, callCount, "expected action not to be called. name: set-context-entry")
				assert.NotNil(t, err, "expected cli action to fail")
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
				fun := func(ctx common.Context, cfgCtxName string, changes map[string]string, cfgManager config.ConfigManager) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command, err := NewCommand(ctx, fakeCfgManager, fun)
				_, err = drivers.CLI().RunCommand(command, configContextName)
				assert.Equal(t, 1, callCount, "expected action not to be called. name: set-context-entry")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, "an error occurred", err.Error())
			})
		})
	})
}
