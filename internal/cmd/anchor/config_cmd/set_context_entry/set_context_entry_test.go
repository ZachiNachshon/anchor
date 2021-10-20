package set_context_entry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SetContextEntryActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
		{
			Name: "set config context entries",
			Func: SetConfigContextEntries,
		},
		{
			Name: "fail to override config entries",
			Func: FailToOverrideConfigEntries,
		},
		{
			Name: "fail to set new config context defaults",
			Func: FailToSetNewConfigContextDefaults,
		},
		{
			Name: "add new config context with defaults",
			Func: AddNewConfigContextWithDefaults,
		},
		{
			Name: "populate config context changes",
			Func: PopulateConfigContextChanges,
		},
		{
			Name: "fail to populate config context changes on bad input type",
			Func: FailToPopulateConfigContextChangesOnBadInputType,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			configContextName := "test-cfg-context"
			fakeCfgMgr := config.CreateFakeConfigManager()
			fakeO := NewOrchestrator(fakeCfgMgr, configContextName, nil)
			runCallCount := 0
			fakeO.runFunc = func(o *setContextEntryOrchestrator, ctx common.Context) error {
				runCallCount++
				return nil
			}
			err := ConfigSetContextEntry(ctx, fakeO)
			assert.Nil(t, err, "expected not to fail")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var SetConfigContextEntries = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				changes := make(map[string]string)
				changes[remoteUrlFlagName] = remoteUrlFlagValue
				changes[remoteBranchFlagName] = remoteBranchFlagValue
				changes[remoteRevisionFlagName] = remoteRevisionFlagValue

				configContextName := "1st-anchorfiles"
				fakeCfgManager := config.CreateFakeConfigManager()
				overrideCfgCallCount := 0
				fakeCfgManager.OverrideConfigMock = func(cfgToUpdate *config.AnchorConfig) error {
					cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, configContextName)
					assert.Equal(t, changes[remoteUrlFlagName], cfgCtx.Context.Repository.Remote.Url)
					assert.Equal(t, changes[remoteBranchFlagName], cfgCtx.Context.Repository.Remote.Branch)
					assert.Equal(t, changes[remoteRevisionFlagName], cfgCtx.Context.Repository.Remote.Revision)
					overrideCfgCallCount++
					return nil
				}

				fakeO := NewOrchestrator(fakeCfgManager, configContextName, changes)
				err := run(fakeO, ctx)
				assert.Nil(t, err, "expected set context entry to succeed")
				assert.Equal(t, 1, overrideCfgCallCount, "expected to be called exactly once")
			})
		})
	})
}

var FailToOverrideConfigEntries = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				changes := make(map[string]string)
				configContextName := "1st-anchorfiles"
				fakeCfgManager := config.CreateFakeConfigManager()
				overrideCfgCallCount := 0
				fakeCfgManager.OverrideConfigMock = func(cfgToUpdate *config.AnchorConfig) error {
					overrideCfgCallCount++
					return fmt.Errorf("failed to override config entries")
				}

				fakeO := NewOrchestrator(fakeCfgManager, configContextName, changes)
				err := run(fakeO, ctx)
				assert.NotNil(t, err, "expected set context entry to fail")
				assert.Equal(t, "failed to override config entries", err.Error())
				assert.Equal(t, 1, overrideCfgCallCount, "expected to be called exactly once")
			})
		})
	})
}

var FailToSetNewConfigContextDefaults = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "new-config-context-name"
				fakeCfgManager := config.CreateFakeConfigManager()
				fakeCfgManager.SetDefaultsPostCreationMock = func(anchorConfig *config.AnchorConfig) error {
					return fmt.Errorf("failed to set defaults post config context creation")
				}
				fakeO := NewOrchestrator(fakeCfgManager, configContextName, nil)
				fakeO.cfgManager = fakeCfgManager
				err := run(fakeO, ctx)
				assert.NotNil(t, err, "expected set context entry to fail")
				assert.Equal(t, err.Error(), "failed to set defaults post config context creation")
			})
		})
	})
}

