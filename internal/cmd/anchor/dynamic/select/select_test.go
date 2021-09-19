package _select

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
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
			Name: "anchor folder item selection: cancel selection successfully",
			Func: AnchorFolderItemSelectionCancelSelectionSuccessfully,
		},
		{
			Name: "anchor folder item selection: fail to prompt",
			Func: AnchorFolderItemSelectionFailToPrompt,
		},
		{
			Name: "anchor folder item selection: fail to extract instructions",
			Func: AnchorFolderItemSelectionFailToExtractInstructions,
		},
		{
			Name: "anchor folder item selection: go back from instruction action selection",
			Func: AnchorFolderItemSelectionGoBackFromInstructionActionSelection,
		},
		{
			Name: "anchor folder item selection: instruction action selection missing error",
			Func: AnchorFolderItemSelectionInstructionActionSelectionMissingError,
		},
		{
			Name: "anchor folder item selection: complete successfully",
			Func: AnchorFolderItemSelectionCompleteSuccessfully,
		},
		{
			Name: "anchor folder item prompt: fail to prompt",
			Func: AnchorFolderItemPromptFailToPrompt,
		},
		{
			Name: "anchor folder item prompt: prompt successfully",
			Func: AnchorFolderItemPromptPromptSuccessfully,
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
			Name: "extract instructions: fail when missing",
			Func: ExtractInstructionsFailWhenMissing,
		},
		{
			Name: "extract instructions: prompt empty when invalid schema",
			Func: ExtractInstructionsPromptEmptyWhenInvalidSchema,
		},
		{
			Name: "extract instructions: prompt enriched actions",
			Func: ExtractInstructionsPromptEnrichedActions,
		},
		{
			Name: "run instruction action: fail on mutual exclusive scripts origin",
			Func: RunInstructionActionFailOnMutualExclusiveScriptPaths,
		},
		{
			Name: "run instruction action: fail on missing script to run",
			Func: RunInstructionActionFailOnMissingScriptToExec,
		},
		{
			Name: "run instruction action: fail to execute action",
			Func: RunInstructionActionFailToExecuteAction,
		},
		{
			Name: "run instruction action: fail to execute verbose action",
			Func: RunInstructionActionFailToExecuteVerboseAction,
		},
		{
			Name: "run instruction action: run action with verbose from flag",
			Func: RunInstructionActionRunActionWithVerboseFromFlag,
		},
		{
			Name: "run instruction action: run action with forced verbose from schema",
			Func: RunInstructionActionRunActionWithForcedVerboseFromSchema,
		},
		{
			Name: "run instruction action: run action interactive",
			Func: RunInstructionActionRunActionInteractive,
		},
		{
			Name: "action exec interactive: fail to exec script",
			Func: ActionExecInteractiveFailToExecScript,
		},
		{
			Name: "action exec interactive: exec script successfully",
			Func: ActionExecInteractiveExecScriptSuccessfully,
		},
		{
			Name: "action exec interactive: fail to exec script file",
			Func: ActionExecInteractiveFailToExecScriptFile,
		},
		{
			Name: "action exec interactive: exec script file successfully",
			Func: ActionExecInteractiveExecScriptFileSuccessfully,
		},
		{
			Name: "action exec interactive: no op when nothing to exec",
			Func: ActionExecInteractiveNoOpWhenNothingToExec,
		},
		{
			Name: "action exec verbose: fail to exec script",
			Func: ActionExecVerboseFailToExecScript,
		},
		{
			Name: "action exec verbose: exec script successfully",
			Func: ActionExecVerboseExecScriptSuccessfully,
		},
		{
			Name: "action exec verbose: fail to exec script file",
			Func: ActionExecVerboseFailToExecScriptFile,
		},
		{
			Name: "action exec verbose: exec script file successfully",
			Func: ActionExecVerboseExecScriptFileSuccessfully,
		},
		{
			Name: "action exec verbose: no op when nothing to exec",
			Func: ActionExecVerboseNoOpWhenNothingToExec,
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
			Name: "instruction action exec flow: ask yes/no question successfully",
			Func: InstructionActionExecAskYesNoQuestionSuccessfully,
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
			Name: "instruction workflow exec: ask yes/no question successfully",
			Func: InstructionWorkflowExecAskYesNoQuestionSuccessfully,
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
			Name: "run instruction workflow: fail to run actions",
			Func: RunInstructionWorkflowFailToRunActions,
		},
		{
			Name: "run instruction workflow: run actions successfully",
			Func: RunInstructionWorkflowRunActionsSuccessfully,
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
			Name: "extract args from script file",
			Func: ExtractArgsFromScriptFile,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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
			assert.Nil(t, err, "expected not to fail anchor folder item status")
			assert.Equal(t, 1, bannerCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var StartTheSelectionFlowWhenRunIsCalled = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			startFolderItemSelectionCallCount := 0
			fakeO.startFolderItemsSelectionFlowFunc = func(o *selectOrchestrator, anchorfilesRepoPath string) *errors.PromptError {
				startFolderItemSelectionCallCount++
				return nil
			}
			err := fakeO.runFunc(fakeO, ctx)
			assert.Nil(t, err, "expected not to fail anchor folder item status")
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

	fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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

		fakeExtractor := extractor.CreateFakeExtractor()
		reg.Set(extractor.Identifier, fakeExtractor)

		fakeParser := parser.CreateFakeParser()
		reg.Set(parser.Identifier, fakeParser)

		fakeShell := shell.CreateFakeShell()
		reg.Set(shell.Identifier, fakeShell)

		fakeInput := input.CreateFakeUserInput()
		reg.Set(input.Identifier, fakeInput)

		fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
		err := fakeO.prepareFunc(fakeO, ctx)

		assert.Nil(t, err)
		assert.NotNil(t, fakeO.l)
		assert.NotNil(t, fakeO.prmpt)
		assert.NotNil(t, fakeO.e)
		assert.NotNil(t, fakeO.prsr)
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
		fakeExtractor := extractor.CreateFakeExtractor()
		fakeParser := parser.CreateFakeParser()
		fakeShell := shell.CreateFakeShell()
		fakeInput := input.CreateFakeUserInput()

		fakeO := NewOrchestrator(stubs.AnchorFolder1Name)

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
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", extractor.Identifier))
		reg.Set(extractor.Identifier, fakeExtractor)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", parser.Identifier))
		reg.Set(parser.Identifier, fakeParser)

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
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			prepareRegistryItemsCallCount := 0
			fakeO.prepareFunc = func(o *selectOrchestrator, ctx common.Context) error {
				prepareRegistryItemsCallCount++
				return fmt.Errorf("failed to prepare registry items")
			}
			err := DynamicSelect(ctx, fakeO)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare registry items", err.Error())
			assert.Equal(t, 1, prepareRegistryItemsCallCount, "expected func to be called exactly once")
		})
	})
}

var FailToRunSelectionForDynamicCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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

var AnchorFolderItemSelectionCancelSelectionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptFolderItemsSelectionFunc = func(o *selectOrchestrator) (*models.AnchorFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				return &models.AnchorFolderItemInfo{
					Name: prompter.CancelActionName,
				}, nil
			}
			err := fakeO.startFolderItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.Nil(t, err, "expected selection to succeed")
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var AnchorFolderItemSelectionFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptFolderItemsSelectionFunc = func(o *selectOrchestrator) (*models.AnchorFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				return nil, errors.NewPromptError(fmt.Errorf("failed to prompt"))
			}
			err := fakeO.startFolderItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection prompt to fail")
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, "failed to prompt", err.GoError().Error())
		})
	})
}

var AnchorFolderItemSelectionFailToExtractInstructions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptFolderItemsSelectionFunc = func(o *selectOrchestrator) (*models.AnchorFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				if promptCallCount == 1 {
					return item1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractorCallCount := 0
			fakeO.extractInstructionsFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				assert.Equal(t, anchorFolderItem.InstructionsPath, item1.InstructionsPath)
				return nil, errors.NewPromptError(fmt.Errorf("failed to extract instructions"))
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrintMissingInstructionsMock = func() {}
			fakeO.prntr = fakePrinter

			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				return nil
			}

			err := fakeO.startFolderItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
		})
	})
}

