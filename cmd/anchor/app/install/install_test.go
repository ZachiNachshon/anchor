package install

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/test"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InstallCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail on illegal amount of args",
			Func: FailOnIllegalAmountOfArgs,
		},
	}
	harness.RunTests(t, tests)
}

var FailOnIllegalAmountOfArgs = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, test.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				cmd := NewCommand(ctx)
				if _, err := drivers.CLI().RunCommand(cmd); err != nil {
					assert.Equal(t, err.Error(), "accepts 1 arg(s), received 0")
				} else {
					assert.Fail(t, "expected to fail on invalid arguments count")
				}
			})
		})
	})
}
