package _select

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/runner"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SelectActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
		{
			Name: "start the selection flow when run is called",
			Func: StartTheSelectionFlowWhenRunIsCalled,
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
			Name: "fail preparing registry items",
			Func: FailPreparingRegistryItems,
		},
		{
			Name: "fail to run selection for dynamic command",
			Func: FailToRunSelectionForDynamicCommand,
		},
		{
			Name: "command folder item selection: cancel selection successfully",
			Func: CommandFolderItemSelectionCancelSelectionSuccessfully,
		},
		{
			Name: "command folder item selection: fail to prompt",
			Func: CommandFolderItemSelectionFailToPrompt,
		},
		{
			Name: "command folder item selection: fail to extract instructions",
			Func: CommandFolderItemSelectionFailToExtractInstructions,
		},
		{
			Name: "command folder item selection: go back from instruction action selection",
			Func: CommandFolderItemSelectionGoBackFromInstructionActionSelection,
		},
		{
			Name: "command folder item selection: instruction action selection missing error",
			Func: CommandFolderItemSelectionInstructionActionSelectionMissingError,
		},
		{
			Name: "command folder item selection: instruction action selection general error",
			Func: CommandFolderItemSelectionInstructionActionSelectionGeneralError,
		},
		{
			Name: "command folder item selection: complete successfully",
			Func: CommandFolderItemSelectionCompleteSuccessfully,
		},
		{
			Name: "command folder item prompt: fail to prompt",
			Func: CommandFolderItemPromptFailToPrompt,
		},
		{
			Name: "command folder item prompt: prompt successfully",
			Func: CommandFolderItemPromptPromptSuccessfully,
		},
		{
			Name: "exec wrap: fail to press any key",
			Func: ExecWrapFailToPressAnyKey,
		},
		{
			Name: "exec wrap: fail to clear screen",
			Func: ExecWrapFailToClearScreen,
		},
		{
			Name: "exec wrap: warp up successfully",
			Func: ExecWrapWrapUpSuccessfully,
		},
		{
			Name: "instruction action selection: fail to prompt",
			Func: InstructionActionSelectionFailToPrompt,
		},
		{
			Name: "instruction action selection: go back",
			Func: InstructionActionSelectionGoBack,
		},
		{
			Name: "instruction action selection: fail workflow selection",
			Func: InstructionActionSelectionFailWorkflowSelection,
		},
		{
			Name: "instruction action selection: succeed workflow selection",
			Func: InstructionActionSelectionSucceedWorkflowSelection,
		},
		{
			Name: "instruction action selection: fail action execution",
			Func: InstructionActionSelectionFailActionExecution,
		},
		{
			Name: "instruction action selection: succeed action execution",
			Func: InstructionActionSelectionSucceedActionExecution,
		},
		{
			Name: "instruction action prompt: fail to prompt",
			Func: InstructionActionPromptFailToPrompt,
		},
		{
			Name: "instruction action prompt: prompt successfully",
			Func: InstructionActionPromptPromptSuccessfully,
		},
		{
			Name: "instruction action exec flow: fail to ask before running",
			Func: InstructionActionExecFailToAskBeforeRunning,
		},
		{
			Name: "instruction action exec flow: fail to run",
			Func: InstructionActionExecFailToRun,
		},
		{
			Name: "instruction action exec flow: fail to wrap after run",
			Func: InstructionActionExecFailToWrapAfterRun,
		},
		{
			Name: "instruction action exec flow: run successfully",
			Func: InstructionActionExecRunSuccessfully,
		},
		{
			Name: "instruction action exec flow: fail to ask yes/no question",
			Func: InstructionActionExecFailToAskYesNoQuestion,
		},
		{
			Name: "instruction action exec: ask yes/no question for application context",
			Func: InstructionActionExecAskYesNoQuestionForApplicationContext,
		},
		{
			Name: "instruction action exec: ask yes/no question for kubernetes context",
			Func: InstructionActionExecAskYesNoQuestionForKubernetesContext,
		},
		{
			Name: "instruction workflow selection: fail to prompt",
			Func: InstructionWorkflowSelectionFailToPrompt,
		},
		{
			Name: "instruction workflow selection: go back",
			Func: InstructionWorkflowSelectionGoBack,
		},
		{
			Name: "instruction workflow selection: fail workflow execution",
			Func: InstructionWorkflowSelectionFailWorkflowExecution,
		},
		{
			Name: "instruction workflow selection: succeed workflow execution",
			Func: InstructionWorkflowSelectionSucceedWorkflowExecution,
		},
		{
			Name: "instruction workflow exec: fail to ask yes/no question",
			Func: InstructionWorkflowExecFailToAskYesNoQuestion,
		},
		{
			Name: "instruction workflow exec: ask yes/no question for application context",
			Func: InstructionWorkflowExecAskYesNoQuestionForApplicationContext,
		},
		{
			Name: "instruction workflow exec: ask yes/no question for kubernetes context",
			Func: InstructionWorkflowExecAskYesNoQuestionForKubernetesContext,
		},
		{
			Name: "instruction workflow prompt: fail to prompt",
			Func: InstructionWorkflowPromptFailToPrompt,
		},
		{
			Name: "instruction workflow prompt: prompt successfully",
			Func: InstructionWorkflowPromptPromptSuccessfully,
		},
		{
			Name: "instruction workflow exec: fail to ask before running",
			Func: InstructionWorkflowExecFailToAskBeforeRunning,
		},
		{
			Name: "instruction workflow exec: tolerate run failures",
			Func: InstructionWorkflowExecTolerateRunFailures,
		},
		{
			Name: "instruction workflow exec: fail to wrap after run",
			Func: InstructionWorkflowExecFailToWrapAfterRun,
		},
		{
			Name: "instruction workflow exec: run successfully",
			Func: InstructionWorkflowExecRunSuccessfully,
		},
		{
			Name: "manage prompt error: missing inner go error",
			Func: ManagePromptErrorMissingInnerGoError,
		},
		{
			Name: "manage prompt error: mitigate interrupt error",
			Func: ManagePromptErrorMitigateInterruptError,
		},
		{
			Name: "manage prompt error: return inner go error",
			Func: ManagePromptErrorReturnInnerGoError,
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
		{
			Name: "do not add back option to workflows prompt selector if already exist",
			Func: DoNotAddBackOptionToWorkflowPromptSelectorIfAlreadyExist,
		},
		{
			Name: "extract args from script file",
			Func: ExtractArgsFromScriptFile,
		},
		{
			Name: "resolve instructions action context successfully",
			Func: ResolveInstructionActionContextSuccessfully,
		},
		{
			Name: "resolve instructions workflow context successfully",
			Func: ResolveInstructionWorkflowContextSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			runnerPrepareCallCount := 0
			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				runnerPrepareCallCount++
				return nil
			}
			fakeO := NewOrchestrator(fakeRunner, cmdFolderName)
			bannerCallCount := 0
			fakeO.bannerFunc = func(o *selectOrchestrator) {
				bannerCallCount++
			}
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *selectOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return nil
			}
			runCallCount := 0
			fakeO.runFunc = func(o *selectOrchestrator, ctx common.Context) *errors.PromptError {
				runCallCount++
				return nil
			}
			err := DynamicSelect(ctx, fakeO)
			assert.Nil(t, err, "expected not to fail command folder item status")
			assert.Equal(t, 1, runnerPrepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, bannerCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var StartTheSelectionFlowWhenRunIsCalled = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			startFolderItemSelectionCallCount := 0
			fakeO.startCommandItemsSelectionFlowFunc = func(o *selectOrchestrator, anchorfilesRepoPath string) *errors.PromptError {
				startFolderItemSelectionCallCount++
				return nil
			}
			err := fakeO.runFunc(fakeO, ctx)
			assert.Nil(t, err, "expected not to fail command folder item status")
			assert.Equal(t, 1, startFolderItemSelectionCallCount, "expected func to be called exactly once")
		})
	})
}