var AnchorFolderItemSelectionGoBackFromInstructionActionSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptFolderItemsSelectionFunc = func(o *selectOrchestrator) (*models.AnchorFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				if promptCallCount == 1 {
					return item1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractorCallCount := 0
			fakeO.extractInstructionsFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				return instRootTestData, nil
			}

			instActionSelectionCallCount := 0
			fakeO.startInstructionActionSelectionFlowFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {
				instActionSelectionCallCount++
				return &models.Action{
					Id: prompter.BackActionName,
				}, nil
			}

			err := fakeO.startFolderItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instActionSelectionCallCount)
		})
	})
}

var AnchorFolderItemSelectionInstructionActionSelectionMissingError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptFolderItemsSelectionFunc = func(o *selectOrchestrator) (*models.AnchorFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				if promptCallCount == 1 {
					return item1, nil
				} else {
					// Use programmatic keyboard interrupt to stop the selection flow
					return nil, errors.NewInterruptError(fmt.Errorf("keyboard interrupt the test flow"))
				}
			}

			extractorCallCount := 0
			fakeO.extractInstructionsFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				return instRootTestData, nil
			}

			instActionSelectionCallCount := 0
			fakeO.startInstructionActionSelectionFlowFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {
				instActionSelectionCallCount++
				return nil, errors.NewInstructionMissingError(fmt.Errorf("missing instruction"))
			}

			err := fakeO.startFolderItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.NotNil(t, err, "expected selection to stop due to keyboard interrupt")
			assert.Equal(t, "keyboard interrupt the test flow", err.GoError().Error())
			assert.Equal(t, 2, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instActionSelectionCallCount)
		})
	})
}

var AnchorFolderItemSelectionCompleteSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptFolderItemsSelectionFunc = func(o *selectOrchestrator) (*models.AnchorFolderItemInfo, *errors.PromptError) {
				promptCallCount++
				return item1, nil
			}

			extractorCallCount := 0
			fakeO.extractInstructionsFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
				extractorCallCount++
				assert.Equal(t, item1, anchorFolderItem)
				return instRootTestData, nil
			}

			instActionSelectionCallCount := 0
			fakeO.startInstructionActionSelectionFlowFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {
				instActionSelectionCallCount++
				assert.Equal(t, item1, anchorFolderItem)
				return stubs.GetInstructionActionById(instructionRoot.Instructions, stubs.AnchorFolder1Item1Action1Id), nil
			}

			err := fakeO.startFolderItemsSelectionFlowFunc(fakeO, ctx.AnchorFilesPath())
			assert.Nil(t, err, "expected selection to complete successfully")
			assert.Equal(t, 1, promptCallCount)
			assert.Equal(t, 1, extractorCallCount)
			assert.Equal(t, 1, instActionSelectionCallCount)
		})
	})
}

var AnchorFolderItemPromptFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()

			fakeLocator := locator.CreateFakeLocator()
			locateAnchorFolderItemsCallCount := 0
			fakeLocator.AnchorFolderItemsMock = func(parentFolderName string) []*models.AnchorFolderItemInfo {
				locateAnchorFolderItemsCallCount++
				return items
			}

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptAnchorFolderItemSelectionMock = func(appsArr []*models.AnchorFolderItemInfo) (*models.AnchorFolderItemInfo, error) {
				promptCallCount++
				return nil, fmt.Errorf("failed to prompt for items")
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.l = fakeLocator
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptFolderItemsSelectionFunc(fakeO)
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected prompt to fail")
			assert.Equal(t, "failed to prompt for items", err.GoError().Error())
			assert.Equal(t, 1, locateAnchorFolderItemsCallCount)
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var AnchorFolderItemPromptPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)

			fakeLocator := locator.CreateFakeLocator()
			locateAnchorFolderItemsCallCount := 0
			fakeLocator.AnchorFolderItemsMock = func(parentFolderName string) []*models.AnchorFolderItemInfo {
				locateAnchorFolderItemsCallCount++
				return items
			}

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptAnchorFolderItemSelectionMock = func(appsArr []*models.AnchorFolderItemInfo) (*models.AnchorFolderItemInfo, error) {
				promptCallCount++
				return item1, nil
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.l = fakeLocator
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptFolderItemsSelectionFunc(fakeO)
			assert.Nil(t, err, "expected prompt to succeed")
			assert.NotNil(t, result, "expected prompt response")
			assert.Equal(t, 1, locateAnchorFolderItemsCallCount)
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

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				promptCallCount++
				return &models.Action{
					Id: prompter.WorkflowsActionName,
				}, nil
			}

			workflowSelectionCallCount := 0
			fakeO.startInstructionWorkflowSelectionFlowFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				assert.Equal(t, item1, anchorFolderItem)
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
				anchorFolderItem *models.AnchorFolderItemInfo,
				workflows []*models.Workflow,
				actions []*models.Action) (*models.Workflow, *errors.PromptError) {
				workflowSelectionCallCount++
				assert.Equal(t, item1, anchorFolderItem)
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				assert.Equal(t, item1, anchorFolderItem)
				promptCallCount++
				return action1, nil
			}

			instActionExecCallCount := 0
			fakeO.startInstructionActionExecutionFlowFunc = func(o *selectOrchestrator, action *models.Action) (*models.Action, *errors.PromptError) {
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionActionSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				actions []*models.Action) (*models.Action, *errors.PromptError) {
				assert.Equal(t, item1, anchorFolderItem)
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
			fakeO.startInstructionActionExecutionFlowFunc = func(o *selectOrchestrator, action *models.Action) (*models.Action, *errors.PromptError) {
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionActionsMock = func(appName string, actions []*models.Action) (*models.Action, error) {
				promptCallCount++
				return nil, fmt.Errorf("failed to prompt for instruction actions")
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionActionsMock = func(appName string, actions []*models.Action) (*models.Action, error) {
				promptCallCount++
				return action1, nil
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.prmpt = fakePrompter

			result, err := fakeO.promptInstructionActionSelectionFunc(fakeO, item1, instRootTestData.Instructions.Actions)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected prompt to succeed")
			assert.Equal(t, action1, result)
			assert.Equal(t, 1, promptCallCount)
		})
	})
}

