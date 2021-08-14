package _select

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/models"

	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/pkg/printer"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/manifoldco/promptui"
)

type AppSelectFunc func(ctx common.Context) error

var AppSelect = func(ctx common.Context) error {
	var o orchestrator.Orchestrator
	if resolved, err := ctx.Registry().SafeGet(orchestrator.Identifier); err != nil {
		return err
	} else {
		o = resolved.(orchestrator.Orchestrator)
	}

	var p printer.Printer
	if resolved, err := ctx.Registry().SafeGet(printer.Identifier); err != nil {
		return err
	} else {
		p = resolved.(printer.Printer)
	}

	p.PrintAnchorBanner()
	promptErr := runApplicationSelectionFlow(o, ctx.AnchorFilesPath())
	if promptErr != nil {
		return managePromptError(promptErr)
	}
	return nil
}

func runApplicationSelectionFlow(o orchestrator.Orchestrator, anchorfilesRepoPath string) *errors.PromptError {
	if app, promptErr := o.OrchestrateApplicationSelection(); promptErr != nil {
		return promptErr
	} else if app.Name == prompter.CancelActionName {
		return nil
	} else {
		instRoot, promptError := o.ExtractInstructions(app, anchorfilesRepoPath)
		if promptError != nil {
			return runApplicationSelectionFlow(o, anchorfilesRepoPath)
		}

		if instructionItem, promptErr := runInstructionActionSelectionFlow(o, app, instRoot); promptErr != nil {
			if promptErr.Code() == errors.InstructionMissingError {
				return runApplicationSelectionFlow(o, anchorfilesRepoPath)
			}
			return promptErr
		} else if instructionItem.Id == prompter.BackActionName {
			return runApplicationSelectionFlow(o, anchorfilesRepoPath)
		}
		return nil
	}
}

func runInstructionActionSelectionFlow(
	o orchestrator.Orchestrator,
	app *models.ApplicationInfo,
	instructionRoot *models.InstructionsRoot) (*models.Action, *errors.PromptError) {

	appendInstructionActionsCustomOptions(instructionRoot.Instructions)
	actions := instructionRoot.Instructions.Actions

	if action, promptErr := o.OrchestrateInstructionActionSelection(app, actions); promptErr != nil {
		return nil, promptErr
	} else if action.Id == prompter.BackActionName {
		logger.Debugf("Selected to go back from instruction actions menu. id: %v", action.Id)
		return action, nil
	} else if action.Id == prompter.WorkflowsActionName {
		appendInstructionWorkflowsCustomOptions(instructionRoot.Instructions)
		workflows := instructionRoot.Instructions.Workflows
		logger.Debugf("Selected to prompt for instruction workflows menu. id: %v", action.Id)
		if _, promptErr := runInstructionWorkflowSelectionFlow(o, app, workflows, actions); promptErr != nil {
			return nil, promptErr
		} else {
			return runInstructionActionSelectionFlow(o, app, instructionRoot)
		}
	} else {
		logger.Debugf("Selected instruction action to run. id: %v", action.Id)
		if _, promptErr := runInstructionActionExecutionFlow(o, action); promptErr != nil {
			return nil, promptErr
		} else {
			return runInstructionActionSelectionFlow(o, app, instructionRoot)
		}
	}
}

func runInstructionWorkflowSelectionFlow(
	o orchestrator.Orchestrator,
	app *models.ApplicationInfo,
	workflows []*models.Workflow,
	actions []*models.Action) (*models.Workflow, *errors.PromptError) {

	if workflow, promptError := o.OrchestrateInstructionWorkflowSelection(app, workflows); promptError != nil {
		return nil, promptError
	} else if workflow.Id == prompter.BackActionName {
		logger.Debugf("Selected to go back from instruction workflow menu. id: %v", workflow.Id)
		return workflow, nil
	} else {
		logger.Debugf("Selected instruction workflow to run. id: %v", workflow.Id)
		if _, promptErr := runInstructionWorkflowExecutionFlow(o, workflow, actions); promptErr != nil {
			return nil, promptErr
		} else {
			return runInstructionWorkflowSelectionFlow(o, app, workflows, actions)
		}
	}
}

func runInstructionActionExecutionFlow(
	o orchestrator.Orchestrator,
	action *models.Action) (*models.Action, *errors.PromptError) {

	if shouldRun, promptError := o.AskBeforeRunningInstructionAction(action); promptError != nil {
		logger.Debugf("failed to ask before running an instruction action. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		if promptErr := o.RunInstructionAction(action); promptErr != nil {
			return nil, promptErr
		}
		if promptErr := o.WrapAfterActionRun(); promptErr != nil {
			return nil, promptErr
		}
	}
	return action, nil
}

func runInstructionWorkflowExecutionFlow(
	o orchestrator.Orchestrator,
	workflow *models.Workflow,
	actions []*models.Action) (*models.Workflow, *errors.PromptError) {

	if shouldRun, promptError := o.AskBeforeRunningInstructionWorkflow(workflow); promptError != nil {
		logger.Debugf("failed to ask before running an instruction workflow. error: %s", promptError.GoError().Error())
		return nil, promptError
	} else if shouldRun {
		if promptErr := o.RunInstructionWorkflow(workflow, actions); promptErr != nil {
			return nil, promptErr
		}
		if promptErr := o.WrapAfterActionRun(); promptErr != nil {
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
