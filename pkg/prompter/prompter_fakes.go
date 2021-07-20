package prompter

import "github.com/ZachiNachshon/anchor/models"

var CreateFakePrompter = func() *fakePrompterImpl {
	return &fakePrompterImpl{}
}

type fakePrompterImpl struct {
	Prompter
	PromptAppsMock         func(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error)
	PromptInstructionsMock func(appName string, instructionsRoot *models.InstructionsRoot) (*models.Action, error)
}

func (p *fakePrompterImpl) PromptApps(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error) {
	appsArr = appendAppsCustomOptions(appsArr)
	return p.PromptAppsMock(appsArr)
}

func (p *fakePrompterImpl) PromptInstructions(appName string, instructionsRoot *models.InstructionsRoot) (*models.Action, error) {
	instructions := instructionsRoot.Instructions
	appendInstructionCustomOptions(instructions)
	return p.PromptInstructionsMock(appName, instructionsRoot)
}