var ExtractInstructionsFailWhenMissing = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)

			extractInstructionsCallCount := 0
			fakeExtractor := extractor.CreateFakeExtractor()
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractInstructionsCallCount++
				return nil, fmt.Errorf("failed to extract instructions")
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.e = fakeExtractor

			result, err := fakeO.extractInstructionsFunc(fakeO, item1, ctx.AnchorFilesPath())
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected instructions extraction to fail")
			assert.Equal(t, errors.InstructionMissingError, err.Code())
			assert.Equal(t, 1, extractInstructionsCallCount)
		})
	})
}

var ExtractInstructionsPromptEmptyWhenInvalidSchema = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)

			extractInstructionsCallCount := 0
			fakeExtractor := extractor.CreateFakeExtractor()
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractInstructionsCallCount++
				assert.Equal(t, item1.InstructionsPath, instructionsPath)
				return nil, nil
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.e = fakeExtractor

			result, err := fakeO.extractInstructionsFunc(fakeO, item1, ctx.AnchorFilesPath())
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected instructions extraction to succeed")
			assert.Equal(t, 1, extractInstructionsCallCount)
			assert.Equal(t, models.EmptyInstructionsRoot(), result)
		})
	})
}

var ExtractInstructionsPromptEnrichedActions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.HarnessAnchorfilesTestRepo(ctx)
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			extractInstructionsCallCount := 0
			fakeExtractor := extractor.CreateFakeExtractor()
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractInstructionsCallCount++
				assert.Equal(t, item1.InstructionsPath, instructionsPath)
				return instRootTestData, nil
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.e = fakeExtractor

			result, err := fakeO.extractInstructionsFunc(fakeO, item1, ctx.AnchorFilesPath())
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected instructions extraction to succeed")
			assert.Equal(t, 1, extractInstructionsCallCount)
			assert.Contains(t, ctx.AnchorFilesPath(), result.Instructions.Actions[0].AnchorfilesRepoPath)
			assert.Contains(t, ctx.AnchorFilesPath(), result.Instructions.Actions[1].AnchorfilesRepoPath)
		})
	})
}

var RunInstructionActionFailOnMutualExclusiveScriptPaths = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			action1.Script = "some script"
			action1.ScriptFile = "/path/to/script"
			err := fakeO.runInstructionActionFunc(fakeO, action1)
			assert.NotNil(t, err, "expected instructions execution to fail")
			assert.Contains(t, err.GoError().Error(), "script / scriptFile are mutual exclusive")
		})
	})
}

