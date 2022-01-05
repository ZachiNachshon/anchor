package runner

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RunnerShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail resolving registry components",
			Func: FailResolvingRegistryComponents,
		},
		{
			Name: "fill instructions globals with defaults if missing from schema",
			Func: FillInstructionsGlobalsWithDefaultsIfMissingFromSchema,
		},
		{
			Name: "action exec: fail to exec script",
			Func: ActionExecFailToExecScript,
		},
		{
			Name: "action exec: exec script successfully",
			Func: ActionExecExecScriptSuccessfully,
		},
		{
			Name: "action exec: fail to exec script file",
			Func: ActionExecFailToExecScriptFile,
		},
		{
			Name: "action exec: exec script file successfully",
			Func: ActionExecExecScriptFileSuccessfully,
		},
		{
			Name: "action exec: no op when nothing to exec",
			Func: ActionExecNoOpWhenNothingToExec,
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
			Name: "extract instructions: fail when missing",
			Func: ExtractInstructionsFailWhenMissing,
		},
		{
			Name: "extract instructions: return empty on invalid schema",
			Func: ExtractInstructionsReturnEmptyOnInvalidSchema,
		},
		{
			Name: "extract instructions: return enriched actions",
			Func: ExtractInstructionsReturnEnrichedActions,
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
			Name: "run instruction workflow: fail to run actions",
			Func: RunInstructionWorkflowFailToRunActions,
		},
		{
			Name: "run instruction workflow: run actions successfully",
			Func: RunInstructionWorkflowRunActionsSuccessfully,
		},
		{
			Name: "extract multiple args from action script file attribute",
			Func: ExtractMultipleArgsFromActionScriptFileAttribute,
		},
		{
			Name: "extract single command from action script file attribute",
			Func: ExtractSingleCommandFromActionScriptFileAttribute,
		},
		{
			Name: "enrich actions with working dir canonical path",
			Func: EnrichActionsWithWorkingDirCanonicalPath,
		},
	}
	harness.RunTests(t, tests)
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		cmdFolderName := stubs.CommandFolder1Name
		fakeO := NewOrchestrator(cmdFolderName)

		err := fakeO.PrepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", extractor.Identifier))

		fakeExtractor := extractor.CreateFakeExtractor()
		reg.Set(extractor.Identifier, fakeExtractor)

		err = fakeO.PrepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", parser.Identifier))

		fakeParser := parser.CreateFakeParser()
		reg.Set(parser.Identifier, fakeParser)

		err = fakeO.PrepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", printer.Identifier))

		fakePrinter := printer.CreateFakePrinter()
		reg.Set(printer.Identifier, fakePrinter)

		err = fakeO.PrepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", shell.Identifier))

		fakeShell := shell.CreateFakeShell()
		reg.Set(shell.Identifier, fakeShell)

		err = fakeO.PrepareFunc(fakeO, ctx)
		assert.Nil(t, err)
	})
}

var FillInstructionsGlobalsWithDefaultsIfMissingFromSchema = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	fillInstructionGlobals(instRootTestData)
	assert.NotNil(t, instRootTestData.Globals, "expected non nil globals object")
	assert.Empty(t, instRootTestData.Globals, "expected empty globals object")
}

var ActionExecFailToExecScript = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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

var ActionExecExecScriptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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

var ActionExecFailToExecScriptFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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
			fakeShell.ExecuteScriptFileSilentlyWithOutputToFileMock = func(relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return fmt.Errorf("fail to execute")
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.s = fakeShell
			fakeO.prntr = fakePrinter

			err := executeInstructionAction(fakeO, action1, "")
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, 1, spinCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, execCallCount, "expected to be called exactly once")
			assert.Equal(t, 1, stopOnFailureCallCount, "expected to be called exactly once")
		})
	})
}

var ActionExecExecScriptFileSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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
			fakeShell.ExecuteScriptFileSilentlyWithOutputToFileMock = func(relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return nil
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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

var ActionExecNoOpWhenNothingToExec = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

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

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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
			fakeShell.ExecuteScriptFileWithOutputToFileMock = func(relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return fmt.Errorf("fail to execute")
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
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
			fakeShell.ExecuteScriptFileWithOutputToFileMock = func(relativeScriptPath string, outputFilePath string, args ...string) error {
				execCallCount++
				assert.Equal(t, "/path/to/script", relativeScriptPath)
				return nil
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

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

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.prntr = fakePrinter

			err := executeInstructionActionVerbose(fakeO, action1, "")
			assert.Nil(t, err)
			assert.Equal(t, 0, startCallCount, "expected no calls to be made")
		})
	})
}

