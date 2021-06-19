package install

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/test"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/app"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InstallCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start the installation flow successfully",
			Func: StartTheInstallationFlowSuccessfully,
		},
		{
			Name: "fail command when install action failed",
			Func: FailCommandWhenInstallActionFailed,
		},
	}
	harness.RunTests(t, tests)
}

var StartTheInstallationFlowSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, test.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &app.ApplicationActions{
					Install: func(ctx common.Context) error {
						callCount++
						return nil
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: install")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailCommandWhenInstallActionFailed = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, test.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &app.ApplicationActions{
					Install: func(ctx common.Context) error {
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
