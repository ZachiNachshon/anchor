package _select

import (
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
	"github.com/manifoldco/promptui"
	"strings"
)

type DynamicSelectFunc func(ctx common.Context, o *selectOrchestrator) error

var DynamicSelect = func(ctx common.Context, o *selectOrchestrator) error {
	err := o.runner.PrepareFunc(o.runner, ctx)
	if err != nil {
		return err
	}
	err = o.prepareFunc(o, ctx)
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
	verboseFlag       bool
	commandFolderName string

	runner *runner.ActionRunnerOrchestrator

	prmpt prompter.Prompter
	l     locator.Locator
	s     shell.Shell
	in    input.UserInput
	prntr printer.Printer

	// --- CLI Command ---
	prepareFunc func(o *selectOrchestrator, ctx common.Context) error
	bannerFunc  func(o *selectOrchestrator)
	runFunc     func(o *selectOrchestrator, ctx common.Context) *errors.PromptError

	// --- Folder Items ---
	startCommandItemsSelectionFlowFunc func(o *selectOrchestrator, anchorfilesRepoPath string) *errors.PromptError
	promptCommandItemsSelectionFunc    func(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError)
	wrapAfterExecutionFunc             func(o *selectOrchestrator) *errors.PromptError

	// --- Action ---
	startInstructionActionSelectionFlowFunc func(
		o *selectOrchestrator,
		commandFolderItem *models.CommandFolderItemInfo,
		instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError)

	promptInstructionActionSelectionFunc func(
		o *selectOrchestrator,
		commandFolderItem *models.CommandFolderItemInfo,
		actions []*models.Action) (*models.Action, *errors.PromptError)

	startInstructionActionExecutionFlowFunc func(
		o *selectOrchestrator,
		globals *models.Globals,
		action *models.Action) (*models.Action, *errors.PromptError)

	askBeforeRunningInstructionActionFunc func(
		o *selectOrchestrator,
		globals *models.Globals,
		action *models.Action) (bool, *errors.PromptError)

	// --- Workflow ---
	startInstructionWorkflowSelectionFlowFunc func(
		o *selectOrchestrator,
		commandFolderItem *models.CommandFolderItemInfo,
		globals *models.Globals,
		workflows []*models.Workflow,
		actions []*models.Action) (*models.Workflow, *errors.PromptError)

	promptInstructionWorkflowSelectionFunc func(
		o *selectOrchestrator,
		commandFolderItem *models.CommandFolderItemInfo,
		workflows []*models.Workflow) (*models.Workflow, *errors.PromptError)

	startInstructionWorkflowExecutionFlowFunc func(
		o *selectOrchestrator,
		globals *models.Globals,
		workflow *models.Workflow,
		actions []*models.Action) (*models.Workflow, *errors.PromptError)

	askBeforeRunningInstructionWorkflowFunc func(
		o *selectOrchestrator,
		globals *models.Globals,
		workflow *models.Workflow) (bool, *errors.PromptError)
}

