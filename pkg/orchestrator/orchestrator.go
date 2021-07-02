package orchestrator

import (
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	extractor2 "github.com/ZachiNachshon/anchor/pkg/extractor"
	locator2 "github.com/ZachiNachshon/anchor/pkg/locator"
	parser2 "github.com/ZachiNachshon/anchor/pkg/parser"
	prompter2 "github.com/ZachiNachshon/anchor/pkg/prompter"
)

type orchestratorImpl struct {
	Orchestrator
	prompter  prompter2.Prompter
	locator   locator2.Locator
	extractor extractor2.Extractor
	parser    parser2.Parser
}

func New(
	pr prompter2.Prompter,
	l locator2.Locator,
	e extractor2.Extractor,
	pa parser2.Parser) Orchestrator {

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
		if app.Name == prompter2.CancelButtonName {
			return &models.PromptItem{
				Id: prompter2.CancelButtonName,
			}, nil
		} else {
			path := app.InstructionsPath
			if instructions, err := o.extractor.ExtractPromptItems(path, o.parser); err != nil {
				return nil, err
			} else {
				if item, err := o.prompter.PromptInstructions(instructions); err != nil {
					return nil, err
				} else {
					if item.Id == prompter2.BackButtonName {
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
