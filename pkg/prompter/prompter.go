package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/screenbuf"
)

const (
	BackActionName      = "back"
	WorkflowsActionName = "workflows..."
	CancelActionName    = "cancel"
	selectorEmoji       = "\U0001F449"
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

	logger.Debugf("Selected app value. index: %d, name: %s", i, appsEnhanced[i].Name)
	return appsEnhanced[i], nil
}

func (p *prompterImpl) PromptInstructions(appName string, instructionsRoot *models.InstructionsRoot) (*models.Action, error) {
	instructions := instructionsRoot.Instructions
	setSearchInstructionsPrompt(appName)
	instSelector := preparePromptInstructionsActions(instructions)

	i, _, err := instSelector.Run()
	if err != nil {
		return nil, err
	}

	logger.Debugf("selected instruction value. index: %d, name: %s", i, instructions.Actions[i].Id)
	return instructions.Actions[i], nil
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
