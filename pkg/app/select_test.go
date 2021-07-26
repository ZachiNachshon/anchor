package app

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/banner"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SelectActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail due to missing orchestrator from registry",
			Func: FailDueToMissingOrchestratorFromRegistry,
		},
		{
			Name: "fail due to missing banner from registry",
			Func: FailDueToMissingBannerFromRegistry,
		},
		{
			Name: "start application selection successfully",
			Func: StartApplicationSelectionSuccessfully,
		},
		{
			Name: "fail to run application selection flow",
			Func: FailToRunApplicationSelectionFlow,
		},
		{
			Name: "cancel the application selection flow successfully",
			Func: CancelApplicationSelectionFlowSuccessfully,
		},
		{
			Name: "run instruction selection after application selected successfully",
			Func: RunInstructionSelectionAfterApplicationSelectedSuccessfully,
		},
		{
			Name: "run app selection again if instructions are missing",
			Func: RunAppSelectionAgainIfInstructionsAreMissing,
		},
		{
			Name: "run app selection again when selecting instructions back option",
			Func: RunAppSelectionAgainWhenSelectingInstructionsBackOption,
		},
		{
			Name: "fail to execute instruction upon selection",
			Func: FailToExecuteInstructionUponSelection,
		},
		{
			Name: "return from instruction selection using the back option",
			Func: ReturnFromInstructionSelectionUsingTheBackOption,
		},
		{
			Name: "fail to run instruction action due to failure to ask user",
			Func: FailToRunInstructionActionDueToFailureToAskUser,
		},
		{
			Name: "fail to run instruction action due to failure to ask user",
			Func: FailToRunInstructionActionDueToFailureToAskUser,
		},
		{
			Name: "not run instruction action due to user option not to",
			Func: NotRunInstructionActionDueToUserOptingNotTo,
		},
		{
			Name: "run instruction action successfully",
			Func: RunInstructionActionSuccessfully,
		},
		{
			Name: "execute action after it was selected successfully",
			Func: ExecuteActionAfterItWasSelectedSuccessfully,
		},
		{
			Name: "do not add workflow option when instructions missing workflows",
			Func: DoNotAddWorkflowOptionWhenInstructionsMissingWorkflows,
		},
		{
			Name: "add back and workflow options to actions prompt selector",
			Func: AddBackAndWorkflowOptionsToActionsPromptSelector,
		},
		{
			Name: "add back option to workflows prompt selector",
			Func: AddBackOptionToWorkflowPromptSelector,
		},
	}
	harness.RunTests(t, tests)
}

var FailDueToMissingOrchestratorFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			err := StartApplicationSelectionFlow(ctx)
			assert.NotNil(t, err, "expected to fail app install flow")
			assert.Contains(t, err.Error(), "orchestrator")
		})
	})
}

var FailDueToMissingBannerFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeOrchestrator := orchestrator.CreateFakeOrchestrator()
			orchestrator.ToRegistry(ctx.Registry(), fakeOrchestrator)

			fakeShell := shell.CreateFakeShell()
			shell.ToRegistry(ctx.Registry(), fakeShell)

			fakeUserInput := input.CreateFakeUserInput()
			input.ToRegistry(ctx.Registry(), fakeUserInput)

			err := StartApplicationSelectionFlow(ctx)
			assert.NotNil(t, err, "expected to fail app install flow")
			assert.Contains(t, err.Error(), "banner")
		})
	})
}

var StartApplicationSelectionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeOrchestrator := orchestrator.CreateFakeOrchestrator()
			appSelectCallCount := 0
			fakeOrchestrator.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				appSelectCallCount++
				// Use programmatic keyboard interrupt to stop the selection flow
				return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
			}
			orchestrator.ToRegistry(ctx.Registry(), fakeOrchestrator)

			fakeShell := shell.CreateFakeShell()
			shell.ToRegistry(ctx.Registry(), fakeShell)

			fakeUserInput := input.CreateFakeUserInput()
			input.ToRegistry(ctx.Registry(), fakeUserInput)

			fakeBanner := banner.CreateFakeBanner()
			fakeBanner.PrintAnchorBannerMock = func() {}
			banner.ToRegistry(ctx.Registry(), fakeBanner)

			err := StartApplicationSelectionFlow(ctx)
			assert.NotNil(t, err, "expected graceful failure due to keyboard interrupt")
			assert.Equal(t, 1, appSelectCallCount)
		})
	})
}

