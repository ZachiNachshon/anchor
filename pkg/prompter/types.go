package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Prompter interface {
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

var appsPromptTemplateDetails = fmt.Sprintf(`{{ if not (eq .Name "%s") }}
{{ "Information:" | blue }}

{{ "Name:" | faint }}	{{ .Name }}
{{ "Overview:" | faint }}	{{ .DirPath }}
{{ else }}
Exit Application
{{ end }}`, CancelActionName)

var actionLengthiestOption = len("Description:")
var leftPadding = 3
var instructionsActionPromptTemplateDetails = `{{ if not (eq .Id "` + BackActionName + `") }}
{{ "Information:\n" | blue }}

{{- if not (eq .Id "") }}
  {{ "Id:" | faint }}` + createCustomSpacesString(actionLengthiestOption-3+leftPadding) + `{{ .Id }}
{{- end }}
{{- if not (eq .Title "") }}
  {{ "Title:" | faint }}` + createCustomSpacesString(actionLengthiestOption-6+leftPadding) + `{{ .Title }}
{{- end }}
{{- if not (eq .Script "") }} 
  {{ "Script:" | faint }}` + createCustomSpacesString(actionLengthiestOption-7+leftPadding) + `{{ "(hidden)" }} 
{{- end }}
{{- if not (eq .ScriptFile "") }}
  {{ "ScriptFile:" | faint }}` + createCustomSpacesString(actionLengthiestOption-11+leftPadding) + `{{ .ScriptFile }} 
{{- end }}
{{- if not (eq .Description "") }}
  {{ "Description:" | faint }}` + createCustomSpacesString(leftPadding) + `{{ .Description }}
{{- end }}
{{ else }}
  Go Back (App Selector)
{{- end }}`

var workflowLengthiestOption = len("Tolerate Failures:")
var instructionsWorkflowPromptTemplateDetails = `{{ if not (eq .Id "` + BackActionName + `") }}
{{ "Information:\n" | blue }}

{{- if not (eq .Id "") }}
  {{ "Id:" | faint }}` + createCustomSpacesString(workflowLengthiestOption-3+leftPadding) + `{{ .Id }}
{{- end }}
{{- if true }}
  {{ "Tolerate Failures:" | faint }}` + createCustomSpacesString(leftPadding) + `{{ .TolerateFailures }}
{{- end }}
{{- if not (eq .Description "") }}
  {{ "Description:" | faint }}` + createCustomSpacesString(workflowLengthiestOption-12+leftPadding) + `{{ .Description }}
{{- end }}
{{- if true }}
  {{ "Action Ids:" | faint }}
  {{ range $element := .ActionIds }}` +
	createCustomSpacesString(workflowLengthiestOption+leftPadding) + `• {{ $element }}
  {{ end }}
{{- end }}
{{ else }}
  Go Back (Actions Selector)
{{- end }}`
