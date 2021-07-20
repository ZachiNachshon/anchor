package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Prompter interface {
	PromptApps(appsArr []*models.ApplicationInfo) (*models.ApplicationInfo, error)
	PromptInstructions(appName string, instructionsRoot *models.InstructionsRoot) (*models.Action, error)
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

var appsPromptTemplateDetails = fmt.Sprintf(`{{ if not (eq .Name "%s") }}
{{ "Information:" | blue }}

{{ "Name:" | faint }}	{{ .Name }}
{{ "Overview:" | faint }}	{{ .DirPath }}
{{ else }}
Exit Application
{{ end }}`, CancelActionName)

var instructionsPromptTemplateDetails = fmt.Sprintf(`{{ if not (eq .Id "%s") }}
{{ "Information:" | blue }}

{{ "Id:" | faint }}	{{ .Id }}
{{ "Title:" | faint }}	{{ .Title }}
{{ "File:" | faint }}	{{ .File }}
{{ else }}
Go Back (App Selector)
{{ end }}`, BackActionName)
