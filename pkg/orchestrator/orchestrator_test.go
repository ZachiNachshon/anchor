package orchestrator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/models"

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
			Name: "fail to prompt for instructions actions",
			Func: FailToPromptForInstructionsActions,
		},
		{
			Name: "select instruction action successfully",
			Func: SelectInstructionActionSuccessfully,
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
			Name: "run instruction action by script file successfully",
			Func: RunInstructionActionByScriptFileSuccessfully,
		},
		{
			Name: "failed to run script file for instruction action",
			Func: FailedToRunScriptFileForInstructionAction,
		},
		{
			Name: "failed to run action due to script and script file mutual exclusivity",
			Func: FailedToRunActionDueToScriptAndScriptFileMutualExclusivity,
		},
		{
			Name: "failed to clear screen after instruction run",
			Func: FailedToRunActionDueToMissingScriptAndScriptFile,
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
			instItem, err := orchestrator.ExtractInstructions(app1, ctx.AnchorFilesPath())
			assert.Nil(t, instItem, "expected instruction to be empty")
			assert.NotNil(t, err, "expected instruction selection to fail")
			assert.Equal(t, 1, extractorCallCount)
		})
	})
}

var FailToPromptForInstructionsActions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			instTestData := stubs.GenerateInstructionsTestData()
			actions := instTestData.Instructions.Actions

			fakePrompter := prompter.CreateFakePrompter()
			instPromptCallCount := 0
			fakePrompter.PromptInstructionActionsMock = func(appName string, actions []*models.Action) (*models.Action, error) {
				instPromptCallCount++
				return nil, fmt.Errorf("failed to prompt for instructions")
			}

			orchestrator := New(fakePrompter, nil, nil, nil, nil, nil)
			instItem, err := orchestrator.OrchestrateInstructionActionSelection(app1, actions)
			assert.Nil(t, instItem, "expected instruction to be empty")
			assert.NotNil(t, err, "expected instruction selection to fail")
			assert.Equal(t, 1, instPromptCallCount)
		})
	})
}

var SelectInstructionActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			actions := instRootTestData.Instructions.Actions
			inst1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakePrompter := prompter.CreateFakePrompter()
			instPromptCallCount := 0
			fakePrompter.PromptInstructionActionsMock = func(appName string, actions []*models.Action) (*models.Action, error) {
				instPromptCallCount++
				return inst1, nil
			}

			orchestrator := New(fakePrompter, nil, nil, nil, nil, nil)
			instItem, err := orchestrator.OrchestrateInstructionActionSelection(app1, actions)
			assert.NotNil(t, instItem, "expected instruction not to be empty")
			assert.Nil(t, err, "expected instruction selection not to fail")
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
			shouldRun, err := orchestrator.AskBeforeRunningInstructionAction(inst1)
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
			shouldRun, err := orchestrator.AskBeforeRunningInstructionAction(inst1)
			assert.Equal(t, true, shouldRun)
			assert.Nil(t, err, "expected instruction selection to succeed")
			assert.Equal(t, 1, userInputCallCount)
		})
	})
}

var RunInstructionActionByScriptFileSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeShell := shell.CreateFakeShell()
			execScriptCallCount := 0
			fakeShell.ExecuteScriptFileWithOutputToFileMock = func(dir string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execScriptCallCount++
				assert.Equal(t, relativeScriptPath, action.ScriptFile)
				return nil
			}

			orchestrator := New(nil, nil, nil, nil, fakeShell, nil)
			err := orchestrator.RunInstructionAction(action)
			assert.Nil(t, err, "expected instruction run to succeed")
			assert.Equal(t, 1, execScriptCallCount)
		})
	})
}

var FailedToRunScriptFileForInstructionAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)

			fakeShell := shell.CreateFakeShell()
			execScriptCallCount := 0
			fakeShell.ExecuteScriptFileWithOutputToFileMock = func(dir string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execScriptCallCount++
				assert.Equal(t, relativeScriptPath, action.ScriptFile)
				return fmt.Errorf("failed to execute script")
			}

			orchestrator := New(nil, nil, nil, nil, fakeShell, nil)
			err := orchestrator.RunInstructionAction(action)
			assert.NotNil(t, err, "expected instruction run to fail")
			assert.Equal(t, "failed to execute script", err.GoError().Error())
			assert.Equal(t, 1, execScriptCallCount)
		})
	})
}

var FailedToRunActionDueToScriptAndScriptFileMutualExclusivity = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)
			action.Script = "echo hello world"

			orchestrator := New(nil, nil, nil, nil, nil, nil)
			err := orchestrator.RunInstructionAction(action)
			assert.NotNil(t, err, "expected instruction action run to fail")
			assert.Contains(t, err.GoError().Error(), "script / scriptFile are mutual exclusive")
		})
	})
}

var FailedToRunActionDueToMissingScriptAndScriptFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)
			action.Script = ""
			action.ScriptFile = ""

			orchestrator := New(nil, nil, nil, nil, nil, nil)
			err := orchestrator.RunInstructionAction(action)
			assert.NotNil(t, err, "expected instruction action run to fail")
			assert.Contains(t, err.GoError().Error(), "missing script or scriptFile")
		})
	})
}