var PrintBanner = func(t *testing.T) {
	fakePrinter := printer.CreateFakePrinter()
	printBannerCallCount := 0
	fakePrinter.PrintAnchorBannerMock = func() {
		printBannerCallCount++
	}

	cmdFolderName := stubs.CommandFolder1Name
	fakeRunner := runner.NewOrchestrator(cmdFolderName)
	fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
	fakeO.prntr = fakePrinter
	fakeO.bannerFunc(fakeO)
	assert.Equal(t, 1, printBannerCallCount, "expected func to be called exactly once")
}

var PrepareRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()

		fakeLocator := locator.CreateFakeLocator()
		reg.Set(locator.Identifier, fakeLocator)

		fakePrompter := prompter.CreateFakePrompter()
		reg.Set(prompter.Identifier, fakePrompter)

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		fakeShell := shell.CreateFakeShell()
		reg.Set(shell.Identifier, fakeShell)

		fakeInput := input.CreateFakeUserInput()
		reg.Set(input.Identifier, fakeInput)

		cmdFolderName := stubs.CommandFolder1Name
		fakeRunner := runner.NewOrchestrator(cmdFolderName)
		fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
		err := fakeO.prepareFunc(fakeO, ctx)

		assert.Nil(t, err)
		assert.NotNil(t, fakeO.l)
		assert.NotNil(t, fakeO.prmpt)
		assert.NotNil(t, fakeO.prntr)
		assert.NotNil(t, fakeO.s)
		assert.NotNil(t, fakeO.in)
	})
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()

		fakeLocator := locator.CreateFakeLocator()
		fakePrompter := prompter.CreateFakePrompter()
		fakePrinter := printer.CreateFakePrinter()
		fakeShell := shell.CreateFakeShell()
		fakeInput := input.CreateFakeUserInput()

		cmdFolderName := stubs.CommandFolder1Name
		fakeRunner := runner.NewOrchestrator(cmdFolderName)
		fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)

		err := fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", locator.Identifier))
		reg.Set(locator.Identifier, fakeLocator)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", prompter.Identifier))
		reg.Set(prompter.Identifier, fakePrompter)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", printer.Identifier))
		reg.Set(printer.Identifier, fakePrinter)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", shell.Identifier))
		reg.Set(shell.Identifier, fakeShell)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", input.Identifier))
		reg.Set(input.Identifier, fakeInput)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
	})
}

