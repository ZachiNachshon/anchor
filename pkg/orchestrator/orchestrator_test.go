package orchestrator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
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
		{
			Name: "fail to prompt for application selection",
			Func: FailToPromptForApplicationSelection,
		},
		{
			Name: "fail to extract instructions and prompt applications selection again",
			Func: FailToExtractInstructionAndPromptAppSelectionAgain,
		},
		{
			Name: "fail to prompt for instruction selection",
			Func: FailToPromptForInstructionSelection,
		},
	}
	harness.RunTests(t, tests)
}

var ExitAppsPromptMenuOnCancelButton = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				appsPromptCallCount++
				return prompter.GetAppByName(appsArr, prompter.CancelButtonName), nil
			}

			orchestrator := New(fakePrompter,
				fakeLocator,
				extractor.CreateFakeExtractor(),
				parser.CreateFakeParser())

			item, err := orchestrator.OrchestrateAppInstructionSelection()
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.Equal(t, 1, locateAppsCallCount)
			assert.Equal(t, 1, appsPromptCallCount)
			assert.EqualValues(t, prompter.CancelButtonName, item.Id)
		})
	})
}

var PerformBasicPromptFromAppToInstructionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instTestData := stubs.GenerateInstructionsTestData()

			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				appsPromptCallCount++
				return prompter.GetAppByName(appsArr, stubs.App1Name), nil
			}

			instructionsPromptCallCount := 0
			fakePrompter.PromptInstructionsMock = func(appName string, instructions *models.Instructions) (*models.PromptItem, error) {
				instructionsPromptCallCount++
				return prompter.GetInstructionItemById(instTestData, stubs.App1InstructionsItem1Id), nil
			}

			fakeExtractor := extractor.CreateFakeExtractor()
			extractorCallCount := 0
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.Instructions, error) {
				extractorCallCount++
				return instTestData, nil
			}

			orchestrator := New(fakePrompter,
				fakeLocator,
				fakeExtractor,
				parser.CreateFakeParser())

			item, err := orchestrator.OrchestrateAppInstructionSelection()
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.Equal(t, 1, locateAppsCallCount)
			assert.Equal(t, 1, appsPromptCallCount)
			assert.Equal(t, 1, instructionsPromptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.EqualValues(t, item, prompter.GetInstructionItemById(instTestData, stubs.App1InstructionsItem1Id))
		})
	})
}

var GoBackFromInstructionsToAppsPrompMenuSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instTestData := stubs.GenerateInstructionsTestData()

			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				appsPromptCallCount++
				if appsPromptCallCount == 1 {
					return prompter.GetAppByName(appsArr, stubs.App1Name), nil
				} else if appsPromptCallCount == 2 {
					return prompter.GetAppByName(appsArr, prompter.CancelButtonName), nil
				}
				return nil, fmt.Errorf("bad test flow")
			}

			instructionsPromptCallCount := 0
			fakePrompter.PromptInstructionsMock = func(appName string, instructions *models.Instructions) (*models.PromptItem, error) {
				instructionsPromptCallCount++
				return prompter.GetInstructionItemById(instTestData, prompter.BackButtonName), nil
			}

			fakeExtractor := extractor.CreateFakeExtractor()
			extractorCallCount := 0
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.Instructions, error) {
				extractorCallCount++
				return instTestData, nil
			}

			orchestrator := New(fakePrompter,
				fakeLocator,
				fakeExtractor,
				parser.CreateFakeParser())

			item, err := orchestrator.OrchestrateAppInstructionSelection()
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.Equal(t, 2, locateAppsCallCount)
			assert.Equal(t, 2, appsPromptCallCount)
			assert.Equal(t, 1, instructionsPromptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.EqualValues(t, item.Id, prompter.CancelButtonName)
		})
	})
}

var FailToPromptForApplicationSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				appsPromptCallCount++
				return nil, fmt.Errorf("failed to prompt for app selection")
			}

			orchestrator := New(fakePrompter,
				fakeLocator,
				extractor.CreateFakeExtractor(),
				parser.CreateFakeParser())

			item, err := orchestrator.OrchestrateAppInstructionSelection()
			assert.NotNil(t, err, "expected orchestrator to fail")
			assert.Equal(t, "failed to prompt for app selection", err.Error())
			assert.Equal(t, 1, locateAppsCallCount)
			assert.Equal(t, 1, appsPromptCallCount)
			assert.Nil(t, item, "expected not to have return value")
		})
	})
}

var FailToExtractInstructionAndPromptAppSelectionAgain = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				appsPromptCallCount++
				if appsPromptCallCount == 1 {
					return prompter.GetAppByName(appsArr, stubs.App1Name), nil
				} else if appsPromptCallCount == 2 {
					return prompter.GetAppByName(appsArr, prompter.CancelButtonName), nil
				}
				return nil, fmt.Errorf("bad test flow")
			}

			fakeExtractor := extractor.CreateFakeExtractor()
			extractorCallCount := 0
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.Instructions, error) {
				extractorCallCount++
				return nil, fmt.Errorf("failed to extract instructions")
			}

			orchestrator := New(fakePrompter,
				fakeLocator,
				fakeExtractor,
				parser.CreateFakeParser())

			item, err := orchestrator.OrchestrateAppInstructionSelection()
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.Equal(t, 2, locateAppsCallCount)
			assert.Equal(t, 2, appsPromptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.EqualValues(t, item.Id, prompter.CancelButtonName)
		})
	})
}

var FailToPromptForInstructionSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instTestData := stubs.GenerateInstructionsTestData()

			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.AppContent {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.AppContent) (*models.AppContent, error) {
				appsPromptCallCount++
				return prompter.GetAppByName(appsArr, stubs.App1Name), nil
			}

			instructionsPromptCallCount := 0
			fakePrompter.PromptInstructionsMock = func(appName string, instructions *models.Instructions) (*models.PromptItem, error) {
				instructionsPromptCallCount++
				return nil, fmt.Errorf("failed to prompt for instructions")
			}

			fakeExtractor := extractor.CreateFakeExtractor()
			extractorCallCount := 0
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.Instructions, error) {
				extractorCallCount++
				return instTestData, nil
			}

			orchestrator := New(fakePrompter,
				fakeLocator,
				fakeExtractor,
				parser.CreateFakeParser())

			item, err := orchestrator.OrchestrateAppInstructionSelection()
			assert.Nil(t, item, "expected not to have return value")
			assert.NotNil(t, err, "expected orchestrator to fail")
			assert.Equal(t, "failed to prompt for instructions", err.Error())
			assert.Equal(t, 1, locateAppsCallCount)
			assert.Equal(t, 1, appsPromptCallCount)
			assert.Equal(t, 1, instructionsPromptCallCount)
			assert.Equal(t, 1, extractorCallCount)
		})
	})
}
