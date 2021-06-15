package prompter

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
)

type orchestratorImpl struct {
	Orchestrator
	prompter  Prompter
	locator   locator.Locator
	extractor extractor.Extractor
	parser    parser.Parser
}

func NewOrchestrator(ctx common.Context) (Orchestrator, error) {
	registry := ctx.Registry()
	prompterInjection, err := FromRegistry(registry)
	if err != nil {
		return nil, err
	}

	locatorInjection, err := locator.FromRegistry(registry)
	if err != nil {
		return nil, err
	}

	extractorInjection, err := extractor.FromRegistry(registry)
	if err != nil {
		return nil, err
	}

	parserInjection, err := parser.FromRegistry(registry)
	if err != nil {
		return nil, err
	}

	return &orchestratorImpl{
		prompter:  prompterInjection,
		locator:   locatorInjection,
		extractor: extractorInjection,
		parser:    parserInjection,
	}, nil
}

func (o *orchestratorImpl) OrchestrateAppInstructionSelection() (*models.PromptItem, error) {
	apps := o.locator.Applications()
	if app, err := o.prompter.PromptApps(apps); err != nil {
		return nil, err
	} else {
		logger.Debugf("Selected application. app: %v", app)
		if app.Name == cancelButtonName {
			return &models.PromptItem{
				Id: cancelButtonName,
			}, nil
		} else {
			path := app.InstructionsPath
			if instructions, err := o.extractor.ExtractPromptItems(path, o.parser); err != nil {
				return nil, err
			} else {
				if item, err := o.prompter.PromptInstructions(instructions); err != nil {
					return nil, err
				} else {
					if item.Id == backButtonName {
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
