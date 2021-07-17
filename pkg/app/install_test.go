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

func Test_InstallActionShould(t *testing.T) {
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
			Name: "start application install flow successfully",
			Func: StartApplicationInstallFlowSuccessfully,
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
			Name: "run instruction flow after application selected successfully",
			Func: RunInstructionFlowAfterApplicationSelectedSuccessfully,
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
			Name: "fail instruction selection flow due to error",
			Func: FailInstructionSelectionDueToError,
		},
		{
			Name: "return from instruction selection using the back option",
			Func: ReturnFromInstructionSelectionUsingTheBackOption,
		},
		{
			Name: "fail to execute instruction it was after selected",
			Func: FailToExecuteInstructionAfterItWasSelected,
		},
		{
			Name: "fail to run instruction execution flow",
			Func: FailToRunInstructionExecutionFlow,
		},
		{
			Name: "run instruction execution flow with ",
			Func: FailToRunInstructionExecutionFlow,
		},
		{
			Name: "run instruction execution flow without running instruction",
			Func: RunInstructionExecutionFlowWithoutRunningInstruction,
		},
		{
			Name: "run instruction execution flow and run instruction successfully",
			Func: RunInstructionExecutionFlowAndRunInstructionSuccessfully,
		},
		{
			Name: "run execution flow after instruction selection successfully",
			Func: RunExecutionFlowAfterInstructionSelectionSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var FailDueToMissingOrchestratorFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			err := StartApplicationInstallFlow(ctx)
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

			err := StartApplicationInstallFlow(ctx)
			assert.NotNil(t, err, "expected to fail app install flow")
			assert.Contains(t, err.Error(), "banner")
		})
	})
}

var StartApplicationInstallFlowSuccessfully = func(t *testing.T) {
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

			err := StartApplicationInstallFlow(ctx)
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
					Name: prompter.CancelButtonName,
				}, nil
			}
			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.Nil(t, err, "expected selection to succeed")
			assert.Equal(t, 1, selectionCallCount)
		})
	})
}

var RunInstructionFlowAfterApplicationSelectedSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			selectedApp := stubs.GetAppByName(apps, stubs.App1Name)
			appSelectCallCount := 0
			o := orchestrator.CreateFakeOrchestrator()
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				appSelectCallCount++
				return selectedApp, nil
			}
			instructionSelectCallCount := 0
			o.OrchestrateInstructionSelectionMock = func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, selectedApp.Name)
				return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
			}
			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 1, appSelectCallCount)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var RunAppSelectionAgainIfInstructionsAreMissing = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			appSelectCallCount := 0
			o := orchestrator.CreateFakeOrchestrator()
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				appSelectCallCount++
				if appSelectCallCount == 1 {
					return app1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}
			instructionSelectCallCount := 0
			o.OrchestrateInstructionSelectionMock = func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return nil, errors.NewInstructionMissingError(fmt.Errorf("missing instructions"))
			}
			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, appSelectCallCount)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var RunAppSelectionAgainWhenSelectingInstructionsBackOption = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			appSelectCallCount := 0
			o := orchestrator.CreateFakeOrchestrator()
			o.OrchestrateApplicationSelectionMock = func() (*models.ApplicationInfo, *errors.PromptError) {
				appSelectCallCount++
				if appSelectCallCount == 1 {
					return app1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}
			instructionSelectCallCount := 0
			o.OrchestrateInstructionSelectionMock = func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return &models.InstructionItem{
					Id: prompter.BackButtonName,
				}, nil
			}
			err := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, appSelectCallCount)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var FailInstructionSelectionDueToError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			o := orchestrator.CreateFakeOrchestrator()
			instructionSelectCallCount := 0
			o.OrchestrateInstructionSelectionMock = func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return nil, errors.New(fmt.Errorf("failed to select instruction item"))
			}
			item, promptError := runInstructionSelectionFlow(app1, o, ctx.AnchorFilesPath())
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
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			o := orchestrator.CreateFakeOrchestrator()
			instructionSelectCallCount := 0
			o.OrchestrateInstructionSelectionMock = func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return &models.InstructionItem{
					Id: prompter.BackButtonName,
				}, nil
			}
			item, promptError := runInstructionSelectionFlow(app1, o, ctx.AnchorFilesPath())
			assert.NotNil(t, item, "expected to receive an input")
			assert.Nil(t, promptError, "expected instruction selection not to fail")
			assert.Equal(t, prompter.BackButtonName, item.Id)
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var FailToExecuteInstructionAfterItWasSelected = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			instructions := stubs.GenerateInstructionsTestData()
			instructions1 := stubs.GetInstructionItemById(instructions, stubs.App1InstructionsItem1Id)
			o := orchestrator.CreateFakeOrchestrator()
			instructionSelectCallCount := 0
			o.OrchestrateInstructionSelectionMock = func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return instructions1, nil
			}
			o.AskBeforeRunningInstructionMock = func(item *models.InstructionItem) (bool, *errors.PromptError) {
				return false, errors.New(fmt.Errorf("failed to ask before running instruction"))
			}
			item, promptError := runInstructionSelectionFlow(app1, o, ctx.AnchorFilesPath())
			assert.Nil(t, item, "expected to receive an empty input")
			assert.NotNil(t, promptError, "expected instruction selection to fail")
			assert.Equal(t, "failed to ask before running instruction", promptError.GoError().Error())
			assert.Equal(t, 1, instructionSelectCallCount)
		})
	})
}

