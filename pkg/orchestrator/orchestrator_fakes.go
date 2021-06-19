package orchestrator

import (
	"github.com/ZachiNachshon/anchor/models"
)

var CreateFakeOrchestrator = func() *fakeOrchestratorImpl {
	return &fakeOrchestratorImpl{}
}

type fakeOrchestratorImpl struct {
	Orchestrator
	OrchestrateAppInstructionSelectionMock func() (*models.PromptItem, error)
}

func (o *fakeOrchestratorImpl) OrchestrateAppInstructionSelection() (*models.PromptItem, error) {
	return o.OrchestrateAppInstructionSelectionMock()
}
