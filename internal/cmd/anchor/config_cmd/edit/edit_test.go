package edit

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EditActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
		{
			Name: "fail to resolve shell from registry",
			Func: FailToResolveShellFromRegistry,
		},
		{
			Name: "fail to edit configuration file",
			Func: FailToEditConfigurationFile,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeShell := shell.CreateFakeShell()
			executeCallCount := 0
			fakeShell.ExecuteTTYMock = func(script string) error {
				executeCallCount++
				return nil
			}
			ctx.Registry().Set(shell.Identifier, fakeShell)

			fakeCfgManager := config.CreateFakeConfigManager()
			fakeCfgManager.GetConfigFilePathMock = func() (string, error) {
				return "/path/to/config", nil
			}

			err := ConfigEdit(ctx, fakeCfgManager)
			assert.Nil(t, err, "expected edit to succeed")
			assert.Equal(t, 1, executeCallCount, "expected to be called exactly once")
		})
	})
}

var FailToResolveShellFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeCfgManager := config.CreateFakeConfigManager()
			fakeCfgManager.GetConfigFilePathMock = func() (string, error) {
				return "/path/to/config", nil
			}
			err := ConfigEdit(ctx, fakeCfgManager)
			assert.NotNil(t, err, "expected to fail on missing registry item")
		})
	})
}

var FailToEditConfigurationFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteTTYMock = func(script string) error {
				return fmt.Errorf("failed to edit")
			}
			ctx.Registry().Set(shell.Identifier, fakeShell)

			fakeCfgManager := config.CreateFakeConfigManager()
			fakeCfgManager.GetConfigFilePathMock = func() (string, error) {
				return "/path/to/config", nil
			}
			err := ConfigEdit(ctx, fakeCfgManager)
			assert.NotNil(t, err, "expected edit to fail")
			assert.Equal(t, "failed to edit", err.Error())
		})
	})
}