var RunInstructionActionFailOnMissingScriptToExec = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)

			action1.Script = ""
			action1.ScriptFile = ""
			err := fakeO.runInstructionActionFunc(fakeO, action1)
			assert.NotNil(t, err, "expected instructions execution to fail")
			assert.Contains(t, err.GoError().Error(), "missing script or scriptFile")
		})
	})
}

var RunInstructionActionFailToExecuteAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.verboseFlag = true

			execActionVerboseCallCount := 0
			fakeO.executeInstructionActionVerboseFunc = func(o *selectOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionVerboseCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to execute action"))
			}

			action1.Script = "some script"
			action1.ScriptFile = ""
			err := fakeO.runInstructionActionFunc(fakeO, action1)
			assert.NotNil(t, err, "expected instructions execution to fail")
			assert.Contains(t, err.GoError().Error(), "failed to execute action")
			assert.Equal(t, 1, execActionVerboseCallCount)
		})
	})
}

var RunInstructionActionFailToExecuteVerboseAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.verboseFlag = false

			execActionCallCount := 0
			fakeO.executeInstructionActionFunc = func(o *selectOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to execute action"))
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"
			err := fakeO.runInstructionActionFunc(fakeO, action1)
			assert.NotNil(t, err, "expected instructions execution to fail")
			assert.Contains(t, err.GoError().Error(), "failed to execute action")
			assert.Equal(t, 1, execActionCallCount, "expected func to be called exactly once")
		})
	})
}

var RunInstructionActionRunActionWithVerboseFromFlag = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.verboseFlag = true
			execActionVerboseCallCount := 0
			fakeO.executeInstructionActionVerboseFunc = func(o *selectOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionVerboseCallCount++
				return nil
			}

			action1.Script = "some script with verbose from flag"
			action1.ScriptFile = ""
			err := fakeO.runInstructionActionFunc(fakeO, action1)
			assert.Nil(t, err)
			assert.Equal(t, 1, execActionVerboseCallCount, "expected to be called exactly once")
		})
	})
}

var RunInstructionActionRunActionWithForcedVerboseFromSchema = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			action1.ForceVerbose = true

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.verboseFlag = false
			execActionVerboseCallCount := 0
			fakeO.executeInstructionActionVerboseFunc = func(o *selectOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionVerboseCallCount++
				return nil
			}

			action1.Script = ""
			action1.ScriptFile = "some script with forced verbose from schema"
			err := fakeO.runInstructionActionFunc(fakeO, action1)
			assert.Nil(t, err)
			assert.Equal(t, 1, execActionVerboseCallCount, "expected to be called exactly once")
		})
	})
}

