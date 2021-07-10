package orchestrator

import (
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

type orchestratorImpl struct {
	Orchestrator
	// TODO: these could get removed and added as a method argument
	prompter  prompter.Prompter
	locator   locator.Locator
	extractor extractor.Extractor
	parser    parser.Parser
}

func New(
	pr prompter.Prompter,
	l locator.Locator,
	e extractor.Extractor,
	pa parser.Parser) Orchestrator {

	return &orchestratorImpl{
		prompter:  pr,
		locator:   l,
		extractor: e,
		parser:    pa,
	}
}

func (o *orchestratorImpl) OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError) {
	apps := o.locator.Applications()
	if app, err := o.prompter.PromptApps(apps); err != nil {
		return nil, errors.New(err)
	} else {
		logger.Debugf("Selected application. app: %v", app)
		if app.Name == prompter.CancelButtonName {
			return &models.ApplicationInfo{
				Name: prompter.CancelButtonName,
			}, nil
		} else {
			return app, nil
		}
	}
}

func (o *orchestratorImpl) OrchestrateInstructionSelection(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
	path := app.InstructionsPath
	if instructions, err := o.extractor.ExtractInstructions(path, o.parser); err != nil {
		logger.Warningf("Missing instructions file. path: %s", path)
		return nil, errors.NewInstructionMissingError(err)
	} else {
		item, err := o.prompter.PromptInstructions(app.Name, instructions)
		if err != nil {
			return nil, errors.New(err)
		}
		return item, nil
	}
}

func (o *orchestratorImpl) AskBeforeRunningInstruction(item *models.InstructionItem, in input.UserInput) (bool, *errors.PromptError) {
	question := prompter.GenerateRunInstructionMessage(item.Id, item.Title)
	if res, err := in.AskYesNoQuestion(question); err != nil {
		return false, errors.New(err)
	} else {
		return res, nil
	}
}

func (o *orchestratorImpl) RunInstruction(item *models.InstructionItem, s shell.Shell) *errors.PromptError {
	// TODO: log to file script output
	logger.Infof("Running: %v...", item.Id)
	return nil
}
