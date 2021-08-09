package _select

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

func Test_SelectCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start select action successfully",
			Func: StartSelectActionSuccessfully,
		},
		{
			Name: "fail select action",
			Func: FailSelectAction,
		},
	}
	harness.RunTests(t, tests)
}

var StartSelectActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return nil
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, fun))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: select")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailSelectAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, fun))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: select")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}
