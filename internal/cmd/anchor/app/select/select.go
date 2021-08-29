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
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"

	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/manifoldco/promptui"
)

type AppSelectFunc func(ctx common.Context, o *selectOrchestrator) error

var AppSelect = func(ctx common.Context, o *selectOrchestrator) error {
	err := o.prepareFunc(o, ctx)
	if err != nil {
		return err
	}
	o.bannerFunc(o)
	promptErr := o.runFunc(o, ctx)
	if promptErr != nil {
		return managePromptError(promptErr)
	}
	return nil
}

type selectOrchestrator struct {
	prmpt prompter.Prompter
	l     locator.Locator
	e     extractor.Extractor
	prsr  parser.Parser
	s     shell.Shell
	in    input.UserInput
	prntr printer.Printer

	// --- CLI Command ---
	prepareFunc func(o *selectOrchestrator, ctx common.Context) error
	bannerFunc  func(o *selectOrchestrator)
	runFunc     func(o *selectOrchestrator, ctx common.Context) *errors.PromptError

	// --- Application ---
	startApplicationSelectionFlowFunc func(o *selectOrchestrator, anchorfilesRepoPath string) *errors.PromptError
	promptApplicationSelectionFunc    func(o *selectOrchestrator) (*models.ApplicationInfo, *errors.PromptError)
	wrapAfterExecutionFunc            func(o *selectOrchestrator) *errors.PromptError

	// --- Action ---
	startInstructionActionSelectionFlowFunc func(
		o *selectOrchestrator,
		app *models.ApplicationInfo,
		instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError)

	promptInstructionActionSelectionFunc func(
		o *selectOrchestrator,
		app *models.ApplicationInfo,
		actions []*models.Action) (*models.Action, *errors.PromptError)

	extractInstructionsFunc func(
		o *selectOrchestrator,
		app *models.ApplicationInfo,
		anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError)

	startInstructionActionExecutionFlowFunc func(
		o *selectOrchestrator,
		action *models.Action) (*models.Action, *errors.PromptError)

	askBeforeRunningInstructionActionFunc func(
		o *selectOrchestrator,
		action *models.Action) (bool, *errors.PromptError)

	runInstructionActionFunc func(
		o *selectOrchestrator,
		action *models.Action) *errors.PromptError

	// --- Workflow ---
	startInstructionWorkflowSelectionFlowFunc func(
		o *selectOrchestrator,
		app *models.ApplicationInfo,
		workflows []*models.Workflow,
		actions []*models.Action) (*models.Workflow, *errors.PromptError)

	promptInstructionWorkflowSelectionFunc func(
		o *selectOrchestrator,
		app *models.ApplicationInfo,
		workflows []*models.Workflow) (*models.Workflow, *errors.PromptError)

	startInstructionWorkflowExecutionFlowFunc func(
		o *selectOrchestrator,
		workflow *models.Workflow,
		actions []*models.Action) (*models.Workflow, *errors.PromptError)

	askBeforeRunningInstructionWorkflowFunc func(
		o *selectOrchestrator,
		workflow *models.Workflow) (bool, *errors.PromptError)

	runInstructionWorkflowFunc func(
		o *selectOrchestrator,
		workflow *models.Workflow,
		actions []*models.Action) *errors.PromptError
}

func NewOrchestrator() *selectOrchestrator {
	return &selectOrchestrator{
		// --- CLI Command ---
		bannerFunc:  banner,
		prepareFunc: prepare,
		runFunc:     run,

		// --- Application ---
		startApplicationSelectionFlowFunc: startApplicationSelectionFlow,
		promptApplicationSelectionFunc:    promptApplicationSelection,
		wrapAfterExecutionFunc:            wrapAfterExecution,

		// --- Action ---
		startInstructionActionSelectionFlowFunc: startInstructionActionSelectionFlow,
		promptInstructionActionSelectionFunc:    promptInstructionActionSelection,
		extractInstructionsFunc:                 extractInstructions,
		startInstructionActionExecutionFlowFunc: startInstructionActionExecutionFlow,
		askBeforeRunningInstructionActionFunc:   askBeforeRunningInstructionAction,
		runInstructionActionFunc:                runInstructionAction,

		// --- Workflow ---
		startInstructionWorkflowSelectionFlowFunc: startInstructionWorkflowSelectionFlow,
		promptInstructionWorkflowSelectionFunc:    promptInstructionWorkflowSelection,
		startInstructionWorkflowExecutionFlowFunc: startInstructionWorkflowExecutionFlow,
		askBeforeRunningInstructionWorkflowFunc:   askBeforeRunningInstructionWorkflow,
		runInstructionWorkflowFunc:                runInstructionWorkflow,
	}
}

