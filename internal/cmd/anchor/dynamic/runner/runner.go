package runner

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"strings"
)

type ActionRunnerOrchestrator struct {
	commandFolderName string
	VerboseFlag       bool

	e     extractor.Extractor
	prsr  parser.Parser
	s     shell.Shell
	prntr printer.Printer

	// --- Internal ---
	PrepareFunc func(o *ActionRunnerOrchestrator, ctx common.Context) error

	ExtractInstructionsFunc func(
		o *ActionRunnerOrchestrator,
		commandFolderItem *models.CommandFolderItemInfo,
		anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError)

	RunInstructionActionFunc func(
		o *ActionRunnerOrchestrator,
		action *models.Action) *errors.PromptError

	RunInstructionWorkflowFunc func(
		o *ActionRunnerOrchestrator,
		workflow *models.Workflow,
		actions []*models.Action) *errors.PromptError

	executeInstructionActionFunc func(
		o *ActionRunnerOrchestrator,
		action *models.Action,
		scriptOutputPath string) *errors.PromptError

	executeInstructionActionVerboseFunc func(
		o *ActionRunnerOrchestrator,
		action *models.Action,
		scriptOutputPath string) *errors.PromptError
}

func NewOrchestrator(commandFolderName string) *ActionRunnerOrchestrator {
	return &ActionRunnerOrchestrator{
		commandFolderName: commandFolderName,
		VerboseFlag:       false,

		// --- Internal ---
		PrepareFunc: prepare,

		// --- Action ---
		ExtractInstructionsFunc:             extractInstructions,
		RunInstructionActionFunc:            runInstructionAction,
		RunInstructionWorkflowFunc:          runInstructionWorkflow,
		executeInstructionActionFunc:        executeInstructionAction,
		executeInstructionActionVerboseFunc: executeInstructionActionVerbose,
	}
}

func prepare(o *ActionRunnerOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(extractor.Identifier); err != nil {
		return err
	} else {
		o.e = resolved.(extractor.Extractor)
	}

	if resolved, err := ctx.Registry().SafeGet(parser.Identifier); err != nil {
		return err
	} else {
		o.prsr = resolved.(parser.Parser)
	}

	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		o.prntr = resolved.(printer.Printer)
	}

	if resolved, err := ctx.Registry().SafeGet(shell.Identifier); err != nil {
		return err
	} else {
		o.s = resolved.(shell.Shell)
	}
	return nil
}

func extractArgsFromScriptFile(scriptFile string) (string, []string) {
	split := strings.SplitN(scriptFile, " ", 2)
	if len(split) == 1 {
		return split[0], nil
	}

	cmd := split[0]
	argsStr := split[1]
	if !strings.ContainsAny(argsStr, "\"") {
		splitSpace := strings.Split(argsStr, " ")
		return cmd, splitSpace
	} else {
		var argsResult []string
		for len(argsStr) > 0 {
			argsStr = strings.Trim(argsStr, " ")
			var splitInner []string
			if strings.Index(argsStr, "\"") == 0 {
				argsStr = strings.TrimPrefix(argsStr, "\"")
				// Cut until the closing quote and continue with the rest of the args string
				splitInner = strings.SplitN(argsStr, "\"", 2)
				extractedArg := splitInner[0]
				extractedArg = strings.TrimSuffix(extractedArg, "\"")
				quotedArg := fmt.Sprintf("\"%s\"", extractedArg)
				argsResult = append(argsResult, quotedArg)
			} else {
				// Cut the first argument and continue with the rest of the args string
				splitInner = strings.SplitN(argsStr, " ", 2)
				extractedArg := splitInner[0]
				argsResult = append(argsResult, extractedArg)
			}

			if splitInner != nil && len(splitInner) > 1 {
				argsStr = splitInner[1]
			} else {
				argsStr = ""
			}
		}
		return cmd, argsResult
	}
}

