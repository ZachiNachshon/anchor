package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Prompter interface {
	PromptApps(appsArr []*models.AppContent) (*models.AppContent, error)
	PromptInstructions(instructions *models.Instructions) (*models.PromptItem, error)
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
	locate := reg.Get(identifier).(Prompter)
	if locate == nil {
		return nil, fmt.Errorf("failed to retrieve prompter from registry")
	}
	return locate, nil
}

type Orchestrator interface {
	OrchestrateAppInstructionSelection() (*models.PromptItem, error)
}

var appsPromptTemplateDetails = `{{ if not (eq .Name "cancel") }}
--------- Information ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Overview:" | faint }}	{{ .DirPath }}
{{ else }}
Cancel application selector
{{ end }}`

var instructionsPromptTemplateDetails = `{{ if not (eq .Id "back") }}
--------- Information ----------
{{ "Id:" | faint }}	{{ .Id }}
{{ "Title:" | faint }}	{{ .Title }}
{{ "File:" | faint }}	{{ .File }}
{{ else }}
Go back to previous step
{{ end }}`
