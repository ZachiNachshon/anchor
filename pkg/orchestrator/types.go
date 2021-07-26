package orchestrator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

const (
	identifier string = "orchestrator"
)

type Orchestrator interface {
	OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError)

	ExtractInstructions(
		app *models.ApplicationInfo,
		anchorfilesRepoPath string) (*models.InstructionsRoot, *errors.PromptError)

	OrchestrateInstructionActionSelection(
		app *models.ApplicationInfo,
		actions []*models.Action) (*models.Action, *errors.PromptError)

	OrchestrateInstructionWorkflowSelection(
		app *models.ApplicationInfo,
		workflows []*models.Workflow) (*models.Workflow, *errors.PromptError)

	AskBeforeRunningInstructionAction(action *models.Action) (bool, *errors.PromptError)
	AskBeforeRunningInstructionWorkflow(workflow *models.Workflow) (bool, *errors.PromptError)

	RunInstructionAction(action *models.Action) *errors.PromptError
	RunInstructionWorkflow(
		workflow *models.Workflow,
		actions []*models.Action) *errors.PromptError

	WrapAfterActionRun() *errors.PromptError
}

func ToRegistry(reg *registry.InjectionsRegistry, locator Orchestrator) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Orchestrator, error) {
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Orchestrator), nil
}
