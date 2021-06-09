package prompter

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
)

type Prompter interface {
	PromptApps(l locator.Locator) (*locator.AppContent, error)
	PromptInstructions(*parser.Instructions) (*parser.PromptItem, error)
}

var appsPromptTemplateDetails = `{{ if not (eq .Name "Back") }}
--------- Information ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Overview:" | faint }}	{{ .DirPath }}
{{ else }}
Go back to previous step
{{ end }}`

var instructionsPromptTemplateDetails = `{{ if not (eq .Name "Back") }}
--------- Information ----------
{{ "Id:" | faint }}	{{ .Id }}
{{ "Title:" | faint }}	{{ .Title }}
{{ "File:" | faint }}	{{ .File }}
{{ else }}
Go back to previous step
{{ end }}`
