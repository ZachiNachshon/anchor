package use_context

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

func Test_UseContextCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start use context action successfully",
			Func: StartUseContextActionSuccessfully,
		},
		{
			Name: "fail due to missing config context name",
			Func: FailDueToMissingConfigContextName,
		},
		{
			Name: "fail use context action",
			Func: FailUseContextAction,
		},
	}
	harness.RunTests(t, tests)
}

var StartUseContextActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				configContextName := "test-cfg-context"
				callCount := 0
				actions := &cfg.ConfigurationActions{
					UseContext: func(ctx common.Context, cfgCtxName string) error {
						callCount++
						return nil
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions), configContextName)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: use-context")
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
					UseContext: func(ctx common.Context, cfgCtxName string) error {
						callCount++
						return nil
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions))
				assert.Equal(t, 0, callCount, "expected action not to be called. name: use-context")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, "accepts 1 arg(s), received 0", err.Error())
			})
		})
	})
}

var FailUseContextAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				configContextName := "test-cfg-context"
				callCount := 0
				actions := &cfg.ConfigurationActions{
					UseContext: func(ctx common.Context, cfgCtxName string) error {
						callCount++
						return fmt.Errorf("an error occurred")
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions), configContextName)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: use-context")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}
