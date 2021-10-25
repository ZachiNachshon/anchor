package use_context

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

func Test_UseContextActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
		{
			Name: "view configuration file",
			Func: ViewConfigurationFile,
		},
		{
			Name: "fail to get config context",
			Func: FailToGetConfigContext,
		},
		{
			Name: "fail to override config entry",
			Func: FailToOverrideConfigEntry,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			configContextName := "test-cfg-context"
			fakeCfgMgr := config.CreateFakeConfigManager()
			fakeO := NewOrchestrator(fakeCfgMgr, configContextName)
			runCallCount := 0
			fakeO.RunFunc = func(o *UseContextOrchestrator, ctx common.Context) error {
				runCallCount++
				return nil
			}
			err := ConfigUseContext(ctx, fakeO)
			assert.Nil(t, err, "expected not to fail")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var ViewConfigurationFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "1st-anchorfiles"
				fakeCfgManager := config.CreateFakeConfigManager()
				overrideCfgCallCount := 0
				fakeCfgManager.OverrideConfigEntryMock = func(entryName string, value interface{}) error {
					overrideCfgCallCount++
					return nil
				}

				fakeO := NewOrchestrator(fakeCfgManager, configContextName)
				err := run(fakeO, ctx)
				assert.Nil(t, err, "expected use context to succeed")
				assert.Equal(t, 1, overrideCfgCallCount, "expected to be called exactly once")
			})
		})
	})
}

var FailToGetConfigContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "invalid-config-context-name"
				fakeCfgManager := config.CreateFakeConfigManager()
				fakeO := NewOrchestrator(fakeCfgManager, configContextName)
				err := run(fakeO, ctx)
				assert.NotNil(t, err, "expected use context to fail")
				assert.Contains(t, err.Error(), "could not identify config context")
			})
		})
	})
}

var FailToOverrideConfigEntry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "1st-anchorfiles"
				fakeCfgManager := config.CreateFakeConfigManager()
				fakeCfgManager.OverrideConfigEntryMock = func(entryName string, value interface{}) error {
					return fmt.Errorf("failed to override config entry")
				}

				fakeO := NewOrchestrator(fakeCfgManager, configContextName)
				err := run(fakeO, ctx)
				assert.NotNil(t, err, "expected use context to fail")
				assert.Equal(t, "failed to override config entry", err.Error())
			})
		})
	})
}
