package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Prompter interface {
	PromptConfigContext(cfgContexts []*config.Context) (*config.Context, error)
	PromptApps(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error)
	PromptInstructionActions(appName string, actions []*models.Action) (*models.Action, error)
	PromptInstructionWorkflows(appName string, workflows []*models.Workflow) (*models.Workflow, error)
}

const (
	identifier string = "prompter"
)

func ToRegistry(reg *registry.InjectionsRegistry, locator Prompter) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Prompter, error) {
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Prompter), nil
}
