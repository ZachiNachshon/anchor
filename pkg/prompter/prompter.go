package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/models"

	"github.com/manifoldco/promptui"
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

const (
	Identifier string = "prompter"
)

type Prompter interface {
	PromptConfigContext(cfgContexts []*config.Context) (*config.Context, error)
	PromptCommandFolderItemSelection(folderItems []*models.CommandFolderItemInfo) (*models.CommandFolderItemInfo, error)
	PromptInstructionActions(folderItem string, actions []*models.Action) (*models.Action, error)
	PromptInstructionWorkflows(folderItem string, workflows []*models.Workflow) (*models.Workflow, error)
}

type prompterImpl struct {
	Prompter

	runConfigCtxSelectorFunc         func(promptui.Select) (int, string, error)
	runCommandFolderItemSelectorFunc func(promptui.Select) (int, string, error)
	runActionSelectorFunc            func(promptui.Select) (int, string, error)
	runWorkflowSelectorFunc          func(promptui.Select) (int, string, error)
}

func New() *prompterImpl {
	return &prompterImpl{
		runConfigCtxSelectorFunc:         runPromptSelector,
		runCommandFolderItemSelectorFunc: runPromptSelector,
		runActionSelectorFunc:            runPromptSelector,
		runWorkflowSelectorFunc:          runPromptSelector,
	}
}

func (p *prompterImpl) PromptConfigContext(cfgContexts []*config.Context) (*config.Context, error) {
	setSearchConfigContextPrompt()
	ctxSelector := preparePromptConfigContextItems(cfgContexts)
	cfgContextsOptions := ctxSelector.Items.([]*config.Context)

	generateConfigContextSelectionMessage()

	i, _, err := p.runConfigCtxSelectorFunc(ctxSelector)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Selected config context value. index: %d, name: %s", i, cfgContextsOptions[i].Name)
	return cfgContextsOptions[i], nil
}

func (p *prompterImpl) PromptCommandFolderItemSelection(folderItems []*models.CommandFolderItemInfo) (*models.CommandFolderItemInfo, error) {
	setSearchFolderItemPrompt()
	folderItemsSelector := preparePromptFolderItemItems(folderItems)
	folderItemsOptions := folderItemsSelector.Items.([]*models.CommandFolderItemInfo)

	i, _, err := p.runCommandFolderItemSelectorFunc(folderItemsSelector)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Selected app value. index: %d, name: %s", i, folderItemsOptions[i].Name)
	return folderItemsOptions[i], nil
}

func (p *prompterImpl) PromptInstructionActions(appName string, actions []*models.Action) (*models.Action, error) {
	setSearchInstructionsPrompt(appName)
	instSelector := preparePromptInstructionsActions(actions)

	i, _, err := p.runActionSelectorFunc(instSelector)
	if err != nil {
		return nil, err
	}

	logger.Debugf("selected instruction action. index: %d, name: %s", i, actions[i].Id)
	return actions[i], nil
}

func (p *prompterImpl) PromptInstructionWorkflows(appName string, workflows []*models.Workflow) (*models.Workflow, error) {
	setSearchInstructionsPrompt(appName + " (workflows)")
	instSelector := preparePromptInstructionsWorkflows(workflows)

	i, _, err := p.runWorkflowSelectorFunc(instSelector)
	if err != nil {
		return nil, err
	}

	logger.Debugf("selected instruction workflow. index: %d, name: %s", i, workflows[i].Id)
	return workflows[i], nil
}

//func ClearPrompter(selector promptui.Select) {
//	buf := screenbuf.New(selector.Stdout)
//	err := buf.Clear()
//	if err != nil {
//		logger.Warningf("failed to clear the screen. error: %s", err.Error())
//	}
//}

func runPromptSelector(selector promptui.Select) (int, string, error) {
	return selector.Run()
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