var FailToRunApplicationSelectionFlow = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			selectionCallCount := 0
			o := orchestrator.CreateFakeOrchestrator()
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				selectionCallCount++
				return nil, errors.New(fmt.Errorf("failed to select an app"))
			}
			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected to fail app selection")
			assert.Equal(t, "failed to select an app", err.GoError().Error())
			assert.Equal(t, 1, selectionCallCount)
		})
	})
}

var CancelApplicationSelectionFlowSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			selectionCallCount := 0
			o := orchestrator.CreateFakeOrchestrator()
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				selectionCallCount++
				return &models.ApplicationInfo{
					Name: prompter.CancelActionName,
				}, nil
			}
			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.Nil(t, err, "expected selection to succeed")
			assert.Equal(t, 1, selectionCallCount)
		})
	})
}

var RunInstructionSelectionAfterApplicationSelectedSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			instRootTestData := stubs.GenerateInstructionsTestData()

			selectedApp := stubs.GetAppByName(apps, stubs.App1Name)
			appSelectionCallCount := 0
			o := orchestrator.CreateFakeOrchestrator()
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				appSelectionCallCount++
				return selectedApp, nil
			}

			extractInstCallCount := 0
			o.ExtractInstructionsMock = func(app *models.ApplicationInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractInstCallCount++
				return instRootTestData, nil
			}

			instructionSelectCallCount := 0
			o.OrchestrateInstructionActionSelectionMock = func(app *models.ApplicationInfo, actions []*models.Action) (*models.Action, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, selectedApp.Name)
				return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
			}

			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 1, appSelectionCallCount)
			assert.Equal(t, 1, extractInstCallCount)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var RunAppSelectionAgainIfInstructionsAreMissing = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			o := orchestrator.CreateFakeOrchestrator()

			appSelectCallCount := 0
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				appSelectCallCount++
				if appSelectCallCount == 1 {
					return app1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractInstCallCount := 0
			o.ExtractInstructionsMock = func(app *models.ApplicationInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractInstCallCount++
				return instRootTestData, nil
			}

			instructionSelectCallCount := 0
			o.OrchestrateInstructionActionSelectionMock = func(app *models.ApplicationInfo, actions []*models.Action) (*models.Action, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return nil, errors.NewInstructionMissingError(fmt.Errorf("missing instructions"))
			}

			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, appSelectCallCount)
			assert.Equal(t, 1, extractInstCallCount)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var RunAppSelectionAgainWhenSelectingInstructionsBackOption = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			o := orchestrator.CreateFakeOrchestrator()

			appSelectCallCount := 0
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				appSelectCallCount++
				if appSelectCallCount == 1 {
					return app1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractInstCallCount := 0
			o.ExtractInstructionsMock = func(app *models.ApplicationInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractInstCallCount++
				return instRootTestData, nil
			}

			instructionSelectCallCount := 0
			o.OrchestrateInstructionActionSelectionMock = func(app *models.ApplicationInfo, actions []*models.Action) (*models.Action, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return &models.Action{
					Id: prompter.BackActionName,
				}, nil
			}

			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, appSelectCallCount)
			assert.Equal(t, 1, extractInstCallCount)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var FailToExecuteInstructionUponSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			o := orchestrator.CreateFakeOrchestrator()

			instructionSelectCallCount := 0
			o.OrchestrateInstructionActionSelectionMock = func(app *models.ApplicationInfo, actions []*models.Action) (*models.Action, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return nil, errors.New(fmt.Errorf("failed to select instruction item"))
			}

			item, promptError := runInstructionActionSelectionFlow(o, app1, instRootTestData)
			assert.Nil(t, item, "expected to receive an empty input")
			assert.NotNil(t, promptError, "expected instruction selection to fail")
			assert.Equal(t, "failed to select instruction item", promptError.GoError().Error())
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var ReturnFromInstructionSelectionUsingTheBackOption = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			o := orchestrator.CreateFakeOrchestrator()

			instructionSelectCallCount := 0
			o.OrchestrateInstructionActionSelectionMock = func(app *models.ApplicationInfo, actions []*models.Action) (*models.Action, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return &models.Action{
					Id: prompter.BackActionName,
				}, nil
			}

			item, promptError := runInstructionActionSelectionFlow(o, app1, instRootTestData)
			assert.NotNil(t, item, "expected to receive an input")
			assert.Nil(t, promptError, "expected instruction selection not to fail")
			assert.Equal(t, prompter.BackActionName, item.Id)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var FailToRunInstructionActionDueToFailureToAskUser = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)
			o := orchestrator.CreateFakeOrchestrator()

			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionActionMock = func(item *models.Action) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				// Do not run the instruction
				return false, errors.New(fmt.Errorf("failed to ask user to press any key"))
			}

			item, promptError := runInstructionActionExecutionFlow(o, action)
			assert.Nil(t, item, "expected to receive an empty input")
			assert.NotNil(t, promptError, "expected instruction execution to fail")
			assert.Equal(t, "failed to ask user to press any key", promptError.GoError().Error())
			assert.Equal(t, 1, askBeforeRunCallCount)
		})
	})
}

