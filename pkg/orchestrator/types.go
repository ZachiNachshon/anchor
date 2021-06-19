package orchestrator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

const (
	identifier string = "orchestrator"
)

type Orchestrator interface {
	OrchestrateAppInstructionSelection() (*models.PromptItem, error)
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
