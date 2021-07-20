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
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OrchestratorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail to prompt for application selection",
			Func: FailToPromptForApplicationSelection,
		},
		{
			Name: "exit apps prompt menu on cancel button",
			Func: ExitAppsPromptMenuOnCancelButton,
		},
		{
			Name: "select application successfully",
			Func: SelectApplicationSuccessfully,
		},
		{
			Name: "fail to extract instruction",
			Func: FailToExtractInstruction,
		},
		{
			Name: "fail to prompt for instructions",
			Func: FailToPromptForInstructions,
		},
		{
			Name: "select instruction successfully",
			Func: SelectInstructionSuccessfully,
		},
		{
			Name: "fail to ask for user input before running instruction",
			Func: FailToAskForUserInputBeforeRunningInstruction,
		},
		{
			Name: "ask for user input before running instructions successfully",
			Func: AskForUserInputBeforeRunningInstructionSuccessfully,
		},
		{
			Name: "run instruction successfully",
			Func: RunInstructionSuccessfully,
		},
		{
			Name: "failed to run instruction due to script execution",
			Func: FailedToRunInstructionDueToScriptExecution,
		},
		{
			Name: "failed to prompt for key press after instruction run",
			Func: FailedToPromptForKeyPressAfterInstructionRun,
		},
		{
			Name: "failed to clear screen after instruction run",
			Func: FailedToClearScreenAfterInstructionRun,
		},
	}
	harness.RunTests(t, tests)
}

var FailToPromptForApplicationSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.ApplicationInfo {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
				appsPromptCallCount++
				return nil, fmt.Errorf("failed to prompt for app selection")
			}

			orchestrator := New(fakePrompter, fakeLocator, nil, nil, nil, nil)
			item, err := orchestrator.OrchestrateApplicationSelection()
			assert.NotNil(t, err, "expected orchestrator to fail")
			assert.Equal(t, "failed to prompt for app selection", err.GoError().Error())
			assert.Equal(t, 1, locateAppsCallCount)
			assert.Equal(t, 1, appsPromptCallCount)
			assert.Nil(t, item, "expected not to have return value")
		})
	})
}

var ExitAppsPromptMenuOnCancelButton = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.ApplicationInfo {
				locateAppsCallCount++
				return stubs.GenerateApplicationTestData()
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
				appsPromptCallCount++
				return stubs.GetAppByName(appsArr, prompter.CancelActionName), nil
			}

			orchestrator := New(fakePrompter, fakeLocator, nil, nil, nil, nil)
			item, err := orchestrator.OrchestrateApplicationSelection()
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.Equal(t, 1, locateAppsCallCount)
			assert.Equal(t, 1, appsPromptCallCount)
			assert.EqualValues(t, prompter.CancelActionName, item.Name)
		})
	})
}

var SelectApplicationSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)

			fakeLocator := locator.CreateFakeLocator(ctx.AnchorFilesPath())
			locateAppsCallCount := 0
			fakeLocator.ApplicationsMock = func() []*models.ApplicationInfo {
				locateAppsCallCount++
				return apps
			}

			fakePrompter := prompter.CreateFakePrompter()
			appsPromptCallCount := 0
			fakePrompter.PromptAppsMock = func(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
				appsPromptCallCount++
				return app1, nil
			}

			orchestrator := New(fakePrompter, fakeLocator, nil, nil, nil, nil)
			item, err := orchestrator.OrchestrateApplicationSelection()
			assert.Nil(t, err, "expected orchestrator to exit successfully")
			assert.Equal(t, 1, locateAppsCallCount)
			assert.Equal(t, 1, appsPromptCallCount)
			assert.EqualValues(t, item.Name, app1.Name)
		})
	})
}

var FailToExtractInstruction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)

			fakeExtractor := extractor.CreateFakeExtractor()
			extractorCallCount := 0
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractorCallCount++
				assert.Equal(t, instructionsPath, app1.InstructionsPath)
				return nil, fmt.Errorf("failed to extract instructions")
			}

			orchestrator := New(nil, nil, fakeExtractor, nil, nil, nil)
			instItem, err := orchestrator.OrchestrateInstructionSelection(app1)
			assert.Nil(t, instItem, "expected instruction to be empty")
			assert.NotNil(t, err, "expected instruction selection to fail")
			assert.Equal(t, 1, extractorCallCount)
		})
	})
}

var FailToPromptForInstructions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			instTestData := stubs.GenerateInstructionsTestData()

			fakeExtractor := extractor.CreateFakeExtractor()
			extractorCallCount := 0
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractorCallCount++
				assert.Equal(t, instructionsPath, app1.InstructionsPath)
				return instTestData, nil
			}

			fakePrompter := prompter.CreateFakePrompter()
			instPromptCallCount := 0
			fakePrompter.PromptInstructionsMock = func(appName string, instructionsRoot *models.InstructionsRoot) (*models.Action, error) {
				instPromptCallCount++
				return nil, fmt.Errorf("failed to prompt for instructions")
			}

			orchestrator := New(fakePrompter, nil, fakeExtractor, nil, nil, nil)
			instItem, err := orchestrator.OrchestrateInstructionSelection(app1)
			assert.Nil(t, instItem, "expected instruction to be empty")
			assert.NotNil(t, err, "expected instruction selection to fail")
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instPromptCallCount)
		})
	})
}

var SelectInstructionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeExtractor := extractor.CreateFakeExtractor()
			extractorCallCount := 0
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractorCallCount++
				assert.Equal(t, instructionsPath, app1.InstructionsPath)
				return instRootTestData, nil
			}

			fakePrompter := prompter.CreateFakePrompter()
			instPromptCallCount := 0
			fakePrompter.PromptInstructionsMock = func(appName string, instructionsRoot *models.InstructionsRoot) (*models.Action, error) {
				instPromptCallCount++
				return inst1, nil
			}

			orchestrator := New(fakePrompter, nil, fakeExtractor, nil, nil, nil)
			instItem, err := orchestrator.OrchestrateInstructionSelection(app1)
			assert.NotNil(t, instItem, "expected instruction not to be empty")
			assert.Nil(t, err, "expected instruction selection not to fail")
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instPromptCallCount)
		})
	})
}

var FailToAskForUserInputBeforeRunningInstruction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeUserInput := input.CreateFakeUserInput()
			userInputCallCount := 0
			fakeUserInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				userInputCallCount++
				return false, fmt.Errorf("failed to ask yes/no question")
			}

			orchestrator := New(nil, nil, nil, nil, nil, fakeUserInput)
			shouldRun, err := orchestrator.AskBeforeRunningInstruction(inst1)
			assert.Equal(t, false, shouldRun)
			assert.NotNil(t, err, "expected instruction selection to fail")
			assert.Equal(t, 1, userInputCallCount)
		})
	})
}

var AskForUserInputBeforeRunningInstructionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeUserInput := input.CreateFakeUserInput()
			userInputCallCount := 0
			fakeUserInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				userInputCallCount++
				return true, nil
			}

			orchestrator := New(nil, nil, nil, nil, nil, fakeUserInput)
			shouldRun, err := orchestrator.AskBeforeRunningInstruction(inst1)
			assert.Equal(t, true, shouldRun)
			assert.Nil(t, err, "expected instruction selection to succeed")
			assert.Equal(t, 1, userInputCallCount)
		})
	})
}

var RunInstructionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeShell := shell.CreateFakeShell()
			execScriptCallCount := 0
			fakeShell.ExecuteScriptWithOutputToFileMock = func(dir string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execScriptCallCount++
				assert.Equal(t, relativeScriptPath, inst1.File)
				return nil
			}

			clearScreenCallCount := 0
			fakeShell.ClearScreenMock = func() error {
				clearScreenCallCount++
				return nil
			}

			fakeUserInput := input.CreateFakeUserInput()
			pressAnyKeyCallCount := 0
			fakeUserInput.PressAnyKeyToContinueMock = func() error {
				pressAnyKeyCallCount++
				return nil
			}

			orchestrator := New(nil, nil, nil, nil, fakeShell, fakeUserInput)
			err := orchestrator.RunInstruction(inst1, ctx.AnchorFilesPath())
			assert.Nil(t, err, "expected instruction run to succeed")
			assert.Equal(t, 1, execScriptCallCount)
			assert.Equal(t, 1, pressAnyKeyCallCount)
			assert.Equal(t, 1, clearScreenCallCount)
		})
	})
}

var FailedToRunInstructionDueToScriptExecution = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeShell := shell.CreateFakeShell()
			execScriptCallCount := 0
			fakeShell.ExecuteScriptWithOutputToFileMock = func(dir string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execScriptCallCount++
				assert.Equal(t, relativeScriptPath, inst1.File)
				return fmt.Errorf("failed to execute script")
			}

			orchestrator := New(nil, nil, nil, nil, fakeShell, nil)
			err := orchestrator.RunInstruction(inst1, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected instruction run to fail")
			assert.Equal(t, "failed to execute script", err.GoError().Error())
			assert.Equal(t, 1, execScriptCallCount)
		})
	})
}

var FailedToPromptForKeyPressAfterInstructionRun = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeShell := shell.CreateFakeShell()
			execScriptCallCount := 0
			fakeShell.ExecuteScriptWithOutputToFileMock = func(dir string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execScriptCallCount++
				assert.Equal(t, relativeScriptPath, inst1.File)
				return nil
			}

			fakeUserInput := input.CreateFakeUserInput()
			pressAnyKeyCallCount := 0
			fakeUserInput.PressAnyKeyToContinueMock = func() error {
				pressAnyKeyCallCount++
				return fmt.Errorf("failed to prompt for key press")
			}

			orchestrator := New(nil, nil, nil, nil, fakeShell, fakeUserInput)
			err := orchestrator.RunInstruction(inst1, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected instruction run to fail")
			assert.Equal(t, "failed to prompt for key press", err.GoError().Error())
			assert.Equal(t, 1, execScriptCallCount)
			assert.Equal(t, 1, pressAnyKeyCallCount)
		})
	})
}

var FailedToClearScreenAfterInstructionRun = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeShell := shell.CreateFakeShell()
			execScriptCallCount := 0
			fakeShell.ExecuteScriptWithOutputToFileMock = func(dir string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execScriptCallCount++
				assert.Equal(t, relativeScriptPath, inst1.File)
				return nil
			}

			clearScreenCallCount := 0
			fakeShell.ClearScreenMock = func() error {
				clearScreenCallCount++
				return fmt.Errorf("failed to clear screen")
			}

			fakeUserInput := input.CreateFakeUserInput()
			pressAnyKeyCallCount := 0
			fakeUserInput.PressAnyKeyToContinueMock = func() error {
				pressAnyKeyCallCount++
				return nil
			}

			orchestrator := New(nil, nil, nil, nil, fakeShell, fakeUserInput)
			err := orchestrator.RunInstruction(inst1, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected instruction run to fail")
			assert.Equal(t, "failed to clear screen", err.GoError().Error())
			assert.Equal(t, 1, execScriptCallCount)
			assert.Equal(t, 1, pressAnyKeyCallCount)
			assert.Equal(t, 1, clearScreenCallCount)
		})
	})
}
