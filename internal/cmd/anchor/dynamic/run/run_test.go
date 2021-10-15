package run

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/runner"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RunActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete dynamic run method successfully",
			Func: CompleteDynamicRunMethodSuccessfully,
		},
		{
			Name: "fail dynamic run due to runner preparation",
			Func: FailDynamicRunDueToRunnerPreparation,
		},
		{
			Name: "fail dynamic run due to command folder item extraction",
			Func: FailDynamicRunDueToCommandFolderItemExtraction,
		},
		{
			Name: "fail dynamic run due to preparation",
			Func: FailDynamicRunDueToPreparation,
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
			Name: "command extraction: fail to find command folder items",
			Func: CommandExtractionFailToFindCommandFolderItems,
		},
		{
			Name: "command extraction: fail to find command items to run",
			Func: CommandExtractionFailToFindCommandItemToRun,
		},
		{
			Name: "command extraction: find command items to run successfully",
			Func: CommandExtractionFindCommandItemToRunSuccessfully,
		},
		{
			Name: "run action: fail to extract instructions",
			Func: RunActionFailToExtractInstructions,
		},
		{
			Name: "run action: fail to find instruction action by id",
			Func: RunActionFailToFindInstructionActionById,
		},
		{
			Name: "run action: fail to run instruction action",
			Func: RunActionFailToFindToRunInstructionAction,
		},
		{
			Name: "run action: run instruction action successfully",
			Func: RunActionRunInstructionActionSuccessfully,
		},
		{
			Name: "run workflow: fail to extract instructions",
			Func: RunWorkflowFailToExtractInstructions,
		},
		{
			Name: "run workflow: fail to find instruction workflow by id",
			Func: RunWorkflowFailToFindInstructionWorkflowById,
		},
		{
			Name: "run workflow: fail to run instruction workflow",
			Func: RunWorkflowFailToFindToRunInstructionWorkflow,
		},
		{
			Name: "run workflow: run instruction workflow successfully",
			Func: RunWorkflowRunInstructionWorkflowSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteDynamicRunMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			runnerPrepareCallCount := 0
			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				runnerPrepareCallCount++
				return nil
			}
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *runOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return nil
			}
			fakeO.extractCommandFolderItemFunc = func(o *runOrchestrator, commandFolderName string, commandFolderItemName string) (*models.CommandFolderItemInfo, error) {
				return nil, nil
			}
			runFuncCallCount := 0
			fakeO.activeRunFunc = func(o *runOrchestrator, ctx common.Context, cmdFolderItem *models.CommandFolderItemInfo, identifier string) error {
				assert.Equal(t, stubs.CommandFolder1Item1Action1Id, identifier)
				runFuncCallCount++
				return nil
			}
			err := DynamicRun(ctx, fakeO, stubs.CommandFolder1Item1Action1Id)
			assert.Nil(t, err, "expected not to fail")
			assert.Equal(t, 1, runnerPrepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runFuncCallCount, "expected func to be called exactly once")
		})
	})
}

var FailDynamicRunDueToRunnerPreparation = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			runnerPrepareCallCount := 0
			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				runnerPrepareCallCount++
				return fmt.Errorf("failed to prepare runner")
			}
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)
			err := DynamicRun(ctx, fakeO, "")
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare runner", err.Error())
			assert.Equal(t, 1, runnerPrepareCallCount, "expected func to be called exactly once")
		})
	})
}

var FailDynamicRunDueToCommandFolderItemExtraction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			runnerPrepareCallCount := 0
			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				runnerPrepareCallCount++
				return nil
			}
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *runOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return nil
			}
			extractCommandFolderItemCallCount := 0
			fakeO.extractCommandFolderItemFunc = func(o *runOrchestrator, commandFolderName string, commandFolderItemName string) (*models.CommandFolderItemInfo, error) {
				extractCommandFolderItemCallCount++
				return nil, fmt.Errorf("failed to extract command folder item")
			}
			err := DynamicRun(ctx, fakeO, "")
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to extract command folder item", err.Error())
			assert.Equal(t, 1, runnerPrepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, extractCommandFolderItemCallCount, "expected func to be called exactly once")
		})
	})
}

var FailDynamicRunDueToPreparation = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cmdFolderName := stubs.CommandFolder1Name
			fakeRunner := runner.NewOrchestrator(cmdFolderName)
			runnerPrepareCallCount := 0
			fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
				runnerPrepareCallCount++
				return nil
			}
			fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *runOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return fmt.Errorf("failed to prepare dynamic run")
			}
			fakeO.extractCommandFolderItemFunc = func(o *runOrchestrator, commandFolderName string, commandFolderItemName string) (*models.CommandFolderItemInfo, error) {
				return nil, nil
			}
			err := DynamicRun(ctx, fakeO, "")
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare dynamic run", err.Error())
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
		})
	})
}

var PrepareRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()

		fakeLocator := locator.CreateFakeLocator()
		reg.Set(locator.Identifier, fakeLocator)

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		cmdFolderName := stubs.CommandFolder1Name
		fakeRunner := runner.NewOrchestrator(cmdFolderName)
		fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)
		err := fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
		assert.NotNil(t, fakeO.prntr)
		assert.NotNil(t, fakeO.l)
	})
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		cmdFolderName := stubs.CommandFolder1Name
		fakeRunner := runner.NewOrchestrator(cmdFolderName)
		fakeRunner.PrepareFunc = func(o *runner.ActionRunnerOrchestrator, ctx common.Context) error {
			return nil
		}
		fakeO := NewOrchestrator(fakeRunner, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)

		err := fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", locator.Identifier))

		fakeLocator := locator.CreateFakeLocator()
		reg.Set(locator.Identifier, fakeLocator)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", printer.Identifier))

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
	})
}

var CommandExtractionFailToFindCommandFolderItems = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		cmdFolderName := stubs.CommandFolder1Name
		cmdFolderItemName := stubs.CommandFolder1Item1Name

		fakeLocator := locator.CreateFakeLocator()
		cmdFolderItemsCallCount := 0
		fakeLocator.CommandFolderItemsMock = func(commandFolderName string) []*models.CommandFolderItemInfo {
			cmdFolderItemsCallCount++
			return nil
		}

		fakePrinter := printer.CreateFakePrinter()
		printMissingCmdCallCount := 0
		fakePrinter.PrintMissingCommandMock = func(name string) {
			printMissingCmdCallCount++
		}

		fakeO := NewOrchestrator(nil, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)
		fakeO.l = fakeLocator
		fakeO.prntr = fakePrinter

		item, err := extractCommandFolderItem(fakeO, cmdFolderName, cmdFolderItemName)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Nil(t, item, "expected extraction to fail")
		assert.Equal(t, 1, cmdFolderItemsCallCount)
		assert.Equal(t, 1, printMissingCmdCallCount)
	})
}

var CommandExtractionFailToFindCommandItemToRun = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeLocator := locator.CreateFakeLocator()
		cmdFolderItemsCallCount := 0
		fakeLocator.CommandFolderItemsMock = func(commandFolderName string) []*models.CommandFolderItemInfo {
			cmdFolderItemsCallCount++
			return stubs.GenerateCommandFolderItemsInfoTestData()
		}

		fakeO := NewOrchestrator(nil, "", "")
		fakeO.l = fakeLocator

		missingItemName := "no-where-to-be-found"
		fakeO.commandFolderItemName = missingItemName
		item, err := extractCommandFolderItem(fakeO, stubs.CommandFolder1Name, missingItemName)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Equal(t, err.Error(), fmt.Sprintf("cannot identify dynamic command folder item: %s", missingItemName))
		assert.Nil(t, item, "expected extraction to fail")
		assert.Equal(t, 1, cmdFolderItemsCallCount)
	})
}

var CommandExtractionFindCommandItemToRunSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeLocator := locator.CreateFakeLocator()
		cmdFolderItemsCallCount := 0
		fakeLocator.CommandFolderItemsMock = func(commandFolderName string) []*models.CommandFolderItemInfo {
			cmdFolderItemsCallCount++
			return stubs.GenerateCommandFolderItemsInfoTestData()
		}

		fakeO := NewOrchestrator(nil, "", "")
		fakeO.l = fakeLocator

		item, err := extractCommandFolderItem(fakeO, stubs.CommandFolder1Name, stubs.CommandFolder1Item1Name)
		assert.Nil(t, err, "expected extraction to succeed")
		assert.NotNil(t, item, "expected extraction to succeed")
		assert.Equal(t, 1, cmdFolderItemsCallCount)
	})
}

var RunActionFailToExtractInstructions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		cmdFolderItemsInfoData := stubs.GenerateCommandFolderItemsInfoTestData()
		cmdItem := stubs.GetCommandFolderItemByName(cmdFolderItemsInfoData, stubs.CommandFolder1Item1Name)

		cmdFolderName := stubs.CommandFolder1Name
		fakeRunner := runner.NewOrchestrator(cmdFolderName)
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return nil, errors.NewPromptError(fmt.Errorf("failed to extract instructions"))
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		err := runAction(fakeO, ctx, cmdItem, stubs.CommandFolder1Item1Action1Id)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Equal(t, 1, extractInstCallCount)
	})
}