var ExtractInstructionsReturnEnrichedActions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.HarnessAnchorfilesTestRepo(ctx)
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)
			instRootTestData := stubs.GenerateInstructionsTestData()

			extractInstructionsCallCount := 0
			fakeExtractor := extractor.CreateFakeExtractor()
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractInstructionsCallCount++
				assert.Equal(t, item1.InstructionsPath, instructionsPath)
				return instRootTestData, nil
			}

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.e = fakeExtractor

			result, err := fakeO.ExtractInstructionsFunc(fakeO, item1, ctx.AnchorFilesPath())
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			action1.Script = "some script"
			action1.ScriptFile = "/path/to/script"
			err := fakeO.RunInstructionActionFunc(fakeO, action1)
			assert.NotNil(t, err, "expected instructions execution to fail")
			assert.Contains(t, err.GoError().Error(), "script / scriptFile are mutual exclusive")
		})
	})
}

var RunInstructionActionFailOnMissingScriptToExec = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)

			action1.Script = ""
			action1.ScriptFile = ""
			err := fakeO.RunInstructionActionFunc(fakeO, action1)
			assert.NotNil(t, err, "expected instructions execution to fail")
			assert.Contains(t, err.GoError().Error(), "missing script or scriptFile")
		})
	})
}

var RunInstructionActionFailToExecuteAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.verboseFlag = true

			execActionVerboseCallCount := 0
			fakeO.executeInstructionActionVerboseFunc = func(o *ActionRunnerOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionVerboseCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to execute action"))
			}

			action1.Script = "some script"
			action1.ScriptFile = ""
			err := fakeO.RunInstructionActionFunc(fakeO, action1)
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.verboseFlag = false

			execActionCallCount := 0
			fakeO.executeInstructionActionFunc = func(o *ActionRunnerOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to execute action"))
			}

			action1.Script = ""
			action1.ScriptFile = "/path/to/script"
			err := fakeO.RunInstructionActionFunc(fakeO, action1)
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
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.verboseFlag = true
			execActionVerboseCallCount := 0
			fakeO.executeInstructionActionVerboseFunc = func(o *ActionRunnerOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionVerboseCallCount++
				return nil
			}

			action1.Script = "some script with verbose from flag"
			action1.ScriptFile = ""
			err := fakeO.RunInstructionActionFunc(fakeO, action1)
			assert.Nil(t, err)
			assert.Equal(t, 1, execActionVerboseCallCount, "expected to be called exactly once")
		})
	})
}

var RunInstructionActionRunActionWithForcedVerboseFromSchema = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)
			action1.ShowOutput = true

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.verboseFlag = false
			execActionVerboseCallCount := 0
			fakeO.executeInstructionActionVerboseFunc = func(o *ActionRunnerOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionVerboseCallCount++
				return nil
			}

			action1.Script = ""
			action1.ScriptFile = "some script with forced verbose from schema"
			err := fakeO.RunInstructionActionFunc(fakeO, action1)
			assert.Nil(t, err)
			assert.Equal(t, 1, execActionVerboseCallCount, "expected to be called exactly once")
		})
	})
}

var RunInstructionActionRunActionInteractive = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			action1 := models.GetInstructionActionById(instRootTestData.Instructions, stubs.CommandFolder1Item1Action1Id)

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.verboseFlag = false
			execActionCallCount := 0
			fakeO.executeInstructionActionFunc = func(o *ActionRunnerOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
				execActionCallCount++
				return nil
			}
			action1.Script = "some script"
			action1.ScriptFile = ""
			err := fakeO.RunInstructionActionFunc(fakeO, action1)
			assert.Nil(t, err)
			assert.Equal(t, 1, execActionCallCount, "expected to be called exactly once")
		})
	})
}