var FailPreparingRegistryItems = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			prepareRunnerCallCount := 0
			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				prepareRunnerCallCount++
				return fmt.Errorf("failed to prepare runner")
			}
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			err := DynamicSelect(ctx, fakeO)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare runner", err.Error())
			assert.Equal(t, 1, prepareRunnerCallCount, "expected func to be called exactly once")

			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				return nil
			}
			prepareRegistryItemsCallCount := 0
			fakeO.prepareFunc = func(o *selectOrchestrator, ctx common.Context) error {
				prepareRegistryItemsCallCount++
				return fmt.Errorf("failed to prepare registry items")
			}
			err = DynamicSelect(ctx, fakeO)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare registry items", err.Error())
			assert.Equal(t, 1, prepareRegistryItemsCallCount, "expected func to be called exactly once")
		})
	})
}

var FailToRunSelectionForDynamicCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				return nil
			}
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			prepareRegistryItemsCallCount := 0
			fakeO.prepareFunc = func(o *selectOrchestrator, ctx common.Context) error {
				prepareRegistryItemsCallCount++
				return nil
			}
			bannerCallCount := 0
			fakeO.bannerFunc = func(o *selectOrchestrator) {
				bannerCallCount++
			}
			runCallCount := 0
			fakeO.runFunc = func(o *selectOrchestrator, ctx common.Context) *errors.PromptError {
				runCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to run selection"))
			}
			err := DynamicSelect(ctx, fakeO)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to run selection", err.Error())
			assert.Equal(t, 1, prepareRegistryItemsCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, bannerCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var CommandFolderItemSelectionCancelSelectionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptCommandItemsSelectionFunc = func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				return &models.CommandFolderItemInfo{
					Name: prompter.CancelActionName,
				}, nil
			}
			err := fakeO.startCommandItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.Nil(t, err, "expected selection to succeed")
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var CommandFolderItemSelectionFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptCommandItemsSelectionFunc = func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				return nil, errors.NewPromptError(fmt.Errorf("failed to prompt"))
			}
			err := fakeO.startCommandItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection prompt to fail")
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, "failed to prompt", err.GoError().Error())
		})
	})
}

var CommandFolderItemSelectionFailToExtractInstructions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptCommandItemsSelectionFunc = func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				if promptCallCount == 1 {
					return item1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractorCallCount := 0
			fakeRunner.ExtractInstructionsFunc = func(
				o *runner.ActionRunnerOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				assert.Equal(t, commandFolderItem.InstructionsPath, item1.InstructionsPath)
				return nil, errors.NewPromptError(fmt.Errorf("failed to extract instructions"))
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrintMissingInstructionsMock = func() {}
			fakeO.prntr = fakePrinter

			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				return nil
			}

			err := fakeO.startCommandItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
		})
	})
}

var CommandFolderItemSelectionGoBackFromInstructionActionSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptCommandItemsSelectionFunc = func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				if promptCallCount == 1 {
					return item1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractorCallCount := 0
			fakeRunner.ExtractInstructionsFunc = func(
				o *runner.ActionRunnerOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				return instRootTestData, nil
			}

			instActionSelectionCallCount := 0
			fakeO.startInstructionActionSelectionFlowFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {
				instActionSelectionCallCount++
				return &models.Action{
					Id: prompter.BackActionName,
				}, nil
			}

			err := fakeO.startCommandItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instActionSelectionCallCount)
		})
	})
}

var CommandFolderItemSelectionInstructionActionSelectionMissingError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptCommandItemsSelectionFunc = func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				if promptCallCount == 1 {
					return item1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractorCallCount := 0
			fakeRunner.ExtractInstructionsFunc = func(
				o *runner.ActionRunnerOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				return instRootTestData, nil
			}

			instActionSelectionCallCount := 0
			fakeO.startInstructionActionSelectionFlowFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {
				instActionSelectionCallCount++
				return nil, errors.NewInstructionMissingError(fmt.Errorf("missing instruction"))
			}

			err := fakeO.startCommandItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instActionSelectionCallCount)
		})
	})
}