var RunActionFailToFindInstructionActionById = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		instTestData := stubs.GenerateInstructionsTestData()
		fakeRunner := runner.NewOrchestrator("")
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return instTestData, nil
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		missingActionId := "no-where-to-be-found"
		err := runAction(fakeO, ctx, nil, missingActionId)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Contains(t, err.Error(), "Cannot identify action by id")
		assert.Equal(t, 1, extractInstCallCount)
	})
}

var RunActionFailToFindToRunInstructionAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		instTestData := stubs.GenerateInstructionsTestData()
		fakeRunner := runner.NewOrchestrator("")
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return instTestData, nil
		}
		runInstCallCount := 0
		fakeRunner.RunInstructionActionFunc = func(o *runner.ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
			runInstCallCount++
			return errors.NewPromptError(fmt.Errorf("failed to run instruction action"))
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		err := runAction(fakeO, ctx, nil, stubs.CommandFolder1Item1Action1Id)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Equal(t, err.Error(), "failed to run instruction action")
		assert.Equal(t, 1, extractInstCallCount)
		assert.Equal(t, 1, runInstCallCount)
	})
}

var RunActionRunInstructionActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		instTestData := stubs.GenerateInstructionsTestData()
		fakeRunner := runner.NewOrchestrator("")
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return instTestData, nil
		}
		runInstCallCount := 0
		fakeRunner.RunInstructionActionFunc = func(o *runner.ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
			runInstCallCount++
			return nil
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		err := runAction(fakeO, ctx, nil, stubs.CommandFolder1Item1Action1Id)
		assert.Nil(t, err, "expected extraction to succeed")
		assert.Equal(t, 1, extractInstCallCount)
		assert.Equal(t, 1, runInstCallCount)
	})
}

var RunWorkflowFailToExtractInstructions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		cmdFolderItemsInfoData := stubs.GenerateCommandFolderItemsInfoTestData()
		cmdItem := stubs.GetCommandFolderItemByName(cmdFolderItemsInfoData, stubs.CommandFolder1Item1Name)

		cmdFolderName := stubs.CommandFolder1Name
		fakeRunner := runner.NewOrchestrator(cmdFolderName)
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return nil, errors.NewPromptError(fmt.Errorf("failed to extract instructions"))
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		err := runWorkflow(fakeO, ctx, cmdItem, stubs.CommandFolder1Item1Workflow1Id)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Equal(t, 1, extractInstCallCount)
	})
}

var RunWorkflowFailToFindInstructionWorkflowById = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		instTestData := stubs.GenerateInstructionsTestData()
		fakeRunner := runner.NewOrchestrator("")
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return instTestData, nil
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		missingWorkflowId := "no-where-to-be-found"
		err := runWorkflow(fakeO, ctx, nil, missingWorkflowId)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Contains(t, err.Error(), "Cannot identify workflow by id")
		assert.Equal(t, 1, extractInstCallCount)
	})
}

var RunWorkflowFailToFindToRunInstructionWorkflow = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		instTestData := stubs.GenerateInstructionsTestData()
		fakeRunner := runner.NewOrchestrator("")
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return instTestData, nil
		}
		runInstCallCount := 0
		fakeRunner.RunInstructionWorkflowFunc = func(o *runner.ActionRunnerOrchestrator, workflow *models.Workflow, actions []*models.Action) *errors.PromptError {
			runInstCallCount++
			return errors.NewPromptError(fmt.Errorf("failed to run instruction workflow"))
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		err := runWorkflow(fakeO, ctx, nil, stubs.CommandFolder1Item1Workflow1Id)
		assert.NotNil(t, err, "expected extraction to fail")
		assert.Equal(t, err.Error(), "failed to run instruction workflow")
		assert.Equal(t, 1, extractInstCallCount)
		assert.Equal(t, 1, runInstCallCount)
	})
}

var RunWorkflowRunInstructionWorkflowSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		instTestData := stubs.GenerateInstructionsTestData()
		fakeRunner := runner.NewOrchestrator("")
		extractInstCallCount := 0
		fakeRunner.ExtractInstructionsFunc = func(o *runner.ActionRunnerOrchestrator, commandFolderItem *models.CommandFolderItemInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {
			extractInstCallCount++
			return instTestData, nil
		}
		runInstCallCount := 0
		fakeRunner.RunInstructionWorkflowFunc = func(o *runner.ActionRunnerOrchestrator, workflow *models.Workflow, actions []*models.Action) *errors.PromptError {
			runInstCallCount++
			return nil
		}
		fakeO := NewOrchestrator(nil, "", "")
		fakeO.runner = fakeRunner

		err := runWorkflow(fakeO, ctx, nil, stubs.CommandFolder1Item1Workflow1Id)
		assert.Nil(t, err, "expected extraction to succeed")
		assert.Equal(t, 1, extractInstCallCount)
		assert.Equal(t, 1, runInstCallCount)
	})
}
