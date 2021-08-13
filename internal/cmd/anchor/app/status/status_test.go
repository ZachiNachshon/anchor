package status

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StatusActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete status flow successfully",
			Func: CompleteStatusFlowSuccessfully,
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
			Name: "print applications with missing instructions status",
			Func: PrintApplicationsWithMissingInstructionsStatus,
		},
		{
			Name: "print applications with valid status",
			Func: PrintApplicationsWithValidStatus,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteStatusFlowSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeOrch := createFakeOrchestrator()
			bannerCallCount := 0
			fakeOrch.bannerMock = func() {
				bannerCallCount++
			}
			prepareCallCount := 0
			fakeOrch.prepareMock = func(ctx common.Context) error {
				prepareCallCount++
				return nil
			}
			runCallCount := 0
			fakeOrch.runMock = func(ctx common.Context) error {
				runCallCount++
				return nil
			}
			err := AppStatus(ctx, fakeOrch)
			assert.Nil(t, err, "expected not to fail app status")
			assert.Equal(t, 1, bannerCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var PrintBanner = func(t *testing.T) {
	fakePrinter := printer.CreateFakePrinter()
	printBannerCallCount := 0
	fakePrinter.PrintAnchorBannerMock = func() {
		printBannerCallCount++
	}

	statusOrch := &statusOrchestratorImpl{
		p: fakePrinter,
	}
	statusOrch.banner()
	assert.Equal(t, 1, printBannerCallCount, "expected func to be called exactly once")
}

var PrepareRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		path := "/some/path"

		fakeLocator := locator.CreateFakeLocator(path)
		reg.Set(locator.Identifier, fakeLocator)

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		fakeExtractor := extractor.CreateFakeExtractor()
		reg.Set(extractor.Identifier, fakeExtractor)

		fakeParser := parser.CreateFakeParser()
		reg.Set(parser.Identifier, fakeParser)

		statusOrch := &statusOrchestratorImpl{}
		err := statusOrch.prepare(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, statusOrch.l)
		assert.NotNil(t, statusOrch.p)
		assert.NotNil(t, statusOrch.e)
		assert.NotNil(t, statusOrch.pa)
	})
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		path := "/some/path"

		fakeLocator := locator.CreateFakeLocator(path)
		fakePrinter := printer.CreateFakePrinter()
		fakeExtractor := extractor.CreateFakeExtractor()
		fakeParser := parser.CreateFakeParser()
		statusOrch := &statusOrchestratorImpl{}

		err := statusOrch.prepare(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", locator.Identifier))
		reg.Set(locator.Identifier, fakeLocator)

		err = statusOrch.prepare(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", printer.Identifier))
		reg.Set(printer.Identifier, fakePrinter)

		err = statusOrch.prepare(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", extractor.Identifier))
		reg.Set(extractor.Identifier, fakeExtractor)

		err = statusOrch.prepare(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", parser.Identifier))
		reg.Set(parser.Identifier, fakeParser)

		err = statusOrch.prepare(ctx)
		assert.Nil(t, err)
	})
}

var PrintApplicationsWithMissingInstructionsStatus = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		path := "/some/path"
		apps := stubs.GenerateApplicationTestData()

		fakeLocator := locator.CreateFakeLocator(path)
		reg.Set(locator.Identifier, fakeLocator)
		locateAppsCallCount := 0
		fakeLocator.ApplicationsMock = func() []*models.ApplicationInfo {
			locateAppsCallCount++
			return apps
		}

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)
		printAppsCallCount := 0
		fakePrinter.PrintApplicationsMock = func(apps []*printer.AppStatusTemplateItem) {
			printAppsCallCount++
			assert.Equal(t, 2, len(apps))
			for _, app := range apps {
				assert.True(t, app.MissingInstructionFile)
				assert.False(t, app.IsValid)
			}
		}

		statusOrch := &statusOrchestratorImpl{
			l: fakeLocator,
			p: fakePrinter,
		}

		err := statusOrch.run(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 1, locateAppsCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, printAppsCallCount, "expected func to be called exactly once")
	})
}

var PrintApplicationsWithValidStatus = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		with.HarnessAnchorfilesTestRepo(ctx)

		path := "/some/path"
		apps := stubs.GenerateApplicationTestData()
		apps[0].InstructionsPath = ctx.AnchorFilesPath() + "/app/first-app/instructions.yaml"
		apps[1].InstructionsPath = ctx.AnchorFilesPath() + "/app/second-app/instructions.yaml"

		fakeLocator := locator.CreateFakeLocator(path)
		reg.Set(locator.Identifier, fakeLocator)
		locateAppsCallCount := 0
		fakeLocator.ApplicationsMock = func() []*models.ApplicationInfo {
			locateAppsCallCount++
			return apps
		}

		fakeExtractor := extractor.CreateFakeExtractor()
		fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
			assert.True(t, instructionsPath == apps[0].InstructionsPath ||
				instructionsPath == apps[1].InstructionsPath)
			return &models.InstructionsRoot{}, nil
		}
		reg.Set(extractor.Identifier, fakeExtractor)

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)
		printAppsCallCount := 0
		fakePrinter.PrintApplicationsMock = func(apps []*printer.AppStatusTemplateItem) {
			printAppsCallCount++
			assert.Equal(t, 2, len(apps))
			for _, app := range apps {
				assert.False(t, app.MissingInstructionFile)
				assert.True(t, app.IsValid)
			}
		}

		statusOrch := &statusOrchestratorImpl{
			l: fakeLocator,
			p: fakePrinter,
			e: fakeExtractor,
		}

		err := statusOrch.run(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 1, locateAppsCallCount, "expected func to be called exactly once")
		assert.Equal(t, 1, printAppsCallCount, "expected func to be called exactly once")
	})
}