func extractInstructions(
	o *ActionRunnerOrchestrator,
	commandFolderItemInfo *models.CommandFolderItemInfo,
	anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {

	path := commandFolderItemInfo.InstructionsPath
	if instructionsRoot, err := o.e.ExtractInstructions(path, o.prsr); err != nil {
		logger.Warningf("Failed to extract instructions from file. error: %s", err.Error())
		return nil, errors.NewInstructionMissingError(err)
	} else {
		if instructionsRoot == nil || instructionsRoot.Instructions == nil {
			// Perform the same prompt selection flow (back action etc..) on empty instructions due to invalid schema
			instructionsRoot = models.EmptyInstructionsRoot()
		} else {
			enrichActionsWithWorkingDirectoryCanonicalPath(anchorfilesRepoPath, instructionsRoot.Instructions.Actions)
			enrichActionsWithRunCommands(o.commandFolderName, commandFolderItemInfo.Name, instructionsRoot.Instructions.Actions)
			enrichWorkflowsWithRunCommands(o.commandFolderName, commandFolderItemInfo.Name, instructionsRoot.Instructions.Workflows)
			fillInstructionGlobals(instructionsRoot)
		}
		return instructionsRoot, nil
	}
}

func runInstructionAction(o *ActionRunnerOrchestrator, action *models.Action) *errors.PromptError {
	logger.Debugf("Running action: %v...", action.Id)
	scriptOutputPath, _ := logger.GetDefaultScriptOutputLogFilePath()

	if len(action.Script) > 0 && len(action.ScriptFile) > 0 {
		return errors.NewSchemaError(fmt.Errorf("script / scriptFile are mutual exclusive, please use either one"))
	} else if len(action.Script) == 0 && len(action.ScriptFile) == 0 {
		return errors.NewSchemaError(fmt.Errorf("missing script or scriptFile, nothing to run - skipping"))
	}

	if o.VerboseFlag || action.ShowOutput {
		return o.executeInstructionActionVerboseFunc(o, action, scriptOutputPath)
	} else {
		return o.executeInstructionActionFunc(o, action, scriptOutputPath)
	}
}

func runInstructionWorkflow(
	o *ActionRunnerOrchestrator,
	workflow *models.Workflow,
	actions []*models.Action) *errors.PromptError {

	logger.Debugf("Running workflow: %v...", workflow.Id)
	for _, actionId := range workflow.ActionIds {
		action := models.GetActionById(actions, actionId)
		// TODO: continue skip if action is missing since stale action ids could be added to the workflow

		if promptErr := o.RunInstructionActionFunc(o, action); promptErr != nil && !workflow.TolerateFailures {
			logger.Errorf("failed to run workflow and failures are not tolerable. "+
				"workflow: %s, action: %s", workflow.Id, action.Id)
			return promptErr
		}
	}
	return nil
}

func executeInstructionAction(o *ActionRunnerOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
	spnr := o.prntr.PrepareRunActionSpinner(action.Id, scriptOutputPath)

	if len(action.Script) > 0 {
		spnr.Spin()
		if err := o.s.ExecuteSilentlyWithOutputToFile(action.Script, scriptOutputPath); err != nil {
			logger.Errorf("failed to run action. id: %s, source: script, error: %s", action.Id, err.Error())
			spnr.StopOnFailure(err)
			return errors.NewPromptError(err)
		}
		spnr.StopOnSuccess()
	} else if len(action.ScriptFile) > 0 {
		filePath, args := extractArgsFromScriptFile(action.ScriptFile)
		spnr.Spin()
		if err := o.s.ExecuteScriptFileSilentlyWithOutputToFile(
			filePath,
			scriptOutputPath,
			args...); err != nil {
			logger.Errorf("failed to run action. id: %s, source: script file, args: %v, error: %s", action.Id, args, err.Error())
			spnr.StopOnFailure(err)
			return errors.NewPromptError(err)
		}
		spnr.StopOnSuccess()
	}
	return nil
}

func executeInstructionActionVerbose(o *ActionRunnerOrchestrator, action *models.Action, scriptOutputPath string) *errors.PromptError {
	plainer := o.prntr.PrepareRunActionPlainer(action.Id)

	if len(action.Script) > 0 {
		plainer.Start()
		if err := o.s.ExecuteWithOutputToFile(action.Script, scriptOutputPath); err != nil {
			plainer.StopOnFailure(err)
			logger.Errorf("failed to run action. id: %s, source: script, error: %s", action.Id, err.Error())
			return errors.NewPromptError(err)
		}
		plainer.StopOnSuccess()
	} else if len(action.ScriptFile) > 0 {
		filePath, args := extractArgsFromScriptFile(action.ScriptFile)
		plainer.Start()
		//if err := o.s.ExecuteScriptFileWithOutputToFile(
		if err := o.s.ExecuteScriptFileWithOutputToFile(
			filePath,
			scriptOutputPath,
			args...); err != nil {

			logger.Errorf("failed to run action. id: %s, source: script file, args: %s, error: %s", action.Id, args, err.Error())
			plainer.StopOnFailure(err)
			return errors.NewPromptError(err)
		}
		plainer.StopOnSuccess()
	}
	return nil
}

func enrichActionsWithWorkingDirectoryCanonicalPath(anchorfilesRepoPath string, actions []*models.Action) {
	if actions == nil {
		return
	}
	for _, action := range actions {
		if action.ScriptFile != "" {
			action.AnchorfilesRepoPath = anchorfilesRepoPath
		}
	}
}

func enrichActionsWithRunCommands(cmdFolderName string, cmdItemName string, actions []*models.Action) {
	cliCmd := globals.CLIRootCommandName
	if actions == nil {
		return
	}
	for _, action := range actions {
		// Add the direct run command within every action description panel for ease of use
		// anchor CMD run INSTRUCTION --action=ACTION_ID
		action.RunCommand = fmt.Sprintf("%s%s%s %s%s%s run %s%s%s --action=%s%s%s",
			colors.Green, cliCmd, colors.Reset,
			colors.Yellow, cmdFolderName, colors.Reset,
			colors.Purple, cmdItemName, colors.Reset,
			colors.Cyan, action.Id, colors.Reset)
	}
}

func enrichWorkflowsWithRunCommands(cmdFolderName string, cmdItemName string, workflows []*models.Workflow) {
	cliCmd := globals.CLIRootCommandName
	if workflows == nil {
		return
	}
	for _, workflow := range workflows {
		// Add the direct run command within every workflow description panel for ease of use
		// anchor CMD run INSTRUCTION --workflow=WORKFLOW_ID
		workflow.RunCommand = fmt.Sprintf("%s%s%s %s%s%s run %s%s%s --workflow=%s%s%s",
			colors.Green, cliCmd, colors.Reset,
			colors.Yellow, cmdFolderName, colors.Reset,
			colors.Purple, cmdItemName, colors.Reset,
			colors.Cyan, workflow.Id, colors.Reset)
	}
}

func fillInstructionGlobals(instRoot *models.InstructionsRoot) {
	if instRoot.Globals == nil {
		instRoot.Globals = models.EmptyGlobals()
	}
}
