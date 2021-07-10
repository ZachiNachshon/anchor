package orchestrator

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

var CreateFakeOrchestrator = func() *fakeOrchestratorImpl {
	return &fakeOrchestratorImpl{}
}

type fakeOrchestratorImpl struct {
	Orchestrator
	OrchestrateApplicationSelectionMock func() (*models.ApplicationInfo, *errors.PromptError)
	OrchestrateInstructionSelectionMock func(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError)
	AskBeforeRunningInstructionMock     func(item *models.InstructionItem, in input.UserInput) (bool, *errors.PromptError)
	RunInstructionMock                  func(item *models.InstructionItem, s shell.Shell) *errors.PromptError
}

func (o *fakeOrchestratorImpl) OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError) {
	return o.OrchestrateApplicationSelectionMock()
}

func (o *fakeOrchestratorImpl) OrchestrateInstructionSelection(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError) {
	return o.OrchestrateInstructionSelectionMock(app)
}

func (o *fakeOrchestratorImpl) AskBeforeRunningInstruction(item *models.InstructionItem, in input.UserInput) (bool, *errors.PromptError) {
	return o.AskBeforeRunningInstructionMock(item, in)
}

func (o *fakeOrchestratorImpl) RunInstruction(item *models.InstructionItem, s shell.Shell) *errors.PromptError {
	return o.RunInstructionMock(item, s)
}
