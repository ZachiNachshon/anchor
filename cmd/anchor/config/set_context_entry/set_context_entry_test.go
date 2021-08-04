package set_context_entry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/resolver"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &cfg.ConfigurationActions{
					SetContextEntry: func(ctx common.Context, cfgCtxName string, changes map[string]string) error {
						callCount++
						return nil
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions), "test-cfg-context")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: set-context-entry")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var StartSetContextEntryActionWithAllFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				configContextName := "test-cfg-context"
				url := "git@github.com:test/flags.git"
				branch := "test-branch"
				revision := "test-revision"
				clonePath := "/test/clone/path"
				autoUpdate := "true"
				localPath := "/test/local/path"

				callCount := 0
				actions := &cfg.ConfigurationActions{
					SetContextEntry: func(ctx common.Context, cfgCtxName string, changes map[string]string) error {
						assert.Equal(t, cfgCtxName, configContextName)

						assert.Contains(t, changes, resolver.RemoteUrlFlagName)
						assert.Equal(t, changes[resolver.RemoteUrlFlagName], url)

						assert.Contains(t, changes, resolver.RemoteBranchFlagName)
						assert.Equal(t, changes[resolver.RemoteBranchFlagName], branch)

						assert.Contains(t, changes, resolver.RemoteRevisionFlagName)
						assert.Equal(t, changes[resolver.RemoteRevisionFlagName], revision)

						assert.Contains(t, changes, resolver.RemoteClonePathFlagName)
						assert.Equal(t, changes[resolver.RemoteClonePathFlagName], clonePath)

						assert.Contains(t, changes, resolver.RemoteAutoUpdateFlagName)
						assert.Equal(t, changes[resolver.RemoteAutoUpdateFlagName], autoUpdate)

						assert.Contains(t, changes, resolver.LocalPathFlagName)
						assert.Equal(t, changes[resolver.LocalPathFlagName], localPath)

						callCount++
						return nil
					},
				}

				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions),
					configContextName,
					fmt.Sprintf("--%s=%s", resolver.RemoteUrlFlagName, url),
					fmt.Sprintf("--%s=%s", resolver.RemoteBranchFlagName, branch),
					fmt.Sprintf("--%s=%s", resolver.RemoteRevisionFlagName, revision),
					fmt.Sprintf("--%s=%s", resolver.RemoteClonePathFlagName, clonePath),
					fmt.Sprintf("--%s=%s", resolver.RemoteAutoUpdateFlagName, autoUpdate),
					fmt.Sprintf("--%s=%s", resolver.LocalPathFlagName, localPath),
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &cfg.ConfigurationActions{
					SetContextEntry: func(ctx common.Context, cfgCtxName string, changes map[string]string) error {
						callCount++
						return nil
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions))
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				configContextName := "test-cfg-context"
				callCount := 0
				actions := &cfg.ConfigurationActions{
					SetContextEntry: func(ctx common.Context, cfgCtxName string, changes map[string]string) error {
						callCount++
						return fmt.Errorf("an error occurred")
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions), configContextName)
				assert.Equal(t, 1, callCount, "expected action not to be called. name: set-context-entry")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, "an error occurred", err.Error())
			})
		})
	})
}
