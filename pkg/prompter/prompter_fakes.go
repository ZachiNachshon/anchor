package prompter

import "github.com/ZachiNachshon/anchor/models"

var CreateFakePrompter = func() *fakePrompterImpl {
	return &fakePrompterImpl{}
}

type fakePrompterImpl struct {
	Prompter
	PromptAppsMock         func(appsArr []*models.AppContent) (*models.AppContent, error)
	PromptInstructionsMock func(instructions *models.Instructions) (*models.PromptItem, error)
}

func (p *fakePrompterImpl) PromptApps(appsArr []*models.AppContent) (*models.AppContent, error) {
	appsArr = appendAppsCustomOptions(appsArr)
	return p.PromptAppsMock(appsArr)
}

func (p *fakePrompterImpl) PromptInstructions(instructions *models.Instructions) (*models.PromptItem, error) {
	appendInstructionCustomOptions(instructions)
	return p.PromptInstructionsMock(instructions)
}
