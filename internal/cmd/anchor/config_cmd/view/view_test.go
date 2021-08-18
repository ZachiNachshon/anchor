package view

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VersionActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		//{
		//	Name: "complete runner method successfully",
		//	Func: CompleteRunnerMethodSuccessfully,
		//},
		{
			Name: "fail to resolve printer from registry",
			Func: FailToResolvePrinterFromRegistry,
		},
		{
			Name: "print configuration successfully",
			Func: PrintConfigurationSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteTTYMock = func(script string) error {
				return nil
			}
			ctx.Registry().Set(shell.Identifier, fakeShell)

			fakeCfgManager := config.CreateFakeConfigManager()
			fakeCfgManager.GetConfigFilePathMock = func() (string, error) {
				return "/path/to/config", nil
			}

			err := ConfigView(ctx, fakeCfgManager)
			assert.Nil(t, err, "expected edit to succeed")
		})
	})
}

var FailToResolvePrinterFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			err := printConfiguration(ctx, "/path/to/config", "config-yaml-text")
			assert.NotNil(t, err, "expected to fail on missing registry item")
		})
	})
}

var PrintConfigurationSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakePrinter := printer.CreateFakePrinter()
			callCount := 0
			fakePrinter.PrintConfigurationMock = func(cfgFilePath string, cfgText string) {
				callCount++
				return
			}
			ctx.Registry().Set(printer.Identifier, fakePrinter)
			err := printConfiguration(ctx, "/path/to/config", "config-yaml-text")
			assert.Nil(t, err, "expected print to succeed")
			assert.Equal(t, 1, callCount, "expected to be called exactly once")
		})
	})
}
