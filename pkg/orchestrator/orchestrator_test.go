package orchestrator

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	extractor2 "github.com/ZachiNachshon/anchor/pkg/extractor"
	locator2 "github.com/ZachiNachshon/anchor/pkg/locator"
	parser2 "github.com/ZachiNachshon/anchor/pkg/parser"
	prompter2 "github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OrchestratorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "exit apps prompt menu on cancel button",
			Func: ExitAppsPromptMenuOnCancelButton,
		},
		{
			Name: "perform basic prompt from app to instruction successfully",
			Func: PerformBasicPromptFromAppToInstructionSuccessfully,
		},
		{
			Name: "go back from instructions to apps prompt menu successfully",
			Func: GoBackFromInstructionsToAppsPrompMenuSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var ExitAppsPromptMenuOnCancelButton = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			registry := ctx.Registry()

			// Given I create a locator to return test data
			fakeLocator := locator2.CreateFakeLocator(ctx.AnchorFilesPath())
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				return stubs.GenerateApplicationTestData()
			}
			locator2.ToRegistry(registry, fakeLocator)

			// And I create an apps prompter
			fakePrompter := prompter2.CreateFakePrompter()
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				return appsArr[0], nil
			}
			prompter2.ToRegistry(registry, fakePrompter)

			// And I create a dummy extractor
			fakeExtractor := extractor2.CreateFakeExtractor()
			extractor2.ToRegistry(registry, fakeExtractor)

			// And I create a dummy parser
			fakeParser := parser2.CreateFakeParser()
			parser2.ToRegistry(registry, fakeParser)

			// When I create a new orchestrator
			orchestrator := New(fakePrompter, fakeLocator, fakeExtractor, fakeParser)
			item, err := orchestrator.OrchestrateAppInstructionSelection()

			// Then I expect the result item to represent a cancel selection
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.EqualValues(t, prompter2.CancelButtonName, item.Id)
		})
	})
}

var PerformBasicPromptFromAppToInstructionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			registry := ctx.Registry()
			instTestData := stubs.GenerateInstructionsTestData()

			// Given I create a locator to return test data
			fakeLocator := locator2.CreateFakeLocator(ctx.AnchorFilesPath())
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				return stubs.GenerateApplicationTestData()
			}
			locator2.ToRegistry(registry, fakeLocator)

			// And I create an apps prompter that selects the 1st test app and 1st test instruction
			fakePrompter := prompter2.CreateFakePrompter()
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				return appsArr[1], nil
			}
			fakePrompter.PromptInstructionsMock = func(instructions *models.Instructions) (*models.PromptItem, error) {
				return instTestData.Items[1], nil
			}
			prompter2.ToRegistry(registry, fakePrompter)

			// And I create a dummy extractor
			fakeExtractor := extractor2.CreateFakeExtractor()
			fakeExtractor.ExtractPromptItemsMock = func(instructionsPath string, p parser2.Parser) (*models.Instructions, error) {
				return instTestData, nil
			}
			extractor2.ToRegistry(registry, fakeExtractor)

			// And I create a dummy parser
			fakeParser := parser2.CreateFakeParser()
			parser2.ToRegistry(registry, fakeParser)

			// When I create a new orchestrator
			orchestrator := New(fakePrompter, fakeLocator, fakeExtractor, fakeParser)
			item, err := orchestrator.OrchestrateAppInstructionSelection()

			// Then I expect the result item to represent a mocked selection
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.EqualValues(t, item, instTestData.Items[1])
		})
	})
}

var GoBackFromInstructionsToAppsPrompMenuSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			registry := ctx.Registry()
			instTestData := stubs.GenerateInstructionsTestData()

			// Given I create a locator to return test data
			fakeLocator := locator2.CreateFakeLocator(ctx.AnchorFilesPath())
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				return stubs.GenerateApplicationTestData()
			}
			locator2.ToRegistry(registry, fakeLocator)

			// And I create an apps prompter that selects the 1st test app and 1st test instruction
			fakePrompter := prompter2.CreateFakePrompter()
			appsPromptMenuCount := 1
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				if appsPromptMenuCount == 1 {
					appsPromptMenuCount++
					return appsArr[1], nil
				} else {
					// The 2nd apps prompt menu should choose cancel option
					return appsArr[0], nil
				}
			}
			fakePrompter.PromptInstructionsMock = func(instructions *models.Instructions) (*models.PromptItem, error) {
				return instTestData.Items[0], nil
			}
			prompter2.ToRegistry(registry, fakePrompter)

			// And I create a dummy extractor
			fakeExtractor := extractor2.CreateFakeExtractor()
			fakeExtractor.ExtractPromptItemsMock = func(instructionsPath string, p parser2.Parser) (*models.Instructions, error) {
				return instTestData, nil
			}
			extractor2.ToRegistry(registry, fakeExtractor)

			// And I create a dummy parser
			fakeParser := parser2.CreateFakeParser()
			parser2.ToRegistry(registry, fakeParser)

			// When I create a new orchestrator
			orchestrator := New(fakePrompter, fakeLocator, fakeExtractor, fakeParser)
			item, err := orchestrator.OrchestrateAppInstructionSelection()

			// Then I expect to go back successfully from instructions to apps prompt menu and select the cancel option
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.EqualValues(t, item.Id, prompter2.CancelButtonName)
		})
	})
}
