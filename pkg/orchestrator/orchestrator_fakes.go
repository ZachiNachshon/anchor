package orchestrator

import (
	errors "github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/pkg/models"
)

var CreateFakeOrchestrator = func() *fakeOrchestratorImpl {
	return &fakeOrchestratorImpl{}
}

type fakeOrchestratorImpl struct {
	Orchestrator
	OrchestrateApplicationSelectionMock         func() (*models.ApplicationInfo, *errors.PromptError)
	ExtractInstructionsMock                     func(app *models.ApplicationInfo, anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError)
	OrchestrateInstructionActionSelectionMock   func(app *models.ApplicationInfo, actions []*models.Action) (*models.Action, *errors.PromptError)
	OrchestrateInstructionWorkflowSelectionMock func(app *models.ApplicationInfo, workflows []*models.Workflow) (*models.Workflow, *errors.PromptError)
	AskBeforeRunningInstructionActionMock       func(action *models.Action) (bool, *errors.PromptError)
	AskBeforeRunningInstructionWorkflowMock     func(workflow *models.Workflow) (bool, *errors.PromptError)
	RunInstructionActionMock                    func(action *models.Action) *errors.PromptError
	RunInstructionWorkflowMock                  func(workflow *models.Workflow, actions []*models.Action) *errors.PromptError
	WrapAfterActionRunMock                      func() *errors.PromptError
}

func (o *fakeOrchestratorImpl) OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError) {
	return o.OrchestrateApplicationSelectionMock()
}

func (o *fakeOrchestratorImpl) ExtractInstructions(
	app *models.ApplicationInfo,
	anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError) {

	return o.ExtractInstructionsMock(app, anchorfilesRepoPath)
}

func (o *fakeOrchestratorImpl) OrchestrateInstructionActionSelection(
	app *models.ApplicationInfo,
	actions []*models.Action) (*models.Action, *errors.PromptError) {

	return o.OrchestrateInstructionActionSelectionMock(app, actions)
}

func (o *fakeOrchestratorImpl) OrchestrateInstructionWorkflowSelection(
	app *models.ApplicationInfo,
	workflows []*models.Workflow) (*models.Workflow, *errors.PromptError) {

	return o.OrchestrateInstructionWorkflowSelectionMock(app, workflows)
}

func (o *fakeOrchestratorImpl) AskBeforeRunningInstructionAction(action *models.Action) (bool, *errors.PromptError) {
	return o.AskBeforeRunningInstructionActionMock(action)
}

func (o *fakeOrchestratorImpl) AskBeforeRunningInstructionWorkflow(workflow *models.Workflow) (bool, *errors.PromptError) {
	return o.AskBeforeRunningInstructionWorkflowMock(workflow)
}

func (o *fakeOrchestratorImpl) RunInstructionAction(action *models.Action) *errors.PromptError {
	return o.RunInstructionActionMock(action)
}

func (o *fakeOrchestratorImpl) RunInstructionWorkflow(workflow *models.Workflow, actions []*models.Action) *errors.PromptError {
	return o.RunInstructionWorkflowMock(workflow, actions)
}

func (o *fakeOrchestratorImpl) WrapAfterActionRun() *errors.PromptError {
	return o.WrapAfterActionRunMock()
}