func prepare(o *selectOrchestrator, ctx common.Context) error {
	if resolved, err := ctx.Registry().SafeGet(locator.Identifier); err != nil {
		return err
	} else {
		o.l = resolved.(locator.Locator)
	}

	if resolved, err := ctx.Registry().SafeGet(prompter.Identifier); err != nil {
		return err
	} else {
		o.prmpt = resolved.(prompter.Prompter)
	}

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
	return nil
}

func banner(o *selectOrchestrator) {
	o.prntr.PrintAnchorBanner()
}

func run(o *selectOrchestrator, ctx common.Context) *errors.PromptError {
	return o.startApplicationSelectionFlowFunc(o, ctx.AnchorFilesPath())
}

func startApplicationSelectionFlow(o *selectOrchestrator, anchorfilesRepoPath string) *errors.PromptError {
	if app, promptErr := o.promptApplicationSelectionFunc(o); promptErr != nil {
		return promptErr
	} else if app.Name == prompter.CancelActionName {
		return nil
	} else {
		instRoot, promptError := o.extractInstructionsFunc(o, app, anchorfilesRepoPath)
		if promptError != nil {
			return o.startApplicationSelectionFlowFunc(o, anchorfilesRepoPath)
		}

		if instructionItem, promptErr := o.startInstructionActionSelectionFlowFunc(o, app, instRoot); promptErr != nil {
			if promptErr.Code() == errors.InstructionMissingError {
				return o.startApplicationSelectionFlowFunc(o, anchorfilesRepoPath)
			}
			return promptErr
		} else if instructionItem.Id == prompter.BackActionName {
			return o.startApplicationSelectionFlowFunc(o, anchorfilesRepoPath)
		}
		return nil
	}
}

func promptApplicationSelection(o *selectOrchestrator) (*models.ApplicationInfo, *errors.PromptError) {
	apps := o.l.Applications()
	if app, err := o.prmpt.PromptApps(apps); err != nil {
		return nil, errors.NewPromptError(err)
	} else {
		return app, nil
	}
}

func wrapAfterExecution(o *selectOrchestrator) *errors.PromptError {
	if inputErr := o.in.PressAnyKeyToContinue(); inputErr != nil {
		logger.Debugf("Failed to prompt user to press any key after instruction action run")
		return errors.NewPromptError(inputErr)
	}
	if err := o.s.ClearScreen(); err != nil {
		logger.Debugf("Failed to clear screen post instruction action run")
		return errors.NewPromptError(err)
	}
	return nil
}

