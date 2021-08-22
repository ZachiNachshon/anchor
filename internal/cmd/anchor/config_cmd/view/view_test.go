package view

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ViewActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
		{
			Name: "fail runner due to preparation",
			Func: FailRunnerDueToPreparation,
		},
		{
			Name: "prepare registry components",
			Func: PrepareRegistryComponents,
		},
		{
			Name: "fail resolving registry components",
			Func: FailResolvingRegistryComponents,
		},
		{
			Name: "view configuration file",
			Func: ViewConfigurationFile,
		},
		{
			Name: "fail to unmarshal config yaml",
			Func: FailToUnmarshalConfigYaml,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeCfgMgr := config.CreateFakeConfigManager()
			fakeO := NewOrchestrator(fakeCfgMgr)
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *viewOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return nil
			}
			runCallCount := 0
			fakeO.runFunc = func(o *viewOrchestrator, ctx common.Context) error {
				runCallCount++
				return nil
			}
			err := ConfigView(ctx, fakeO)
			assert.Nil(t, err, "expected not to fail")
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var FailRunnerDueToPreparation = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeCfgMgr := config.CreateFakeConfigManager()
			fakeO := NewOrchestrator(fakeCfgMgr)
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *viewOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return fmt.Errorf("failed to prepare runner")
			}
			err := ConfigView(ctx, fakeO)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare runner", err.Error())
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
		})
	})
}

var PrepareRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		fakeCfgMgr := config.CreateFakeConfigManager()
		fakeO := NewOrchestrator(fakeCfgMgr)
		err := fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
		assert.NotNil(t, fakeO.prntr)
	})
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		fakeCfgMgr := config.CreateFakeConfigManager()
		fakeO := NewOrchestrator(fakeCfgMgr)

		err := fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", printer.Identifier))

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
	})
}

var ViewConfigurationFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				fakeCfgManager := config.CreateFakeConfigManager()
				fakeCfgManager.GetConfigFilePathMock = func() (string, error) {
					return "/path/to/config", nil
				}

				fakePrinter := printer.CreateFakePrinter()
				printCfgCallCount := 0
				fakePrinter.PrintConfigurationMock = func(cfgFilePath string, cfgText string) {
					printCfgCallCount++
				}

				fakeO := NewOrchestrator(fakeCfgManager)
				fakeO.prntr = fakePrinter

				err := run(fakeO, ctx)
				assert.Nil(t, err, "expected edit to succeed")
				assert.Equal(t, 1, printCfgCallCount, "expected to be called exactly once")
			})
		})
	})
}

var FailToUnmarshalConfigYaml = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			ctx.(common.ConfigSetter).SetConfig(nil)
			fakeCfgManager := config.CreateFakeConfigManager()
			fakeO := NewOrchestrator(fakeCfgManager)
			err := run(fakeO, ctx)
			assert.NotNil(t, err, "expected config unmarshal to fail")
		})
	})
}