func NewOrchestrator(runner *runner.ActionRunnerOrchestrator, commandFolderName string) *selectOrchestrator {
	return &selectOrchestrator{
		commandFolderName: commandFolderName,
		verboseFlag:       false,
		runner:            runner,

		// --- CLI Command ---
		bannerFunc:  banner,
		prepareFunc: prepare,
		runFunc:     run,

		// --- Anchor Folder Item ---
		startCommandItemsSelectionFlowFunc: startCommandItemsSelectionFlow,
		promptCommandItemsSelectionFunc:    promptCommandItemsSelection,
		wrapAfterExecutionFunc:             wrapAfterExecution,

		// --- Action ---
		startInstructionActionSelectionFlowFunc: startInstructionActionSelectionFlow,
		promptInstructionActionSelectionFunc:    promptInstructionActionSelection,
		startInstructionActionExecutionFlowFunc: startInstructionActionExecutionFlow,
		askBeforeRunningInstructionActionFunc:   askBeforeRunningInstructionAction,

		// --- Workflow ---
		startInstructionWorkflowSelectionFlowFunc: startInstructionWorkflowSelectionFlow,
		promptInstructionWorkflowSelectionFunc:    promptInstructionWorkflowSelection,
		startInstructionWorkflowExecutionFlowFunc: startInstructionWorkflowExecutionFlow,
		askBeforeRunningInstructionWorkflowFunc:   askBeforeRunningInstructionWorkflow,
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

	if resolved, err := ctx.Registry().SafeGet(input.Identifier); err != nil {
		return err
	} else {
		o.in = resolved.(input.UserInput)
	}
	return nil
}

func banner(o *selectOrchestrator) {
	o.prntr.PrintAnchorBanner()
}

func run(o *selectOrchestrator, ctx common.Context) *errors.PromptError {
	return o.startCommandItemsSelectionFlowFunc(o, ctx.AnchorFilesPath())
}

func startCommandItemsSelectionFlow(o *selectOrchestrator, anchorfilesRepoPath string) *errors.PromptError {
	if commandFolder, promptErr := o.promptCommandItemsSelectionFunc(o); promptErr != nil {
		return promptErr
	} else if commandFolder.Name == prompter.CancelActionName {
		return nil
	} else {
		instRoot, promptError := o.runner.ExtractInstructionsFunc(o.runner, commandFolder, anchorfilesRepoPath)
		if promptError != nil {
			o.prntr.PrintMissingInstructions()
			_ = o.wrapAfterExecutionFunc(o)
			return o.startCommandItemsSelectionFlowFunc(o, anchorfilesRepoPath)
		}

		if instructionItem, promptErr := o.startInstructionActionSelectionFlowFunc(o, commandFolder, instRoot); promptErr != nil {
			if promptErr.Code() == errors.InstructionMissingError {
				return o.startCommandItemsSelectionFlowFunc(o, anchorfilesRepoPath)
			}
			return promptErr
		} else if instructionItem.Id == prompter.BackActionName {
			return o.startCommandItemsSelectionFlowFunc(o, anchorfilesRepoPath)
		}
		return nil
	}
}

func promptCommandItemsSelection(o *selectOrchestrator) (*models.CommandFolderItemInfo, *errors.PromptError) {
	folderItems := o.l.CommandFolderItems(o.commandFolderName)
	if app, err := o.prmpt.PromptCommandFolderItemSelection(folderItems); err != nil {
		return nil, errors.NewPromptError(err)
	} else {
		return app, nil
	}
}

func wrapAfterExecution(o *selectOrchestrator) *errors.PromptError {
	o.prntr.PrintEmptyLines(1)
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

func startInstructionActionSelectionFlow(
	o *selectOrchestrator,
	commandFolderItem *models.CommandFolderItemInfo,
	instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {

	appendInstructionActionsCustomOptions(instructionRoot.Instructions)
	actions := instructionRoot.Instructions.Actions

	if action, promptErr := o.promptInstructionActionSelectionFunc(o, commandFolderItem, actions); promptErr != nil {
		return nil, promptErr
	} else if action.Id == prompter.BackActionName {
		logger.Debugf("Selected to go back from instruction actions menu. id: %v", action.Id)
		return action, nil
	} else if action.Id == prompter.WorkflowsActionName {
		appendInstructionWorkflowsCustomOptions(instructionRoot.Instructions)
		workflows := instructionRoot.Instructions.Workflows
		logger.Debugf("Selected to prompt for instruction workflows menu. id: %v", action.Id)
		if _, promptErr = o.startInstructionWorkflowSelectionFlowFunc(o, commandFolderItem, instructionRoot.Globals, workflows, actions); promptErr != nil {
			return nil, promptErr
		} else {
			return o.startInstructionActionSelectionFlowFunc(o, commandFolderItem, instructionRoot)
		}
	} else {
		logger.Debugf("Selected instruction action to run. id: %v", action.Id)
		if _, promptErr = o.startInstructionActionExecutionFlowFunc(o, instructionRoot.Globals, action); promptErr != nil {
			return nil, promptErr
		} else {
			return o.startInstructionActionSelectionFlowFunc(o, commandFolderItem, instructionRoot)
		}
	}
}

func promptInstructionActionSelection(
	o *selectOrchestrator,
	commandFolderItem *models.CommandFolderItemInfo,
	actions []*models.Action) (*models.Action, *errors.PromptError) {

	item, err := o.prmpt.PromptInstructionActions(commandFolderItem.Name, actions)
	if err != nil {
		return nil, errors.NewPromptError(err)
	}
	return item, nil
}

func promptInstructionWorkflowSelection(
	o *selectOrchestrator,
	commandFolderItem *models.CommandFolderItemInfo,
	workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {

	item, err := o.prmpt.PromptInstructionWorkflows(commandFolderItem.Name, workflows)
	if err != nil {
		return nil, errors.NewPromptError(err)
	}
	return item, nil
}

func startInstructionWorkflowSelectionFlow(
	o *selectOrchestrator,
	commandFolderItem *models.CommandFolderItemInfo,
	globals *models.Globals,
	workflows []*models.Workflow,
	actions []*models.Action) (*models.Workflow, *errors.PromptError) {

	if workflow, promptError := o.promptInstructionWorkflowSelectionFunc(o, commandFolderItem, workflows); promptError != nil {
		return nil, promptError
	} else if workflow.Id == prompter.BackActionName {
		logger.Debugf("Selected to go back from instruction workflow menu. id: %v", workflow.Id)
		return workflow, nil
	} else {
		logger.Debugf("Selected instruction workflow to run. id: %v", workflow.Id)
		if _, promptErr := o.startInstructionWorkflowExecutionFlowFunc(o, globals, workflow, actions); promptErr != nil {
			return nil, promptErr
		} else {
			return o.startInstructionWorkflowSelectionFlowFunc(o, commandFolderItem, globals, workflows, actions)
		}
	}
}

func startInstructionActionExecutionFlow(
	o *selectOrchestrator,
	globals *models.Globals,
	action *models.Action) (*models.Action, *errors.PromptError) {

	if shouldRun, promptError := o.askBeforeRunningInstructionActionFunc(o, globals, action); promptError != nil {
		logger.Debugf("failed to ask before running an instruction action. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		// Do not break selection flow upon action failure, print warning and continue
		if promptErr := o.runner.RunInstructionActionFunc(o.runner, action); promptErr != nil && promptErr.Code() == errors.SchemaError {
			// Print only errors which aren't in direct relation to the script execution, these are handled differently
			o.prntr.PrintError(promptErr.GoError().Error())
		}
		if promptErr := o.wrapAfterExecutionFunc(o); promptErr != nil {
			return nil, promptErr
		}
	}
	return action, nil
}

func askBeforeRunningInstructionAction(
	o *selectOrchestrator,
	globals *models.Globals,
	action *models.Action) (bool, *errors.PromptError) {

	var question = ""
	instContext := getInstructionActionContext(globals, action)
	if instContext == models.ApplicationContext {
		question = prompter.GenerateRunInstructionMessage(action.Id, "action", action.Title)
	} else if instContext == models.KubernetesContext {
		question = prompter.GenerateKubernetesRunInstructionMessage(o.s, action.Id, "action", action.Title)
	}
	if res, err := o.in.AskYesNoQuestion(question); err != nil {
		return false, errors.NewPromptError(err)
	} else {
		return res, nil
	}
}

func askBeforeRunningInstructionWorkflow(
	o *selectOrchestrator,
	globals *models.Globals,
	workflow *models.Workflow) (bool, *errors.PromptError) {

	var question = ""
	instContext := getInstructionWorkflowContext(globals, workflow)
	if instContext == models.ApplicationContext {
		question = prompter.GenerateRunInstructionMessage(workflow.Id, "workflow", workflow.Title)
	} else if instContext == models.KubernetesContext {
		question = prompter.GenerateKubernetesRunInstructionMessage(o.s, workflow.Id, "workflow", workflow.Title)
	}
	if res, err := o.in.AskYesNoQuestion(question); err != nil {
		return false, errors.NewPromptError(err)
	} else {
		return res, nil
	}
}

func startInstructionWorkflowExecutionFlow(
	o *selectOrchestrator,
	globals *models.Globals,
	workflow *models.Workflow,
	actions []*models.Action) (*models.Workflow, *errors.PromptError) {

	if shouldRun, promptError := o.askBeforeRunningInstructionWorkflowFunc(o, globals, workflow); promptError != nil {
		logger.Debugf("failed to ask before running an instruction workflow. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		if promptErr := o.runner.RunInstructionWorkflowFunc(o.runner, workflow, actions); promptErr != nil {
			//return nil, promptErr
			// Do nothing, don't break the application flow, log error to file and prompt for any input to continue
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

	if ac := models.GetActionById(actions, prompter.BackActionName); ac != nil {
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

	if wf := models.GetWorkflowById(workflows, prompter.BackActionName); wf != nil {
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

func extractArgsFromScriptFile(scriptFile string) (string, []string) {
	split := strings.Split(scriptFile, " ")
	if len(split) > 1 {
		return split[0], split[1:]
	}
	return split[0], nil
}

func getInstructionActionContext(globals *models.Globals, action *models.Action) string {
	// action context always take precedence over global context
	if len(action.Context) > 0 {
		return action.Context
	} else if len(globals.Context) > 0 {
		return globals.Context
	}
	// default to application context
	return models.ApplicationContext
}

func getInstructionWorkflowContext(globals *models.Globals, workflow *models.Workflow) string {
	// workflow context always take precedence over global context
	if len(workflow.Context) > 0 {
		return workflow.Context
	} else if len(globals.Context) > 0 {
		return globals.Context
	}
	// default to application context
	return models.ApplicationContext
}
