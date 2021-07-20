package orchestrator

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
)

var CreateFakeOrchestrator = func() *fakeOrchestratorImpl {
	return &fakeOrchestratorImpl{}
}

type fakeOrchestratorImpl struct {
	Orchestrator
	OrchestrateApplicationSelectionMock func() (*models.ApplicationInfo, *errors.PromptError)
	OrchestrateInstructionSelectionMock func(app *models.ApplicationInfo) (*models.Action, *errors.PromptError)
	AskBeforeRunningInstructionMock     func(item *models.Action) (bool, *errors.PromptError)
	RunInstructionMock                  func(item *models.Action, repoPath string) *errors.PromptError
}

func (o *fakeOrchestratorImpl) OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError) {
	return o.OrchestrateApplicationSelectionMock()
}

func (o *fakeOrchestratorImpl) OrchestrateInstructionSelection(app *models.ApplicationInfo) (*models.Action, *errors.PromptError) {
	return o.OrchestrateInstructionSelectionMock(app)
}

func (o *fakeOrchestratorImpl) AskBeforeRunningInstruction(item *models.Action) (bool, *errors.PromptError) {
	return o.AskBeforeRunningInstructionMock(item)
}

func (o *fakeOrchestratorImpl) RunInstruction(item *models.Action, repoPath string) *errors.PromptError {
	return o.RunInstructionMock(item, repoPath)
}
