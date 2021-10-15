package prompter

import (
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/pkg/models"
)

var CreateFakePrompter = func() *fakePrompterImpl {
	return &fakePrompterImpl{}
}

type fakePrompterImpl struct {
	Prompter
	PromptConfigContextMock              func(cfgContexts []*config.Context) (*config.Context, error)
	PromptCommandFolderItemSelectionMock func(folderItems []*models.CommandFolderItemInfo) (*models.CommandFolderItemInfo, error)
	PromptInstructionActionsMock         func(folderItem string, actions []*models.Action) (*models.Action, error)
	PromptInstructionWorkflowsMock       func(folderItem string, workflows []*models.Workflow) (*models.Workflow, error)
}

func (p *fakePrompterImpl) PromptConfigContext(cfgContexts []*config.Context) (*config.Context, error) {
	return p.PromptConfigContextMock(cfgContexts)
}

func (p *fakePrompterImpl) PromptCommandFolderItemSelection(appsArr []*models.CommandFolderItemInfo) (*models.CommandFolderItemInfo, error) {
	appsArr = appendFolderItemsCustomOptions(appsArr)
	return p.PromptCommandFolderItemSelectionMock(appsArr)
}

func (p *fakePrompterImpl) PromptInstructionActions(folderItem string, actions []*models.Action) (*models.Action, error) {
	return p.PromptInstructionActionsMock(folderItem, actions)
}

func (p *fakePrompterImpl) PromptInstructionWorkflows(folderItem string, workflows []*models.Workflow) (*models.Workflow, error) {
	return p.PromptInstructionWorkflowsMock(folderItem, workflows)
}
