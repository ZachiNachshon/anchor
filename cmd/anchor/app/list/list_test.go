package list

import (
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

func Test_ListCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start the list flow successfully",
			Func: StartTheListFlowSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var StartTheListFlowSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, test.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &app.ApplicationActions{
					List: func(ctx common.Context) error {
						callCount++
						return nil
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions))
				assert.Equal(t, 1, callCount, "expected list action to be called exactly once")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

//var ScanAndListAllAppsSuccessfully = func(t *testing.T) {
//	with.Context(func(ctx common.Context) {
//		with.Logging(ctx, t, func(logger logger.Logger) {
//			yamlConfigText := test.GetDefaultTestConfigText()
//			with.Config(ctx, yamlConfigText, func(config config.AnchorConfig) {
//				scanCallCount := 0
//				fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
//				fakeLocator.ScanMock = func() error {
//					scanCallCount += 1
//					return nil
//				}
//				appsCallCount := 0
//				fakeLocator.ApplicationsMock = func() []*models.AppContent {
//					appsCallCount += 1
//					return nil
//				}
//				locator.ToRegistry(ctx.Registry(), fakeLocator)
//
//				printCallCount := 0
//				fakePrinter := printer.CreateFakePrinter()
//				fakePrinter.PrintApplicationsMock = func(apps []*models.AppContent) {
//					printCallCount += 1
//				}
//				printer.ToRegistry(ctx.Registry(), fakePrinter)
//
//				if _, err := drivers.CLI().RunCommand(NewCommand(ctx)); err != nil {
//					assert.Failf(t, err.Error(), err.Error())
//				} else {
//					assert.Equal(t, 1, scanCallCount, "expected locator scan to be called exactly once")
//					assert.Equal(t, 1, appsCallCount, "expected locator get all apps func to be called exactly once")
//					assert.Equal(t, 1, printCallCount, "expected printer print to be called exactly once")
//				}
//			})
//		})
//	})
//}
