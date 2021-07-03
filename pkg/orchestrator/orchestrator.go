package orchestrator

import (
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
)

type orchestratorImpl struct {
	Orchestrator
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

func (o *orchestratorImpl) OrchestrateAppInstructionSelection() (*models.PromptItem, error) {
	apps := o.locator.Applications()
	if app, err := o.prompter.PromptApps(apps); err != nil {
		return nil, err
	} else {
		logger.Debugf("Selected application. app: %v", app)
		if app.Name == prompter.CancelButtonName {
			return &models.PromptItem{
				Id: prompter.CancelButtonName,
			}, nil
		} else {
			path := app.InstructionsPath
			if instructions, err := o.extractor.ExtractPromptItems(path, o.parser); err != nil {
				return nil, err
			} else {
				if item, err := o.prompter.PromptInstructions(instructions); err != nil {
					return nil, err
				} else {
					if item.Id == prompter.BackButtonName {
						return o.OrchestrateAppInstructionSelection()
					} else {
						logger.Debugf("Selected instruction to run. id: %v", item.Id)
						return item, nil
					}
				}
			}
		}
	}
}
