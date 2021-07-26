package prompter

import "github.com/ZachiNachshon/anchor/models"

var CreateFakePrompter = func() *fakePrompterImpl {
	return &fakePrompterImpl{}
}

type fakePrompterImpl struct {
	Prompter
	PromptAppsMock                 func(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error)
	PromptInstructionActionsMock   func(appName string, actions []*models.Action) (*models.Action, error)
	PromptInstructionWorkflowsMock func(appName string, workflows []*models.Workflow) (*models.Workflow, error)
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