func extractInstructions(
	o *selectOrchestrator,
	app *models.ApplicationInfo,
	anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {

	path := app.InstructionsPath
	if instructionsRoot, err := o.e.ExtractInstructions(path, o.prsr); err != nil {
		logger.Warningf("Failed to extract instructions from file. error: %s", err.Error())
		return nil, errors.NewInstructionMissingError(err)
	} else {
		if instructionsRoot == nil || instructionsRoot.Instructions == nil {
			// Perform the same prompt selection flow (back action etc..) on empty instructions due to invalid schema
			instructionsRoot = models.EmptyInstructionsRoot()
		} else {
			enrichActionsWithWorkingDirectoryCanonicalPath(anchorfilesRepoPath, instructionsRoot.Instructions.Actions)
		}
		return instructionsRoot, nil
	}
}

func runInstructionAction(o *selectOrchestrator, action *models.Action) *errors.PromptError {
	logger.Debugf("Running action: %v...", action.Id)
	scriptOutputPath, _ := logger.GetDefaultScriptOutputLogFilePath()

	if len(action.Script) > 0 && len(action.ScriptFile) > 0 {
		return errors.NewPromptError(fmt.Errorf("script / scriptFile are mutual exclusive, please use either one"))
	} else if len(action.Script) == 0 && len(action.ScriptFile) == 0 {
		return errors.NewPromptError(fmt.Errorf("missing script or scriptFile, nothing to run - skipping"))
	}

	if len(action.Script) > 0 {
		if err := o.s.ExecuteWithOutputToFile(action.Script, scriptOutputPath); err != nil {
			return errors.NewPromptError(err)
		}
	} else if len(action.ScriptFile) > 0 {
		if err := o.s.ExecuteScriptFileWithOutputToFile(
			action.AnchorfilesRepoPath,
			action.ScriptFile,
			scriptOutputPath); err != nil {

			return errors.NewPromptError(err)
		}
	}
	return nil
}

func startInstructionActionSelectionFlow(
	o *selectOrchestrator,
	app *models.ApplicationInfo,
	instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {

	appendInstructionActionsCustomOptions(instructionRoot.Instructions)
	actions := instructionRoot.Instructions.Actions

	if action, promptErr := o.promptInstructionActionSelectionFunc(o, app, actions); promptErr != nil {
		return nil, promptErr
	} else if action.Id == prompter.BackActionName {
		logger.Debugf("Selected to go back from instruction actions menu. id: %v", action.Id)
		return action, nil
	} else if action.Id == prompter.WorkflowsActionName {
		appendInstructionWorkflowsCustomOptions(instructionRoot.Instructions)
		workflows := instructionRoot.Instructions.Workflows
		logger.Debugf("Selected to prompt for instruction workflows menu. id: %v", action.Id)
		if _, promptErr := o.startInstructionWorkflowSelectionFlowFunc(o, app, workflows, actions); promptErr != nil {
			return nil, promptErr
		} else {
			return o.startInstructionActionSelectionFlowFunc(o, app, instructionRoot)
		}
	} else {
		logger.Debugf("Selected instruction action to run. id: %v", action.Id)
		if _, promptErr = o.startInstructionActionExecutionFlowFunc(o, action); promptErr != nil {
			return nil, promptErr
		} else {
			return o.startInstructionActionSelectionFlowFunc(o, app, instructionRoot)
		}
	}
}

func promptInstructionActionSelection(
	o *selectOrchestrator,
	app *models.ApplicationInfo,
	actions []*models.Action) (*models.Action, *errors.PromptError) {

	item, err := o.prmpt.PromptInstructionActions(app.Name, actions)
	if err != nil {
		return nil, errors.NewPromptError(err)
	}
	return item, nil
}

func promptInstructionWorkflowSelection(
	o *selectOrchestrator,
	app *models.ApplicationInfo,
	workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {

	item, err := o.prmpt.PromptInstructionWorkflows(app.Name, workflows)
	if err != nil {
		return nil, errors.NewPromptError(err)
	}
	return item, nil
}

func startInstructionWorkflowSelectionFlow(
	o *selectOrchestrator,
	app *models.ApplicationInfo,
	workflows []*models.Workflow,
	actions []*models.Action) (*models.Workflow, *errors.PromptError) {

	if workflow, promptError := o.promptInstructionWorkflowSelectionFunc(o, app, workflows); promptError != nil {
		return nil, promptError
	} else if workflow.Id == prompter.BackActionName {
		logger.Debugf("Selected to go back from instruction workflow menu. id: %v", workflow.Id)
		return workflow, nil
	} else {
		logger.Debugf("Selected instruction workflow to run. id: %v", workflow.Id)
		if _, promptErr := o.startInstructionWorkflowExecutionFlowFunc(o, workflow, actions); promptErr != nil {
			return nil, promptErr
		} else {
			return o.startInstructionWorkflowSelectionFlowFunc(o, app, workflows, actions)
		}
	}
}

func startInstructionActionExecutionFlow(
	o *selectOrchestrator,
	action *models.Action) (*models.Action, *errors.PromptError) {
	if shouldRun, promptError := o.askBeforeRunningInstructionActionFunc(o, action); promptError != nil {
		logger.Debugf("failed to ask before running an instruction action. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		if promptErr := o.runInstructionActionFunc(o, action); promptErr != nil {
			return nil, promptErr
		}
		if promptErr := o.wrapAfterExecutionFunc(o); promptErr != nil {
			return nil, promptErr
		}
	}
	return action, nil
}

func askBeforeRunningInstructionAction(
	o *selectOrchestrator,
	action *models.Action) (bool, *errors.PromptError) {
	question := prompter.GenerateRunInstructionMessage(action.Id, "action", action.Title)
	if res, err := o.in.AskYesNoQuestion(question); err != nil {
		return false, errors.NewPromptError(err)
	} else {
		return res, nil
	}
}

func askBeforeRunningInstructionWorkflow(
	o *selectOrchestrator,
	workflow *models.Workflow) (bool, *errors.PromptError) {
	// TODO: Change description since it might be too long
	question := prompter.GenerateRunInstructionMessage(workflow.Id, "workflow", workflow.Description)
	if res, err := o.in.AskYesNoQuestion(question); err != nil {
		return false, errors.NewPromptError(err)
	} else {
		return res, nil
	}
}

func runInstructionWorkflow(
	o *selectOrchestrator,
	workflow *models.Workflow,
	actions []*models.Action) *errors.PromptError {

	logger.Debugf("Running workflow: %v...", workflow.Id)
	for _, actionId := range workflow.ActionIds {
		action := models.GetInstructionActionById(actions, actionId)
		if promptErr := o.runInstructionActionFunc(o, action); promptErr != nil && !workflow.TolerateFailures {
			logger.Debugf("failed to run workflow and failures are not tolerable. "+
				"workflow: %s, action: %s", workflow.Id, action.Id)
			return promptErr
		}
	}
	return nil
}

func startInstructionWorkflowExecutionFlow(
	o *selectOrchestrator,
	workflow *models.Workflow,
	actions []*models.Action) (*models.Workflow, *errors.PromptError) {

	if shouldRun, promptError := o.askBeforeRunningInstructionWorkflowFunc(o, workflow); promptError != nil {
		logger.Debugf("failed to ask before running an instruction workflow. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		if promptErr := o.runInstructionWorkflowFunc(o, workflow, actions); promptErr != nil {
			return nil, promptErr
		}
		if promptErr := o.wrapAfterExecutionFunc(o); promptErr != nil {
			return nil, promptErr
		}
	}
	return workflow, nil
}

func managePromptError(promptErr *errors.PromptError) error {
	err := promptErr.GoError()
	if err == nil {
		logger.Debug("Prompt error returned but does not contain an inner Go error")
		return err
	}
	if err == promptui.ErrInterrupt {
		logger.Debug("exit due to keyboard interrupt")
		return nil
	} else {
		logger.Debug(err.Error())
		return err
	}
}

func appendInstructionActionsCustomOptions(instructions *models.Instructions) {
	actions := instructions.Actions

	if ac := models.GetInstructionActionById(actions, prompter.BackActionName); ac != nil {
		return
	}

	enrichedActionsList := make([]*models.Action, 0, len(actions)+2)
	backAction := &models.Action{
		Id: prompter.BackActionName,
	}
	enrichedActionsList = append(enrichedActionsList, backAction)

	if len(instructions.Workflows) > 0 {
		workflowsAction := &models.Action{
			Id: prompter.WorkflowsActionName,
		}
		enrichedActionsList = append(enrichedActionsList, workflowsAction)
	}

	enrichedActionsList = append(enrichedActionsList, actions...)
	instructions.Actions = enrichedActionsList
}

func appendInstructionWorkflowsCustomOptions(instructions *models.Instructions) {
	workflows := instructions.Workflows

	if wf := models.GetInstructionWorkflowById(workflows, prompter.BackActionName); wf != nil {
		return
	}

	enrichedWorkflowsList := make([]*models.Workflow, 0, len(workflows)+1)
	backAction := &models.Workflow{
		Id: prompter.BackActionName,
	}
	enrichedWorkflowsList = append(enrichedWorkflowsList, backAction)
	enrichedWorkflowsList = append(enrichedWorkflowsList, workflows...)
	instructions.Workflows = enrichedWorkflowsList
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
