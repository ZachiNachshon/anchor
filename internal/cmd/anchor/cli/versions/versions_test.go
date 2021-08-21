package versions

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VersionActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
		{
			Name: "fail runner due to preparation",
			Func: FailRunnerDueToPreparation,
		},
		{
			Name: "print banner",
			Func: PrintBanner,
		},
		{
			Name: "prepare registry components",
			Func: PrepareRegistryComponents,
		},
		{
			Name: "fail resolving registry components",
			Func: FailResolvingRegistryComponents,
		},
		{
			Name: "print cli versions",
			Func: PrintCliVersions,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeO := NewOrchestrator()
			bannerCallCount := 0
			fakeO.bannerFunc = func(o *versionsOrchestrator) {
				bannerCallCount++
			}
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *versionsOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return nil
			}
			runCallCount := 0
			fakeO.runFunc = func(o *versionsOrchestrator, ctx common.Context) error {
				runCallCount++
				return nil
			}
			err := CliVersions(ctx, fakeO)
			assert.Nil(t, err, "expected not to fail")
			assert.Equal(t, 1, bannerCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var FailRunnerDueToPreparation = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeO := NewOrchestrator()
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *versionsOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return fmt.Errorf("failed to prepare runner")
			}
			err := CliVersions(ctx, fakeO)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare runner", err.Error())
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
		})
	})
}

var PrintBanner = func(t *testing.T) {
	fakePrinter := printer.CreateFakePrinter()
	printBannerCallCount := 0
	fakePrinter.PrintAnchorBannerMock = func() {
		printBannerCallCount++
	}

	fakeO := NewOrchestrator()
	fakeO.prntr = fakePrinter
	fakeO.bannerFunc(fakeO)
	assert.Equal(t, 1, printBannerCallCount, "expected func to be called exactly once")
}

var PrepareRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		fakeO := NewOrchestrator()
		err := fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
		assert.NotNil(t, fakeO.prntr)
	})
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		fakePrinter := printer.CreateFakePrinter()
		fakeO := NewOrchestrator()

		err := fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", printer.Identifier))
		reg.Set(printer.Identifier, fakePrinter)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
	})
}

var PrintCliVersions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeO := NewOrchestrator()
		err := fakeO.runFunc(fakeO, ctx)
		assert.Nil(t, err)
	})
}