var RunInstructionActionRunActionInteractive = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.verboseFlag = false
			execActionCallCount := 0
			fakeO.executeInstructionActionFunc = func(o *selectOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionCallCount++
				return nil
			}
			action1.Script = "some script"
			action1.ScriptFile = ""
			err := fakeO.runInstructionActionFunc(fakeO, action1)
			assert.Nil(t, err)
			assert.Equal(t, 1, execActionCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecInteractiveFailToExecScript = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakeSpinner := printer.CreateFakePrinterSpinner()
			spinCallCount := 0
			fakeSpinner.SpinMock = func() {
				spinCallCount++
			}

			stopOnFailureCallCount := 0
			fakeSpinner.StopOnFailureMock = func(err error) {
				stopOnFailureCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionSpinnerMock = func(actionId string, scriptOutputPath string) printer.PrinterSpinner {
				return fakeSpinner
			}

			execCallCount := 0
			fakeShell.ExecuteSilentlyWithOutputToFileMock = func(script string, outputFilePath string) error {
				execCallCount++
				assert.Equal(t, "some script", script)
				return fmt.Errorf("fail to execute")
			}

			action1.Script = "some script"
			action1.ScriptFile = ""

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionAction(fakeO, action1, "")
			assert.NotNil(t, err)
			assert.Equal(t, 1, spinCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnFailureCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecInteractiveExecScriptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakeSpinner := printer.CreateFakePrinterSpinner()
			spinCallCount := 0
			fakeSpinner.SpinMock = func() {
				spinCallCount++
			}

			stopOnSuccessCallCount := 0
			fakeSpinner.StopOnSuccessMock = func() {
				stopOnSuccessCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionSpinnerMock = func(actionId string, scriptOutputPath string) printer.PrinterSpinner {
				return fakeSpinner
			}

			execCallCount := 0
			fakeShell.ExecuteSilentlyWithOutputToFileMock = func(script string, outputFilePath string) error {
				execCallCount++
				assert.Equal(t, "some script", script)
				return nil
			}

			action1.Script = "some script"
			action1.ScriptFile = ""

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionAction(fakeO, action1, "")
			assert.Nil(t, err)
			assert.Equal(t, 1, spinCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnSuccessCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecInteractiveFailToExecScriptFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakeSpinner := printer.CreateFakePrinterSpinner()
			spinCallCount := 0
			fakeSpinner.SpinMock = func() {
				spinCallCount++
			}

			stopOnFailureCallCount := 0
			fakeSpinner.StopOnFailureMock = func(err error) {
				stopOnFailureCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionSpinnerMock = func(actionId string, scriptOutputPath string) printer.PrinterSpinner {
				return fakeSpinner
			}

			execCallCount := 0
			fakeShell.ExecuteScriptFileSilentlyWithOutputToFileMock = func(workingDirectory string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return fmt.Errorf("fail to execute")
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionAction(fakeO, action1, "")
			assert.NotNil(t, err, "expected  to fail")
			assert.Equal(t, 1, spinCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnFailureCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecInteractiveExecScriptFileSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakeSpinner := printer.CreateFakePrinterSpinner()
			spinCallCount := 0
			fakeSpinner.SpinMock = func() {
				spinCallCount++
			}

			stopOnSuccessCallCount := 0
			fakeSpinner.StopOnSuccessMock = func() {
				stopOnSuccessCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionSpinnerMock = func(actionId string, scriptOutputPath string) printer.PrinterSpinner {
				return fakeSpinner
			}

			execCallCount := 0
			fakeShell.ExecuteScriptFileSilentlyWithOutputToFileMock = func(workingDirectory string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return nil
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionAction(fakeO, action1, "")
			assert.Nil(t, err)
			assert.Equal(t, 1, spinCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnSuccessCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecInteractiveNoOpWhenNothingToExec = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeSpinner := printer.CreateFakePrinterSpinner()
			spinCallCount := 0
			fakeSpinner.SpinMock = func() {
				spinCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionSpinnerMock = func(actionId string, scriptOutputPath string) printer.PrinterSpinner {
				return fakeSpinner
			}

			action1.Script = ""
			action1.ScriptFile = ""

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.prntr = fakePrinter

			err := executeInstructionAction(fakeO, action1, "")
			assert.Nil(t, err)
			assert.Equal(t, 0, spinCallCount, "expected no calls to be made")
		})
	})
}

var ActionExecVerboseFailToExecScript = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakePlainer := printer.CreateFakePrinterPlainer()
			startCallCount := 0
			fakePlainer.StartMock = func() {
				startCallCount++
			}

			stopOnFailureCallCount := 0
			fakePlainer.StopOnFailureMock = func(err error) {
				stopOnFailureCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionPlainerMock = func(actionId string) printer.PrinterPlainer {
				return fakePlainer
			}

			execCallCount := 0
			fakeShell.ExecuteWithOutputToFileMock = func(script string, outputFilePath string) error {
				execCallCount++
				assert.Equal(t, "some script", script)
				return fmt.Errorf("fail to execute")
			}

			action1.Script = "some script"
			action1.ScriptFile = ""

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.verboseFlag = true
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionActionVerbose(fakeO, action1, "")
			assert.NotNil(t, err)
			assert.Equal(t, 1, startCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnFailureCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecVerboseExecScriptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakePlainer := printer.CreateFakePrinterPlainer()
			startCallCount := 0
			fakePlainer.StartMock = func() {
				startCallCount++
			}

			stopOnSuccessCallCount := 0
			fakePlainer.StopOnSuccessMock = func() {
				stopOnSuccessCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionPlainerMock = func(actionId string) printer.PrinterPlainer {
				return fakePlainer
			}

			execCallCount := 0
			fakeShell.ExecuteWithOutputToFileMock = func(script string, outputFilePath string) error {
				execCallCount++
				assert.Equal(t, "some script", script)
				return nil
			}

			action1.Script = "some script"
			action1.ScriptFile = ""

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionActionVerbose(fakeO, action1, "")
			assert.Nil(t, err)
			assert.Equal(t, 1, startCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnSuccessCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecVerboseFailToExecScriptFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakePlainer := printer.CreateFakePrinterPlainer()
			startCallCount := 0
			fakePlainer.StartMock = func() {
				startCallCount++
			}

			stopOnFailureCallCount := 0
			fakePlainer.StopOnFailureMock = func(err error) {
				stopOnFailureCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionPlainerMock = func(actionId string) printer.PrinterPlainer {
				return fakePlainer
			}

			execCallCount := 0
			fakeShell.ExecuteScriptFileWithOutputToFileMock = func(workingDirectory string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return fmt.Errorf("fail to execute")
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionActionVerbose(fakeO, action1, "")
			assert.NotNil(t, err, "expected  to fail")
			assert.Equal(t, 1, startCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnFailureCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecVerboseExecScriptFileSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			fakeShell := shell.CreateFakeShell()

			fakePlainer := printer.CreateFakePrinterPlainer()
			startCallCount := 0
			fakePlainer.StartMock = func() {
				startCallCount++
			}

			stopOnSuccessCallCount := 0
			fakePlainer.StopOnSuccessMock = func() {
				stopOnSuccessCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionPlainerMock = func(actionId string) printer.PrinterPlainer {
				return fakePlainer
			}

			execCallCount := 0
			fakeShell.ExecuteScriptFileWithOutputToFileMock = func(workingDirectory string, relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return nil
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionActionVerbose(fakeO, action1, "")
			assert.Nil(t, err)
			assert.Equal(t, 1, startCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnSuccessCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecVerboseNoOpWhenNothingToExec = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakePlainer := printer.CreateFakePrinterPlainer()
			startCallCount := 0
			fakePlainer.StartMock = func() {
				startCallCount++
			}

			fakePrinter := printer.CreateFakePrinter()
			fakePrinter.PrepareRunActionPlainerMock = func(actionId string) printer.PrinterPlainer {
				return fakePlainer
			}

			action1.Script = ""
			action1.ScriptFile = ""

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.prntr = fakePrinter

			err := executeInstructionActionVerbose(fakeO, action1, "")
			assert.Nil(t, err)
			assert.Equal(t, 0, startCallCount, "expected no calls to be made")
		})
	})
}

var InstructionActionExecFailToAskBeforeRunning = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return false, errors.NewPromptError(fmt.Errorf("failed to ask"))
			}

			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, action1)
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
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeO.runInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) *errors.PromptError {
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
			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, action1)
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
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeO.runInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) *errors.PromptError {
				runInstructionCallCount++
				return nil
			}

			wrapAfterExecCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapAfterExecCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to wrap"))
			}

			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, action1)
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
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) (bool, *errors.PromptError) {
				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeO.runInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) *errors.PromptError {
				runInstructionCallCount++
				return nil
			}

			wrapAfterExecCallCount := 0
			fakeO.wrapAfterExecutionFunc = func(o *selectOrchestrator) *errors.PromptError {
				wrapAfterExecCallCount++
				return nil
			}

			result, err := fakeO.startInstructionActionExecutionFlowFunc(fakeO, action1)
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
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return false, fmt.Errorf("failed to ask")
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionActionFunc(fakeO, action1)
			assert.False(t, result)
			assert.NotNil(t, err, "expected ask yes/no question to fail")
			assert.Equal(t, "failed to ask", err.GoError().Error())
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionActionExecAskYesNoQuestionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return true, nil
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionActionFunc(fakeO, action1)
			assert.True(t, result)
			assert.Nil(t, err, "expected ask yes/no question to succeed")
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowSelectionFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {
				promptCallCount++
				return nil, errors.NewInterruptError(fmt.Errorf("failed to prompt"))
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {
				promptCallCount++
				return &models.Workflow{
					Id: prompter.BackActionName,
				}, nil
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {
				promptCallCount++
				return &models.Workflow{
					Id: prompter.WorkflowsActionName,
				}, nil
			}

			workflowExecCallCount := 0
			fakeO.startInstructionWorkflowExecutionFlowFunc = func(
				o *selectOrchestrator,
				workflow *models.Workflow,
				actions []*models.Action) (*models.Workflow, *errors.PromptError) {
				workflowExecCallCount++
				return nil, errors.NewPromptError(fmt.Errorf("failed workflow execution"))
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			promptCallCount := 0
			fakeO.promptInstructionWorkflowSelectionFunc = func(
				o *selectOrchestrator,
				anchorFolderItem *models.AnchorFolderItemInfo,
				workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {

				assert.Equal(t, item1, anchorFolderItem)
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
				workflow *models.Workflow,
				actions []*models.Action) (*models.Workflow, *errors.PromptError) {
				workflowExecCallCount++
				assert.Equal(t, app1Workflow1, workflow)
				return nil, nil
			}

			result, err := fakeO.startInstructionWorkflowSelectionFlowFunc(fakeO, item1,
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
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return false, fmt.Errorf("failed to ask")
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionWorkflowFunc(fakeO, app1Workflow1)
			assert.False(t, result)
			assert.NotNil(t, err, "expected ask yes/no question to fail")
			assert.Equal(t, "failed to ask", err.GoError().Error())
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowExecAskYesNoQuestionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			fakeInput := input.CreateFakeUserInput()
			askYesNoCallCount := 0
			fakeInput.AskYesNoQuestionMock = func(question string) (bool, error) {
				askYesNoCallCount++
				return true, nil
			}

			fakeO.in = fakeInput
			result, err := fakeO.askBeforeRunningInstructionWorkflowFunc(fakeO, app1Workflow1)
			assert.True(t, result)
			assert.Nil(t, err, "expected ask yes/no question to succeed")
			assert.Equal(t, 1, askYesNoCallCount, "expected func to be called exactly once")
		})
	})
}

var InstructionWorkflowPromptFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionWorkflowsMock = func(appName string, workflows []*models.Workflow) (*models.Workflow, error) {
				promptCallCount++
				return nil, fmt.Errorf("failed to prompt for instruction workflow")
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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
			items := stubs.GenerateAnchorFolderItemsInfoTestData()
			item1 := stubs.GetAnchorFolderItemByName(items, stubs.AnchorFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			promptCallCount := 0
			fakePrompter := prompter.CreateFakePrompter()
			fakePrompter.PromptInstructionWorkflowsMock = func(appName string, workflows []*models.Workflow) (*models.Workflow, error) {
				promptCallCount++
				return app1Workflow1, nil
			}

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
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
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return false, errors.NewPromptError(fmt.Errorf("failed to ask"))
			}

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
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
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeO.runInstructionWorkflowFunc = func(
				o *selectOrchestrator,
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

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
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
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeO.runInstructionWorkflowFunc = func(
				o *selectOrchestrator,
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

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
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
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			askBeforeCallCount := 0
			fakeO.askBeforeRunningInstructionWorkflowFunc = func(
				o *selectOrchestrator,
				workflow *models.Workflow) (bool, *errors.PromptError) {

				askBeforeCallCount++
				return true, nil
			}

			runInstructionCallCount := 0
			fakeO.runInstructionWorkflowFunc = func(
				o *selectOrchestrator,
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

			result, err := fakeO.startInstructionWorkflowExecutionFlowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected instructions action execution to succeed")
			assert.Equal(t, 1, askBeforeCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runInstructionCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, wrapAfterExecCallCount, "expected func to be called exactly once")
			assert.Equal(t, result, app1Workflow1)
		})
	})
}

var RunInstructionWorkflowFailToRunActions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			runInstActionCallCount := 0
			fakeO.runInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) *errors.PromptError {
				runInstActionCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to run action"))
			}

			err := fakeO.runInstructionWorkflowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
			assert.NotNil(t, err)
			assert.Equal(t, "failed to run action", err.GoError().Error())
			assert.Equal(t, 1, runInstActionCallCount, "expected func to be called exactly once")
		})
	})
}

var RunInstructionWorkflowRunActionsSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := stubs.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.AnchorFolder1Name)
			runInstActionCallCount := 0
			fakeO.runInstructionActionFunc = func(o *selectOrchestrator, action *models.Action) *errors.PromptError {
				runInstActionCallCount++
				act := stubs.GetInstructionActionById(instRootTestData.Instructions, action.Id)
				assert.NotNil(t, act, "expected action to exist")
				return nil
			}

			err := fakeO.runInstructionWorkflowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
			assert.Nil(t, err)
			assert.Equal(t, 2, runInstActionCallCount, "expected func to be called the amount of actions")
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