var NotRunInstructionActionDueToUserOptingNotTo = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)
			o := orchestrator.CreateFakeOrchestrator()

			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionActionMock = func(item *models.Action) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				return false, nil
			}

			item, promptError := runInstructionActionExecutionFlow(o, action)
			assert.Nil(t, promptError, "expected instruction execution not to fail")
			assert.NotNil(t, item, "expected to receive a valid input")
			assert.Equal(t, action.Id, item.Id)
			assert.Equal(t, 1, askBeforeRunCallCount)
		})
	})
}

var RunInstructionActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)
			o := orchestrator.CreateFakeOrchestrator()

			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionActionMock = func(item *models.Action) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			o.RunInstructionActionMock = func(action *models.Action) *errors.PromptError {
				runInstructionCallCount++
				return nil
			}

			wrapAfterActionCallCount := 0
			o.WrapAfterActionRunMock = func() *errors.PromptError {
				wrapAfterActionCallCount++
				return nil
			}

			item, promptError := runInstructionActionExecutionFlow(o, action)
			assert.Nil(t, promptError, "expected instruction execution not to fail")
			assert.NotNil(t, item, "expected to receive a valid input")
			assert.Equal(t, action.Id, item.Id)
			assert.Equal(t, 1, askBeforeRunCallCount)
			assert.Equal(t, 1, runInstructionCallCount)
			assert.Equal(t, 1, wrapAfterActionCallCount)
		})
	})
}

var ExecuteActionAfterItWasSelectedSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			action := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.App1Action1Id)
			o := orchestrator.CreateFakeOrchestrator()

			instructionSelectCallCount := 0
			o.OrchestrateInstructionActionSelectionMock = func(app *models.ApplicationInfo, actions []*models.Action) (*models.Action, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return action, nil
			}

			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionActionMock = func(item *models.Action) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				// Use programmatic keyboard interrupt to stop the selection flow
				return false, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
			}

			item, promptError := runInstructionActionSelectionFlow(o, app1, instRootTestData)
			assert.Nil(t, item, "expected to receive an empty input")
			assert.NotNil(t, promptError, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", promptError.GoError().Error())
			assert.Equal(t, 1, instructionSelectCallCount)
			assert.Equal(t, 1, askBeforeRunCallCount)
		})
	})
}
var DoNotAddWorkflowOptionWhenInstructionsMissingWorkflows = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	instRootTestData.Instructions.Workflows = nil
	appendInstructionActionsCustomOptions(instRootTestData.Instructions)
	assert.EqualValues(t, prompter.BackActionName, instRootTestData.Instructions.Actions[0].Id)
	assert.NotEqual(t, prompter.WorkflowsActionName, instRootTestData.Instructions.Actions[1].Id)
}

var AddBackAndWorkflowOptionsToActionsPromptSelector = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	appendInstructionActionsCustomOptions(instRootTestData.Instructions)
	assert.EqualValues(t, prompter.BackActionName, instRootTestData.Instructions.Actions[0].Id)
	assert.EqualValues(t, prompter.WorkflowsActionName, instRootTestData.Instructions.Actions[1].Id)
}

var AddBackOptionToWorkflowPromptSelector = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	appendInstructionWorkflowsCustomOptions(instRootTestData.Instructions)
	assert.EqualValues(t, prompter.BackActionName, instRootTestData.Instructions.Workflows[0].Id)
}
