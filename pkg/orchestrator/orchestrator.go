package orchestrator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/models"

	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"

	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

const (
	Identifier string = "orchestrator"
)

type Orchestrator interface {
	OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError)

	ExtractInstructions(
		app *models.ApplicationInfo,
		anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError)

	OrchestrateInstructionActionSelection(
		app *models.ApplicationInfo,
		actions []*models.Action) (*models.Action, *errors.PromptError)

	OrchestrateInstructionWorkflowSelection(
		app *models.ApplicationInfo,
		workflows []*models.Workflow) (*models.Workflow, *errors.PromptError)

	AskBeforeRunningInstructionAction(action *models.Action) (bool, *errors.PromptError)
	AskBeforeRunningInstructionWorkflow(workflow *models.Workflow) (bool, *errors.PromptError)

	RunInstructionAction(action *models.Action) *errors.PromptError
	RunInstructionWorkflow(
		workflow *models.Workflow,
		actions []*models.Action) *errors.PromptError

	WrapAfterActionRun() *errors.PromptError
}

type orchestratorImpl struct {
	Orchestrator
	prompter  prompter.Prompter
	locator   locator.Locator
	extractor extractor.Extractor
	parser    parser.Parser
	shell     shell.Shell
	input     input.UserInput
}

func New(
	pr prompter.Prompter,
	l locator.Locator,
	e extractor.Extractor,
	pa parser.Parser,
	s shell.Shell,
	in input.UserInput) Orchestrator {

	return &orchestratorImpl{
		prompter:  pr,
		locator:   l,
		extractor: e,
		parser:    pa,
		shell:     s,
		input:     in,
	}
}

func (o *orchestratorImpl) OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError) {
	apps := o.locator.Applications()
	if app, err := o.prompter.PromptApps(apps); err != nil {
		return nil, errors.New(err)
	} else {
		if app.Name == prompter.CancelActionName {
			return &models.ApplicationInfo{
				Name: prompter.CancelActionName,
			}, nil
		} else {
			return app, nil
		}
	}
}

func (o *orchestratorImpl) ExtractInstructions(
	app *models.ApplicationInfo,
	anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {

	path := app.InstructionsPath
	if instructionsRoot, err := o.extractor.ExtractInstructions(path, o.parser); err != nil {
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

func (o *orchestratorImpl) OrchestrateInstructionActionSelection(
	app *models.ApplicationInfo,
	actions []*models.Action) (*models.Action, *errors.PromptError) {

	item, err := o.prompter.PromptInstructionActions(app.Name, actions)
	if err != nil {
		return nil, errors.New(err)
	}
	return item, nil
}

func (o *orchestratorImpl) OrchestrateInstructionWorkflowSelection(
	app *models.ApplicationInfo,
	workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {

	item, err := o.prompter.PromptInstructionWorkflows(app.Name, workflows)
	if err != nil {
		return nil, errors.New(err)
	}
	return item, nil
}

func (o *orchestratorImpl) AskBeforeRunningInstructionAction(action *models.Action) (bool, *errors.PromptError) {
	question := prompter.GenerateRunInstructionMessage(action.Id, "action", action.Title)
	if res, err := o.input.AskYesNoQuestion(question); err != nil {
		return false, errors.New(err)
	} else {
		return res, nil
	}
}

func (o *orchestratorImpl) AskBeforeRunningInstructionWorkflow(workflow *models.Workflow) (bool, *errors.PromptError) {
	// TODO: Change description since it might be too long
	question := prompter.GenerateRunInstructionMessage(workflow.Id, "workflow", workflow.Description)
	if res, err := o.input.AskYesNoQuestion(question); err != nil {
		return false, errors.New(err)
	} else {
		return res, nil
	}
}

func (o *orchestratorImpl) RunInstructionAction(action *models.Action) *errors.PromptError {
	logger.Debugf("Running action: %v...", action.Id)
	scriptOutputPath, _ := logger.GetDefaultScriptOutputLogFilePath()

	if len(action.Script) > 0 && len(action.ScriptFile) > 0 {
		return errors.New(fmt.Errorf("script / scriptFile are mutual exclusive, please use either one"))
	} else if len(action.Script) == 0 && len(action.ScriptFile) == 0 {
		return errors.New(fmt.Errorf("missing script or scriptFile, nothing to run - skipping"))
	}

	if len(action.Script) > 0 {
		if err := o.shell.ExecuteWithOutputToFile(action.Script, scriptOutputPath); err != nil {
			return errors.New(err)
		}
	} else if len(action.ScriptFile) > 0 {
		if err := o.shell.ExecuteScriptFileWithOutputToFile(
			action.AnchorfilesRepoPath,
			action.ScriptFile,
			scriptOutputPath); err != nil {

			return errors.New(err)
		}
	}
	return nil
}

func (o *orchestratorImpl) RunInstructionWorkflow(
	workflow *models.Workflow,
	actions []*models.Action) *errors.PromptError {

	logger.Debugf("Running workflow: %v...", workflow.Id)
	for _, actionId := range workflow.ActionIds {
		action := models.GetInstructionActionById(actions, actionId)
		if promptErr := o.RunInstructionAction(action); promptErr != nil && !workflow.TolerateFailures {
			logger.Debugf("failed to run workflow and failures are not tolerable. "+
				"workflow: %s, action: %s", workflow.Id, action.Id)
			return promptErr
		}
	}
	return nil
}

func (o *orchestratorImpl) WrapAfterActionRun() *errors.PromptError {
	if inputErr := o.input.PressAnyKeyToContinue(); inputErr != nil {
		logger.Debugf("Failed to prompt user to press any key after instruction action run")
		return errors.New(inputErr)
	}
	if err := o.shell.ClearScreen(); err != nil {
		logger.Debugf("Failed to clear screen post instruction action run")
		return errors.New(err)
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
