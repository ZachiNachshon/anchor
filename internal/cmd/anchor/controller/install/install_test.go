package install

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InstallActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			err := ControllerInstall(ctx)
			assert.Nil(t, err, "expected runner method not to fail")
		})
	})
}
