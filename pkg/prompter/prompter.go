package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/screenbuf"
)

const (
	BackActionName          = "back"
	WorkflowsActionName     = "workflows..."
	CancelActionName        = "cancel"
	selectorEmoji           = "\U0001F449"
	CheckMarkEmoji          = "\U00002705"
	CrossMarkEmoji          = "\U0000274C"
	selectorEmojiCharLength = 3
)

type prompterImpl struct {
	Prompter
}

func New() Prompter {
	return &prompterImpl{}
}

func (p *prompterImpl) PromptConfigContext(cfgContexts []*config.Context) (*config.Context, error) {
	setSearchConfigContextPrompt()
	ctxSelector := preparePromptConfigContextItems(cfgContexts)
	cfgContextsOptions := ctxSelector.Items.([]*config.Context)

	generateConfigContextSelectionMessage()

	i, _, err := ctxSelector.Run()
	if err != nil {
		return nil, err
	}

	logger.Debugf("Selected config context value. index: %d, name: %s", i, cfgContextsOptions[i].Name)
	return cfgContextsOptions[i], nil
}

func (p *prompterImpl) PromptApps(apps []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
	setSearchAppPrompt()
	appsSelector := preparePromptAppsItems(apps)
	appsOptions := appsSelector.Items.([]*models.ApplicationInfo)

	i, _, err := appsSelector.Run()
	if err != nil {
		return nil, err
	}

	logger.Debugf("Selected app value. index: %d, name: %s", i, appsOptions[i].Name)
	return appsOptions[i], nil
}

func (p *prompterImpl) PromptInstructionActions(appName string, actions []*models.Action) (*models.Action, error) {
	setSearchInstructionsPrompt(appName)
	instSelector := preparePromptInstructionsActions(actions)

	i, _, err := instSelector.Run()
	if err != nil {
		return nil, err
	}

	logger.Debugf("selected instruction action. index: %d, name: %s", i, actions[i].Id)
	return actions[i], nil
}

func (p *prompterImpl) PromptInstructionWorkflows(appName string, workflows []*models.Workflow) (*models.Workflow, error) {
	setSearchInstructionsPrompt(appName + " (workflows)")
	instSelector := preparePromptInstructionsWorkflows(workflows)

	i, _, err := instSelector.Run()
	if err != nil {
		return nil, err
	}

	logger.Debugf("selected instruction workflow. index: %d, name: %s", i, workflows[i].Id)
	return workflows[i], nil
}

func ClearPrompter(selector promptui.Select) {
	buf := screenbuf.New(selector.Stdout)
	err := buf.Clear()
	if err != nil {
		logger.Warningf("failed to clear the screen. error: %s", err.Error())
	}
}

func createPaddingLeftString(length int) string {
	// This is an example output "\"%-23s\""
	return "\"%-" + fmt.Sprintf("%d", length) + "s\""
}

func createCustomSpacesString(length int) string {
	spaces := ""
	for i := 0; i < length; i++ {
		spaces += " "
	}
	return spaces
}
