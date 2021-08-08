package versions

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

func Test_VersionsCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start versions action successfully",
			Func: StartVersionsActionSuccessfully,
		},
		{
			Name: "fail versions action",
			Func: FailVersionsAction,
		},
	}
	harness.RunTests(t, tests)
}

var StartVersionsActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return nil
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, fun))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: versions")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailVersionsAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, fun))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: versions")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}
