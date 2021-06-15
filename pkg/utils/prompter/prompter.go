package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
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

func (p *prompterImpl) PromptApps(apps []*models.AppContent) (*models.AppContent, error) {
	setSearchAppPrompt()
	appsSelector := preparePromptAppsItems(apps)
	appsEnhanced := appsSelector.Items.([]*models.AppContent)

	i, _, err := appsSelector.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and run apps prompt selector. error: %s", err.Error())
	}
	logger.Debugf("selected app value. index: %d, name: %s", i, appsEnhanced[i].Name)
	return appsEnhanced[i], nil
}

func (p *prompterImpl) PromptInstructions(instructions *models.Instructions) (*models.PromptItem, error) {
	setSearchInstructionsPrompt()
	instSelector := preparePromptInstructionsItems(instructions)
	i, _, err := instSelector.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and run instruction prompt selector. error: %s", err.Error())
	}
	logger.Debugf("selected instruction value. index: %d, name: %s", i, instructions.Items[i].Id)
	return instructions.Items[i], nil
}