var CommandFolderItemSelectionInstructionActionSelectionGeneralError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			fakeRunner := runner.NewOrchestrator("")
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			promptCallCount := 0
			fakeO.promptCommandItemsSelectionFunc = func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				return item1, nil
			}

			extractorCallCount := 0
			fakeRunner.ExtractInstructionsFunc = func(
				o *runner.ActionRunnerOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				return nil, nil
			}

			instActionSelectionCallCount := 0
			fakeO.startInstructionActionSelectionFlowFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {
				instActionSelectionCallCount++
				return nil, errors.NewPromptError(fmt.Errorf("failed to start instruction action selection flow"))
			}

			err := fakeO.startCommandItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to fail")
			assert.Equal(t, "failed to start instruction action selection flow", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instActionSelectionCallCount)
		})
	})
}

var CommandFolderItemSelectionCompleteSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptCommandItemsSelectionFunc = func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				return item1, nil
			}

			extractorCallCount := 0
			fakeRunner.ExtractInstructionsFunc = func(
				o *runner.ActionRunnerOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				assert.Equal(t, item1, commandFolderItem)
				return instRootTestData, nil
			}

			instActionSelectionCallCount := 0
			fakeO.startInstructionActionSelectionFlowFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {
				instActionSelectionCallCount++
				assert.Equal(t, item1, commandFolderItem)
				return models.GetInstructionActionById(instructionRoot.Instructions, stubs.CommandFolder1Item1Action1Id), nil
			}

			err := fakeO.startCommandItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.Nil(t, err, "expected selection to complete successfully")
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instActionSelectionCallCount)
		})
	})
}

var CommandFolderItemPromptFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()

			fakeLocator := locator.CreateFakeLocator()
			locateCommandFolderItemsCallCount := 0
			fakeLocator.CommandFolderItemsMock = func(commandFolderName string) []*models.CommandFolderItemInfo {
				locateCommandFolderItemsCallCount++
				return items
			}

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptCommandFolderItemSelectionMock = func(appsArr []*models.CommandFolderItemInfo) (*models.CommandFolderItemInfo, error) {
				promptCallCount++
				return nil, fmt.Errorf("failed to prompt for items")
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.l = fakeLocator
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptCommandItemsSelectionFunc(fakeO)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected prompt to fail")
			assert.Equal(t, "failed to prompt for items", err.GoError().Error())
			assert.Equal(t, 1, locateCommandFolderItemsCallCount)
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var CommandFolderItemPromptPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)

			fakeLocator := locator.CreateFakeLocator()
			locateCommandFolderItemsCallCount := 0
			fakeLocator.CommandFolderItemsMock = func(commandFolderName string) []*models.CommandFolderItemInfo {
				locateCommandFolderItemsCallCount++
				return items
			}

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptCommandFolderItemSelectionMock = func(appsArr []*models.CommandFolderItemInfo) (*models.CommandFolderItemInfo, error) {
				promptCallCount++
				return item1, nil
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.l = fakeLocator
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptCommandItemsSelectionFunc(fakeO)
			assert.Nil(t, err, "expected prompt to succeed")
			assert.NotNil(t, result, "expected prompt response")
			assert.Equal(t, 1, locateCommandFolderItemsCallCount)
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, item1, result)
		})
	})
}

var ExecWrapFailToPressAnyKey = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeInput := input.CreateFakeUserInput()
			pressAnyKeyCallCount := 0
			fakeInput.PressAnyKeyToContinueMock = func() error {
				pressAnyKeyCallCount++
				return fmt.Errorf("failed to press any key")
			}
			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrintEmptyLinesMock = func(count int) {}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.in = fakeInput
			fakeO.prntr = fakePrinter

			err := fakeO.wrapAfterExecutionFunc(fakeO)
			assert.NotNil(t, err, "expected wrap up to fail")
			assert.Equal(t, "failed to press any key", err.GoError().Error())
			assert.Equal(t, 1, pressAnyKeyCallCount)
		})
	})
}

var ExecWrapFailToClearScreen = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeInput := input.CreateFakeUserInput()
			pressAnyKeyCallCount := 0
			fakeInput.PressAnyKeyToContinueMock = func() error {
				pressAnyKeyCallCount++
				return nil
			}
			fakeShell := shell.CreateFakeShell()
			clearScreenCallCount := 0
			fakeShell.ClearScreenMock = func() error {
				clearScreenCallCount++
				return fmt.Errorf("failed to clean screen")
			}
			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrintEmptyLinesMock = func(count int) {}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.in = fakeInput
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := fakeO.wrapAfterExecutionFunc(fakeO)
			assert.NotNil(t, err, "expected clear screen to fail")
			assert.Equal(t, "failed to clean screen", err.GoError().Error())
			assert.Equal(t, 1, pressAnyKeyCallCount)
			assert.Equal(t, 1, clearScreenCallCount)
		})
	})
}

var ExecWrapWrapUpSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeInput := input.CreateFakeUserInput()
			pressAnyKeyCallCount := 0
			fakeInput.PressAnyKeyToContinueMock = func() error {
				pressAnyKeyCallCount++
				return nil
			}
			fakeShell := shell.CreateFakeShell()
			clearScreenCallCount := 0
			fakeShell.ClearScreenMock = func() error {
				clearScreenCallCount++
				return nil
			}
			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrintEmptyLinesMock = func(count int) {}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.in = fakeInput
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := fakeO.wrapAfterExecutionFunc(fakeO)
			assert.Nil(t, err, "expected clear screen to succeed")
			assert.Equal(t, 1, pressAnyKeyCallCount)
			assert.Equal(t, 1, clearScreenCallCount)
		})
	})
}

var InstructionActionSelectionFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				promptCallCount++
				return nil, errors.NewInterruptError(fmt.Errorf("failed to prompt"))
			}

			result, err := fakeO.startInstructionActionSelectionFlowFunc(fakeO, item1, instRootTestData)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected prompt to fail")
			assert.Equal(t, "failed to prompt", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionActionSelectionGoBack = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				promptCallCount++
				return &models.Action{
					Id: prompter.BackActionName,
				}, nil
			}

			result, err := fakeO.startInstructionActionSelectionFlowFunc(fakeO, item1, instRootTestData)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected prompt to succeed")
			assert.Equal(t, prompter.BackActionName, result.Id)
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionActionSelectionFailWorkflowSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				promptCallCount++
				return &models.Action{
					Id: prompter.WorkflowsActionName,
				}, nil
			}

			workflowSelectionCallCount := 0
			fakeO.startInstructionWorkflowSelectionFlowFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				globals *models.Globals,
				workflows []*models.Workflow,
				actions []*models.Action) (*models.Workflow, *errors.PromptError) {
				workflowSelectionCallCount++
				return nil, errors.NewPromptError(fmt.Errorf("failed workflow selection"))
			}

			result, err := fakeO.startInstructionActionSelectionFlowFunc(fakeO, item1, instRootTestData)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected workflow selection to fail")
			assert.Equal(t, "failed workflow selection", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, 1, workflowSelectionCallCount)
		})
	})
}

var InstructionActionSelectionSucceedWorkflowSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				assert.Equal(t, item1, commandFolderItem)
				promptCallCount++
				if promptCallCount == 1 {
					return &models.Action{
						Id: prompter.WorkflowsActionName,
					}, nil
				} else {
					return &models.Action{
						Id: prompter.BackActionName,
					}, nil
				}
			}

			workflowSelectionCallCount := 0
			fakeO.startInstructionWorkflowSelectionFlowFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				globals *models.Globals,
				workflows []*models.Workflow,
				actions []*models.Action) (*models.Workflow, *errors.PromptError) {
				workflowSelectionCallCount++
				assert.Equal(t, item1, commandFolderItem)
				return nil, nil
			}

			result, err := fakeO.startInstructionActionSelectionFlowFunc(fakeO, item1, instRootTestData)
			assert.NotNil(t, result)
			assert.Equal(t, prompter.BackActionName, result.Id)
			assert.Nil(t, err, "expected workflow selection to succeed")
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, workflowSelectionCallCount)
		})
	})
}

var InstructionActionSelectionFailActionExecution = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				assert.Equal(t, item1, commandFolderItem)
				promptCallCount++
				return action1, nil
			}

			instActionExecCallCount := 0
			fakeO.startInstructionActionExecutionFlowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				action *models.Action) (*models.Action, *errors.PromptError) {
				instActionExecCallCount++
				assert.Equal(t, action1, action)
				return nil, errors.NewPromptError(fmt.Errorf("failed to exec instruction action"))
			}

			result, err := fakeO.startInstructionActionSelectionFlowFunc(fakeO, item1, instRootTestData)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected inst action exec to fail")
			assert.Equal(t, "failed to exec instruction action", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, 1, instActionExecCallCount)
		})
	})
}

var InstructionActionSelectionSucceedActionExecution = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				assert.Equal(t, item1, commandFolderItem)
				promptCallCount++
				if promptCallCount == 1 {
					return action1, nil
				} else {
					return &models.Action{
						Id: prompter.BackActionName,
					}, nil
				}
			}

			instActionExecCallCount := 0
			fakeO.startInstructionActionExecutionFlowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				action *models.Action) (*models.Action, *errors.PromptError) {
				instActionExecCallCount++
				assert.Equal(t, action1, action)
				return nil, nil
			}

			result, err := fakeO.startInstructionActionSelectionFlowFunc(fakeO, item1, instRootTestData)
			assert.NotNil(t, result)
			assert.Equal(t, prompter.BackActionName, result.Id)
			assert.Nil(t, err, "expected workflow selection to succeed")
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, instActionExecCallCount)
		})
	})
}

var InstructionActionPromptFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionActionsMock = func(appName string, actions []*models.Action) (*models.Action, error) {
				promptCallCount++
				return nil, fmt.Errorf("failed to prompt for instruction actions")
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptInstructionActionSelectionFunc(fakeO, item1, instRootTestData.Instructions.Actions)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected prompt to fail")
			assert.Equal(t, "failed to prompt for instruction actions", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionActionPromptPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionActionsMock = func(appName string, actions []*models.Action) (*models.Action, error) {
				promptCallCount++
				return action1, nil
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptInstructionActionSelectionFunc(fakeO, item1, instRootTestData.Instructions.Actions)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected prompt to succeed")
			assert.Equal(t, action1, result)
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionActionExecFailToAskBeforeRunning = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, globals *models.Globals, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return false, errors.NewPromptError(fmt.Errorf("failed to ask"))
			}

			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, models.EmptyGlobals(), action1)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected instructions action execution to fail")
			assert.Equal(t, "failed to ask", err.GoError().Error())
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionActionExecFailToRun = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, globals *models.Globals, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeRunner.RunInstructionActionFunc = func(o *runner.ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
				runInstructionCallCount++
				// Schema error should not fail the selection flow
				return errors.NewSchemaError(fmt.Errorf("failed to run instruction"))
			}

			wrapCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapCallCount++
				return nil
			}

			fakePrinter := printer.CreateFakePrinter()
			printErrorCallCount := 0
			fakePrinter.PrintErrorMock = func(message string) {
				printErrorCallCount++
			}

			fakeO.prntr = fakePrinter
			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, models.EmptyGlobals(), action1)
			assert.Nil(t, err, "expected not to fail instructions action selection")
			assert.NotNil(t, result)
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runInstructionCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, printErrorCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, wrapCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionActionExecFailToWrapAfterRun = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, globals *models.Globals, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeRunner.RunInstructionActionFunc = func(o *runner.ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
				runInstructionCallCount++
				return nil
			}

			wrapAfterExecCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapAfterExecCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to wrap"))
			}

			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, models.EmptyGlobals(), action1)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected instructions action execution to fail")
			assert.Equal(t, "failed to wrap", err.GoError().Error())
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runInstructionCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, wrapAfterExecCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionActionExecRunSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, globals *models.Globals, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeRunner.RunInstructionActionFunc = func(o *runner.ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
				runInstructionCallCount++
				return nil
			}

			wrapAfterExecCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapAfterExecCallCount++
				return nil
			}

			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, models.EmptyGlobals(), action1)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected instructions action execution to succeed")
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runInstructionCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, wrapAfterExecCallCount, "expected func to be called exactly once")
			assert.Equal(t, result, action1)
		})
	})
}

var InstructionActionExecFailToAskYesNoQuestion = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return false, fmt.Errorf("failed to ask")
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionActionFunc(fakeO, models.EmptyGlobals(), action1)
			assert.False(t, result)
			assert.NotNil(t, err, "expected ask yes/no question to fail")
			assert.Equal(t, "failed to ask", err.GoError().Error())
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionActionExecAskYesNoQuestionForApplicationContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
			action1.Context = models.ApplicationContext
			globals := models.EmptyGlobals()

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return true, nil
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionActionFunc(fakeO, globals, action1)
			assert.True(t, result)
			assert.Nil(t, err, "expected ask yes/no question to succeed")
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionActionExecAskYesNoQuestionForKubernetesContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
			action1.Context = models.KubernetesContext
			globals := models.EmptyGlobals()

			fakeShell := shell.CreateFakeShell()
			execReturnOutputCallCount := 0
			fakeShell.ExecuteReturnOutputMock = func(script string) (string, error) {
				execReturnOutputCallCount++
				return "", nil
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return true, nil
			}

			fakeO.in = fakeInput
			fakeO.s = fakeShell
			result, err := fakeO.askBeforeRunningInstructionActionFunc(fakeO, globals, action1)
			assert.True(t, result)
			assert.Nil(t, err, "expected ask yes/no question to succeed")
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
			assert.Equal(t, 5, execReturnOutputCallCount)
		})
	})
}

var InstructionWorkflowSelectionFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {
				promptCallCount++
				return nil, errors.NewInterruptError(fmt.Errorf("failed to prompt"))
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
				models.EmptyGlobals(),
				instRootTestData.Instructions.Workflows,
				instRootTestData.Instructions.Actions)

			assert.Nil(t, result)
			assert.NotNil(t, err, "expected prompt to fail")
			assert.Equal(t, "failed to prompt", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionWorkflowSelectionGoBack = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {
				promptCallCount++
				return &models.Workflow{
					Id: prompter.BackActionName,
				}, nil
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
				models.EmptyGlobals(),
				instRootTestData.Instructions.Workflows,
				instRootTestData.Instructions.Actions)

			assert.NotNil(t, result)
			assert.Nil(t, err, "expected prompt to succeed")
			assert.Equal(t, prompter.BackActionName, result.Id)
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionWorkflowSelectionFailWorkflowExecution = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {
				promptCallCount++
				return &models.Workflow{
					Id: prompter.WorkflowsActionName,
				}, nil
			}

			workflowExecCallCount := 0
			fakeO.startInstructionWorkflowExecutionFlowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				workflow *models.Workflow,
				actions []*models.Action) (*models.Workflow, *errors.PromptError) {
				workflowExecCallCount++
				return nil, errors.NewPromptError(fmt.Errorf("failed workflow execution"))
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
				models.EmptyGlobals(),
				instRootTestData.Instructions.Workflows,
				instRootTestData.Instructions.Actions)

			assert.Nil(t, result)
			assert.NotNil(t, err, "expected workflow selection to fail")
			assert.Equal(t, "failed workflow execution", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, 1, workflowExecCallCount)
		})
	})
}

var InstructionWorkflowSelectionSucceedWorkflowExecution = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				commandFolderItem *models.CommandFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {

				assert.Equal(t, item1, commandFolderItem)
				promptCallCount++
				if promptCallCount == 1 {
					return app1Workflow1, nil
				} else {
					return &models.Workflow{
						Id: prompter.BackActionName,
					}, nil
				}
			}

			workflowExecCallCount := 0
			fakeO.startInstructionWorkflowExecutionFlowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				workflow *models.Workflow,
				actions []*models.Action) (*models.Workflow, *errors.PromptError) {
				workflowExecCallCount++
				assert.Equal(t, app1Workflow1, workflow)
				return nil, nil
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
				models.EmptyGlobals(),
				instRootTestData.Instructions.Workflows,
				instRootTestData.Instructions.Actions)

			assert.NotNil(t, result)
			assert.Equal(t, prompter.BackActionName, result.Id)
			assert.Nil(t, err, "expected workflow selection to succeed")
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, workflowExecCallCount)
		})
	})
}

var InstructionWorkflowExecFailToAskYesNoQuestion = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return false, fmt.Errorf("failed to ask")
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionWorkflowFunc(fakeO, models.EmptyGlobals(), app1Workflow1)
			assert.False(t, result)
			assert.NotNil(t, err, "expected ask yes/no question to fail")
			assert.Equal(t, "failed to ask", err.GoError().Error())
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowExecAskYesNoQuestionForApplicationContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)
			globals := models.EmptyGlobals()
			app1Workflow1.Context = models.ApplicationContext

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return true, nil
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionWorkflowFunc(fakeO, globals, app1Workflow1)
			assert.True(t, result)
			assert.Nil(t, err, "expected ask yes/no question to succeed")
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowExecAskYesNoQuestionForKubernetesContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)
			globals := models.EmptyGlobals()
			app1Workflow1.Context = models.KubernetesContext

			fakeShell := shell.CreateFakeShell()
			execReturnOutputCallCount := 0
			fakeShell.ExecuteReturnOutputMock = func(script string) (string, error) {
				execReturnOutputCallCount++
				return "", nil
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return true, nil
			}

			fakeO.in = fakeInput
			fakeO.s = fakeShell
			result, err := fakeO.askBeforeRunningInstructionWorkflowFunc(fakeO, globals, app1Workflow1)
			assert.True(t, result)
			assert.Nil(t, err, "expected ask yes/no question to succeed")
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
			assert.Equal(t, 5, execReturnOutputCallCount)
		})
	})
}

var InstructionWorkflowPromptFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionWorkflowsMock = func(appName string, workflows []*models.Workflow) (*models.Workflow, error) {
				promptCallCount++
				return nil, fmt.Errorf("failed to prompt for instruction workflow")
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptInstructionWorkflowSelectionFunc(fakeO, item1, instRootTestData.Instructions.Workflows)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected prompt to fail")
			assert.Equal(t, "failed to prompt for instruction workflow", err.GoError().Error())
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionWorkflowPromptPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionWorkflowsMock = func(appName string, workflows []*models.Workflow) (*models.Workflow, error) {
				promptCallCount++
				return app1Workflow1, nil
			}

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptInstructionWorkflowSelectionFunc(fakeO, item1, instRootTestData.Instructions.Workflows)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected prompt to succeed")
			assert.Equal(t, app1Workflow1, result)
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var InstructionWorkflowExecFailToAskBeforeRunning = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return false, errors.NewPromptError(fmt.Errorf("failed to ask"))
			}

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, models.EmptyGlobals(), app1Workflow1, instRootTestData.Instructions.Actions)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected instructions workflow execution to fail")
			assert.Equal(t, "failed to ask", err.GoError().Error())
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowExecTolerateRunFailures = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeRunner.RunInstructionWorkflowFunc = func(
				o *runner.ActionRunnerOrchestrator,
				workflow *models.Workflow,
				actions []*models.Action) *errors.PromptError {

				runInstructionCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to run instruction workflow"))
			}

			wrapAfterExecCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapAfterExecCallCount++
				return nil
			}

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, models.EmptyGlobals(), app1Workflow1, instRootTestData.Instructions.Actions)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected instructions workflow execution to succeed")
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runInstructionCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, wrapAfterExecCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowExecFailToWrapAfterRun = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeRunner.RunInstructionWorkflowFunc = func(
				o *runner.ActionRunnerOrchestrator,
				workflow *models.Workflow,
				actions []*models.Action) *errors.PromptError {

				runInstructionCallCount++
				return nil
			}

			wrapAfterExecCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapAfterExecCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to wrap"))
			}

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, models.EmptyGlobals(), app1Workflow1, instRootTestData.Instructions.Actions)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected instructions action execution to fail")
			assert.Equal(t, "failed to wrap", err.GoError().Error())
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runInstructionCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, wrapAfterExecCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowExecRunSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				globals *models.Globals,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeRunner.RunInstructionWorkflowFunc = func(
				o *runner.ActionRunnerOrchestrator,
				workflow *models.Workflow,
				actions []*models.Action) *errors.PromptError {

				runInstructionCallCount++
				return nil
			}

			wrapAfterExecCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapAfterExecCallCount++
				return nil
			}

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, models.EmptyGlobals(), app1Workflow1, instRootTestData.Instructions.Actions)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected instructions action execution to succeed")
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runInstructionCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, wrapAfterExecCallCount, "expected func to be called exactly once")
			assert.Equal(t, result, app1Workflow1)
		})
	})
}

var ManagePromptErrorMissingInnerGoError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptError := errors.NewPromptError(nil)
			err := managePromptError(promptError)
			assert.Nil(t, err, "expect no inner go error")
		})
	})
}

var ManagePromptErrorMitigateInterruptError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptError := errors.NewPromptError(promptui.ErrInterrupt)
			err := managePromptError(promptError)
			assert.Nil(t, err, "expect no inner go error")
		})
	})
}

var ManagePromptErrorReturnInnerGoError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptError := errors.NewPromptError(fmt.Errorf("inner go error"))
			err := managePromptError(promptError)
			assert.NotNil(t, err, "expect error to exist")
			assert.Equal(t, err.Error(), "inner go error")
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

var DoNotAddBackOptionToWorkflowPromptSelectorIfAlreadyExist = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	appendInstructionWorkflowsCustomOptions(instRootTestData.Instructions)
	appendInstructionWorkflowsCustomOptions(instRootTestData.Instructions)
	assert.EqualValues(t, prompter.BackActionName, instRootTestData.Instructions.Workflows[0].Id)
	assert.NotEqualValues(t, prompter.BackActionName, instRootTestData.Instructions.Workflows[1].Id)
}

var ExtractArgsFromScriptFile = func(t *testing.T) {
	expectedPath := "/path/to/script"
	scriptNoArgs := expectedPath
	path, args := extractArgsFromScriptFile(scriptNoArgs)
	assert.Equal(t, "/path/to/script", path)
	assert.Nil(t, args)

	scriptMultipleArgs := fmt.Sprintf("%s %s %s", expectedPath, "--create", "--deploy")
	path, args = extractArgsFromScriptFile(scriptMultipleArgs)
	assert.Equal(t, "/path/to/script", path)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, "--create", args[0])
	assert.Equal(t, "--deploy", args[1])
}

var ResolveInstructionActionContextSuccessfully = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	action := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
	globals := models.EmptyGlobals()

	globals.Context = models.KubernetesContext
	action.Context = ""
	ctxResult := getInstructionActionContext(globals, action)
	assert.Equal(t, ctxResult, models.KubernetesContext)

	globals.Context = ""
	action.Context = models.ApplicationContext
	ctxResult = getInstructionActionContext(globals, action)
	assert.Equal(t, ctxResult, models.ApplicationContext)

	globals.Context = ""
	action.Context = ""
	ctxResult = getInstructionActionContext(globals, action)
	assert.Equal(t, ctxResult, models.ApplicationContext)
}

var ResolveInstructionWorkflowContextSuccessfully = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	workflow := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)
	globals := models.EmptyGlobals()

	globals.Context = models.KubernetesContext
	workflow.Context = ""
	ctxResult := getInstructionWorkflowContext(globals, workflow)
	assert.Equal(t, ctxResult, models.KubernetesContext)

	globals.Context = ""
	workflow.Context = models.ApplicationContext
	ctxResult = getInstructionWorkflowContext(globals, workflow)
	assert.Equal(t, ctxResult, models.ApplicationContext)

	globals.Context = ""
	workflow.Context = ""
	ctxResult = getInstructionWorkflowContext(globals, workflow)
	assert.Equal(t, ctxResult, models.ApplicationContext)
}