var FailToRunInstructionExecutionFlow = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instructions := stubs.GenerateInstructionsTestData()
			instructions1 := stubs.GetInstructionItemById(instructions, stubs.App1InstructionsItem1Id)
			o := orchestrator.CreateFakeOrchestrator()
			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionMock = func(item *models.InstructionItem) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				// Do not run the instruction
				return false, errors.New(fmt.Errorf("failed to ask user to press any key"))
			}

			item, promptError := runInstructionExecutionFlow(instructions1, o, ctx.AnchorFilesPath())
			assert.Nil(t, item, "expected to receive an empty input")
			assert.NotNil(t, promptError, "expected instruction execution to fail")
			assert.Equal(t, "failed to ask user to press any key", promptError.GoError().Error())
			assert.Equal(t, 1, askBeforeRunCallCount)
		})
	})
}

var RunInstructionExecutionFlowWithoutRunningInstruction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instructions := stubs.GenerateInstructionsTestData()
			instructions1 := stubs.GetInstructionItemById(instructions, stubs.App1InstructionsItem1Id)
			o := orchestrator.CreateFakeOrchestrator()
			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionMock = func(item *models.InstructionItem) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				// Do not run the instruction
				return false, nil
			}

			item, promptError := runInstructionExecutionFlow(instructions1, o, ctx.AnchorFilesPath())
			assert.Nil(t, promptError, "expected instruction execution not to fail")
			assert.NotNil(t, item, "expected to receive a valid input")
			assert.Equal(t, instructions1.Id, item.Id)
			assert.Equal(t, 1, askBeforeRunCallCount)
		})
	})
}

var RunInstructionExecutionFlowAndRunInstructionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instructions := stubs.GenerateInstructionsTestData()
			instructions1 := stubs.GetInstructionItemById(instructions, stubs.App1InstructionsItem1Id)
			o := orchestrator.CreateFakeOrchestrator()
			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionMock = func(item *models.InstructionItem) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				// Do not run the instruction
				return true, nil
			}
			runInstructionCallCount := 0
			o.RunInstructionMock = func(item *models.InstructionItem, repoPath string) *errors.PromptError {
				runInstructionCallCount++
				return nil
			}
			item, promptError := runInstructionExecutionFlow(instructions1, o, ctx.AnchorFilesPath())
			assert.Nil(t, promptError, "expected instruction execution not to fail")
			assert.NotNil(t, item, "expected to receive a valid input")
			assert.Equal(t, instructions1.Id, item.Id)
			assert.Equal(t, 1, askBeforeRunCallCount)
			assert.Equal(t, 1, runInstructionCallCount)
		})
	})
}

var RunExecutionFlowAfterInstructionSelectionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			apps := stubs.GenerateApplicationTestData()
			app1 := stubs.GetAppByName(apps, stubs.App1Name)
			instructions := stubs.GenerateInstructionsTestData()
			instructions1 := stubs.GetInstructionItemById(instructions, stubs.App1InstructionsItem1Id)
			o := orchestrator.CreateFakeOrchestrator()
			instructionSelectCallCount := 0
			o.OrchestrateInstructionSelectionMock = func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
				instructionSelectCallCount++
				assert.Equal(t, app.Name, app1.Name)
				return instructions1, nil
			}
			askBeforeRunCallCount := 0
			o.AskBeforeRunningInstructionMock = func(item *models.InstructionItem) (bool, *errors.PromptError) {
				askBeforeRunCallCount++
				// Use programmatic keyboard interrupt to stop the selection flow
				return false, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
			}
			item, promptError := runInstructionSelectionFlow(app1, o, ctx.AnchorFilesPath())
			assert.Nil(t, item, "expected to receive an empty input")
			assert.NotNil(t, promptError, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", promptError.GoError().Error())
			assert.Equal(t, 1, instructionSelectCallCount)
			assert.Equal(t, 1, askBeforeRunCallCount)
		})
	})
}
