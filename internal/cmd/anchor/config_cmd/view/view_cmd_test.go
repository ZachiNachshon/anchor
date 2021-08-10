package view

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

func Test_ViewCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start view action successfully",
			Func: StartViewActionSuccessfully,
		},
		{
			Name: "fail view action",
			Func: FailViewAction,
		},
	}
	harness.RunTests(t, tests)
}

var StartViewActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				fakeCfgManager := config.CreateFakeConfigManager()
				callCount := 0
				fun := func(ctx common.Context, cfgManager config.ConfigManager) error {
					callCount++
					return nil
				}
				command, err := NewCommand(ctx, fakeCfgManager, fun)
				out, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: print")
				assert.Nil(t, err, "expected cli action to succeed")
				logger.Info(out)
			})
		})
	})
}

var FailViewAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, cfgManager config.ConfigManager) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command, err := NewCommand(ctx, fakeCfgManager, fun)
				_, err = drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: print")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}