var PopulateConfigContextChanges = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "1st-anchorfiles"
				cfgCtx := config.TryGetConfigContext(cfg.Config.Contexts, configContextName)

				url := "git@github.com:test/flags.git"
				branch := "test-branch"
				revision := "test-revision"
				clonePath := "/test/clone/path"
				autoUpdate := "true"
				localPath := "/test/local/path"

				changes := make(map[string]string)
				changes[remoteUrlFlagName] = url
				changes[remoteBranchFlagName] = branch
				changes[remoteRevisionFlagName] = revision
				changes[remoteClonePathFlagName] = clonePath
				changes[remoteAutoUpdateFlagName] = autoUpdate
				changes[localPathFlagName] = localPath
				err := populateConfigContextChanges(cfgCtx, changes)
				assert.Nil(t, err, "expected use context to fail")
				assert.Equal(t, cfgCtx.Context.Repository.Remote.Url, url)
				assert.Equal(t, cfgCtx.Context.Repository.Remote.Branch, branch)
				assert.Equal(t, cfgCtx.Context.Repository.Remote.Revision, revision)
				assert.Equal(t, cfgCtx.Context.Repository.Remote.ClonePath, clonePath)
				assert.Equal(t, cfgCtx.Context.Repository.Remote.AutoUpdate, true)
				assert.Equal(t, cfgCtx.Context.Repository.Local.Path, localPath)
			})
		})
	})
}

var FailToPopulateConfigContextChangesOnBadInputType = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				changes := make(map[string]string)
				changes[remoteAutoUpdateFlagName] = "bad-boolean"
				configContextName := "1st-anchorfiles"
				fakeCfgManager := config.CreateFakeConfigManager()
				overrideCfgCallCount := 0
				fakeCfgManager.OverrideConfigMock = func(cfgToUpdate *config.AnchorConfig) error {
					overrideCfgCallCount++
					return nil
				}
				fakeO := NewOrchestrator(fakeCfgManager, configContextName, changes)
				err := run(fakeO, ctx)
				assert.NotNil(t, err, "expected set context entry to fail")
				assert.Equal(t, 0, overrideCfgCallCount, "expected not to be called")
			})
		})
	})
}

var AddNewConfigContextWithDefaults = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				changes := make(map[string]string)
				changes[remoteUrlFlagName] = "git@github.com:test/flags.git"
				configContextName := "new-cfg-context"
				fakeCfgManager := config.CreateFakeConfigManager()
				setDefaultsPostCreationCallCount := 0
				fakeCfgManager.SetDefaultsPostCreationMock = func(anchorConfig *config.AnchorConfig) error {
					setDefaultsPostCreationCallCount++
					newCfgCtxAdded := config.TryGetConfigContext(anchorConfig.Config.Contexts, configContextName)
					assert.NotNil(t, newCfgCtxAdded, "expected contexts to contains the new cfg context")
					newCfgCtxAdded.Context.Repository.Remote.Branch = config.DefaultRemoteBranch
					return nil
				}
				overrideCfgCallCount := 0
				fakeCfgManager.OverrideConfigMock = func(cfgToUpdate *config.AnchorConfig) error {
					overrideCfgCallCount++
					newCfgCtxAdded := config.TryGetConfigContext(cfgToUpdate.Config.Contexts, configContextName)
					assert.Equal(t, config.DefaultRemoteBranch, newCfgCtxAdded.Context.Repository.Remote.Branch)
					return nil
				}
				fakeO := NewOrchestrator(fakeCfgManager, configContextName, changes)
				err := run(fakeO, ctx)
				assert.Nil(t, err, "expected set context entry to succeed")
				assert.Equal(t, 1, setDefaultsPostCreationCallCount, "expected not to be called")
				assert.Equal(t, 1, overrideCfgCallCount, "expected not to be called")
			})
		})
	})
}