var RunInstructionWorkflowFailToRunActions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			instRootTestData := stubs.GenerateInstructionsTestData()
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			runInstActionCallCount := 0
			fakeO.RunInstructionActionFunc = func(o *ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
				runInstActionCallCount++
				return errors.NewPromptError(fmt.Errorf("failed to run action"))
			}

			err := fakeO.RunInstructionWorkflowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
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
			app1Workflow1 := models.GetInstructionWorkflowById(instRootTestData.Instructions, stubs.CommandFolder1Item1Workflow1Id)

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			runInstActionCallCount := 0
			fakeO.RunInstructionActionFunc = func(o *ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
				runInstActionCallCount++
				act := models.GetInstructionActionById(instRootTestData.Instructions, action.Id)
				assert.NotNil(t, act, "expected action to exist")
				return nil
			}

			err := fakeO.RunInstructionWorkflowFunc(fakeO, app1Workflow1, instRootTestData.Instructions.Actions)
			assert.Nil(t, err)
			assert.Equal(t, 2, runInstActionCallCount, "expected func to be called the amount of actions")
		})
	})
}

var ExtractInstructionsFailWhenMissing = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)

			extractInstructionsCallCount := 0
			fakeExtractor := extractor.CreateFakeExtractor()
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractInstructionsCallCount++
				return nil, fmt.Errorf("failed to extract instructions")
			}

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.e = fakeExtractor

			result, err := fakeO.ExtractInstructionsFunc(fakeO, item1, ctx.AnchorFilesPath())
			assert.Nil(t, result)
			assert.NotNil(t, err, "expected instructions extraction to fail")
			assert.Equal(t, errors.InstructionMissingError, err.Code())
			assert.Equal(t, 1, extractInstructionsCallCount)
		})
	})
}

var ExtractInstructionsReturnEmptyOnInvalidSchema = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			items := stubs.GenerateCommandFolderItemsInfoTestData()
			item1 := stubs.GetCommandFolderItemByName(items, stubs.CommandFolder1Item1Name)

			extractInstructionsCallCount := 0
			fakeExtractor := extractor.CreateFakeExtractor()
			fakeExtractor.ExtractInstructionsMock = func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
				extractInstructionsCallCount++
				assert.Equal(t, item1.InstructionsPath, instructionsPath)
				return nil, nil
			}

			fakeO := NewOrchestrator(stubs.CommandFolder1Name)
			fakeO.e = fakeExtractor

			result, err := fakeO.ExtractInstructionsFunc(fakeO, item1, ctx.AnchorFilesPath())
			assert.NotNil(t, result)
			assert.Nil(t, err, "expected instructions extraction to succeed")
			assert.Equal(t, 1, extractInstructionsCallCount)
			assert.Equal(t, models.EmptyInstructionsRoot(), result)
		})
	})
}

var ExtractMultipleArgsFromActionScriptFileAttribute = func(t *testing.T) {
	scriptFile := "/path/to/script ${PWD} --custom=flag"
	cmd, args := extractArgsFromScriptFile(scriptFile)
	assert.NotEmpty(t, cmd)
	assert.Equal(t, "/path/to/script", cmd)
	assert.Len(t, args, 2)
	assert.Equal(t, "${PWD}", args[0])
	assert.Equal(t, "--custom=flag", args[1])
}

var ExtractSingleCommandFromActionScriptFileAttribute = func(t *testing.T) {
	scriptFile := "/path/to/script"
	cmd, args := extractArgsFromScriptFile(scriptFile)
	assert.NotEmpty(t, cmd)
	assert.Equal(t, "/path/to/script", cmd)
	assert.Nil(t, args)
}

var EnrichActionsWithWorkingDirCanonicalPath = func(t *testing.T) {
	anchorfilesRepoPath := "/path/to/anchorfiles/repo"
	instData := stubs.GenerateInstructionsTestData()
	enrichActionsWithWorkingDirectoryCanonicalPath(anchorfilesRepoPath, nil)
	for _, action := range instData.Instructions.Actions {
		assert.Empty(t, action.AnchorfilesRepoPath)
	}

	enrichActionsWithWorkingDirectoryCanonicalPath(anchorfilesRepoPath, instData.Instructions.Actions)
	for _, action := range instData.Instructions.Actions {
		assert.NotEmpty(t, action.AnchorfilesRepoPath)
		assert.Equal(t, anchorfilesRepoPath, action.AnchorfilesRepoPath)
	}
}
