package orchestrator

import (
	"github.com/ZachiNachshon/anchor/config"
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
		logger.Warningf("Failed to extract instructions from file. path: %s, error: %s", path, err.Error())
		return nil, errors.NewInstructionMissingError(err)
	} else {
		item, err := o.prompter.PromptInstructions(app.Name, instructions)
		if err != nil {
			return nil, errors.New(err)
		}
		return item, nil
	}
}

func (o *orchestratorImpl) AskBeforeRunningInstruction(item *models.InstructionItem) (bool, *errors.PromptError) {
	question := prompter.GenerateRunInstructionMessage(item.Id, item.Title)
	if res, err := o.input.AskYesNoQuestion(question); err != nil {
		return false, errors.New(err)
	} else {
		return res, nil
	}
}

func (o *orchestratorImpl) RunInstruction(item *models.InstructionItem, repoPath string) *errors.PromptError {
	logger.Debugf("Running: %v...", item.Id)
	scriptRunPath, _ := config.GetDefaultScriptRunLogFilePath()
	if err := o.shell.ExecuteScriptWithOutputToFile(repoPath, item.File, scriptRunPath); err != nil {
		return errors.New(err)
	} else {
		if inputErr := o.input.PressAnyKeyToContinue(); inputErr != nil {
			logger.Debugf("Failed to prompt user to press any key after instruction run")
			return errors.New(inputErr)
		}
		if err = o.shell.ClearScreen(); err != nil {
			logger.Debugf("Failed to clear screen post instruction run")
			return errors.New(err)
		}
		return nil
	}
}
