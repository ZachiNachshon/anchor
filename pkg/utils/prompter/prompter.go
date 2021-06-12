package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

const (
	backButtonName   = "back"
	cancelButtonName = "cancel"
)

type prompterImpl struct {
	Prompter
}

func New() Prompter {
	return &prompterImpl{}
}

func (p *prompterImpl) PromptApps(l locator.Locator) (*models.AppContent, error) {
	setSearchAppPrompt()
	appsArr := l.Applications()
	appsSelector := PreparePromptAppsItems(appsArr)
	i, _, err := appsSelector.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and run apps prompt selector. error: %s", err.Error())
	}
	logger.Debugf("selected app value. index: %d, name: %s", i+1, appsArr[i].Name)
	return appsArr[i], nil
}

func (p *prompterImpl) PromptInstructions(instructions *models.Instructions) (*models.PromptItem, error) {
	setSearchInstructionsPrompt()
	instSelector := PreparePromptInstructionsItems(instructions)
	i, _, err := instSelector.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and run instruction prompt selector. error: %s", err.Error())
	}
	logger.Debugf("selected instruction value. index: %d, name: %s", i+1, instructions.Items[i].Id)
	return instructions.Items[i], nil
}
