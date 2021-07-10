package prompter

import "github.com/ZachiNachshon/anchor/models"

var CreateFakePrompter = func() *fakePrompterImpl {
	return &fakePrompterImpl{}
}

type fakePrompterImpl struct {
	Prompter
	PromptAppsMock         func(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error)
	PromptInstructionsMock func(appName string, instructions *models.Instructions) (*models.InstructionItem, error)
}

func (p *fakePrompterImpl) PromptApps(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
	appsArr = appendAppsCustomOptions(appsArr)
	return p.PromptAppsMock(appsArr)
}

func (p *fakePrompterImpl) PromptInstructions(appName string, instructions *models.Instructions) (*models.InstructionItem, error) {
	appendInstructionCustomOptions(instructions)
	return p.PromptInstructionsMock(appName, instructions)
}
