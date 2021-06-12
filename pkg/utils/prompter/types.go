package prompter

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
)

type Prompter interface {
	PromptApps(l locator.Locator) (*models.AppContent, error)
	PromptInstructions(*models.Instructions) (*models.PromptItem, error)
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
