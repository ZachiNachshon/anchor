package list

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/test"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ListCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		//{
		//	Name: "call print func exactly once",
		//	Func: CallPrintFuncExactlyOnce,
		//},
	}
	harness.RunTests(t, tests)
}

var CallPrintFuncExactlyOnce = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := test.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config config.AnchorConfig) {
				scanCallCount := 0
				fake := locator.FakeLocatorLoader(ctx.AnchorFilesPath())
				fake.ScanMock = func() error {
					scanCallCount += 1
					return nil
				}
				appsCallCount := 0
				fake.ApplicationsMock = func() []*models.AppContent {
					appsCallCount += 1
					return nil
				}
				locator.ToRegistry(ctx.Registry(), fake)
				cmd := NewCommand(ctx)
				if _, err := drivers.CLI().RunCommand(cmd); err != nil {
					assert.Failf(t, err.Error(), err.Error())
				} else {
					assert.Equal(t, scanCallCount, 1, "expected print func to be called exactly once")
				}
			})
		})
	})
}
