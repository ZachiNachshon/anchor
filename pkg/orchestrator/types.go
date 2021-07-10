package orchestrator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/errors"
	"github.com/ZachiNachshon/anchor/pkg/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

const (
	identifier string = "orchestrator"
)

type Orchestrator interface {
	OrchestrateApplicationSelection() (*models.ApplicationInfo, *errors.PromptError)
	OrchestrateInstructionSelection(app *models.ApplicationInfo) (*models.InstructionItem, *errors.PromptError)
	AskBeforeRunningInstruction(item *models.InstructionItem, in input.UserInput) (bool, *errors.PromptError)
	RunInstruction(item *models.InstructionItem, s shell.Shell) *errors.PromptError
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
