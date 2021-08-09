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
	PromptConfigContextMock        func(cfgContexts []*config.Context) (*config.Context, error)
	PromptAppsMock                 func(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error)
	PromptInstructionActionsMock   func(appName string, actions []*models.Action) (*models.Action, error)
	PromptInstructionWorkflowsMock func(appName string, workflows []*models.Workflow) (*models.Workflow, error)
}

func (p *fakePrompterImpl) PromptConfigContext(cfgContexts []*config.Context) (*config.Context, error) {
	return p.PromptConfigContextMock(cfgContexts)
}

func (p *fakePrompterImpl) PromptApps(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
	appsArr = appendAppsCustomOptions(appsArr)
	return p.PromptAppsMock(appsArr)
}

func (p *fakePrompterImpl) PromptInstructionActions(appName string, actions []*models.Action) (*models.Action, error) {
	return p.PromptInstructionActionsMock(appName, actions)
}

func (p *fakePrompterImpl) PromptInstructionWorkflows(appName string, workflows []*models.Workflow) (*models.Workflow, error) {
	return p.PromptInstructionWorkflowsMock(appName, workflows)
}
