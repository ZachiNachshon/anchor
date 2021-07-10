package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/screenbuf"
)

const (
	BackButtonName   = "back"
	CancelButtonName = "cancel"
	selectorEmoji    = "\U0001F449"
)

type prompterImpl struct {
	Prompter
}

func New() Prompter {
	return &prompterImpl{}
}

func (p *prompterImpl) PromptApps(apps []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
	setSearchAppPrompt()
	appsSelector := preparePromptAppsItems(apps)
	appsEnhanced := appsSelector.Items.([]*models.ApplicationInfo)

	i, _, err := appsSelector.Run()
	if err != nil {
		return nil, err
	}

	logger.Debugf("selected app value. index: %d, name: %s", i, appsEnhanced[i].Name)
	return appsEnhanced[i], nil
}

func (p *prompterImpl) PromptInstructions(appName string, instructions *models.Instructions) (*models.InstructionItem, error) {
	setSearchInstructionsPrompt(appName)
	instSelector := preparePromptInstructionsItems(instructions)

	i, _, err := instSelector.Run()
	if err != nil {
		return nil, err
	}

	logger.Debugf("selected instruction value. index: %d, name: %s", i, instructions.Items[i].Id)
	return instructions.Items[i], nil
}

func ClearScreen(selector promptui.Select) {
	buf := screenbuf.New(selector.Stdout)
	err := buf.Clear()
	if err != nil {
		logger.Warningf("failed to clear the screen. error: %s", err.Error())
	}
}

func createPaddingString(length int) string {
	// This is an example output "\"%-23s\""
	return "\"%-" + fmt.Sprintf("%v", length) + "s\""
}
