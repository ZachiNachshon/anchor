package view

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/cfg"
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &cfg.ConfigurationActions{
					View: func(ctx common.Context) error {
						callCount++
						return nil
					},
				}
				out, err := drivers.CLI().RunCommand(NewCommand(ctx, actions))
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &cfg.ConfigurationActions{
					View: func(ctx common.Context) error {
						callCount++
						return fmt.Errorf("an error occurred")
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: print")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}
