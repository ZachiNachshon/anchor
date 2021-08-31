package printer

import "github.com/ZachiNachshon/anchor/pkg/prompter"

type AppStatusTemplateItem struct {
	Name                     string
	IsValid                  bool
	MissingInstructionFile   bool
	InvalidInstructionFormat bool
}

var configViewTemplate = `{{ "Configuration Path:" | cyan }} 
{{ .ConfigFilePath | yellow }} 

{{ "Configuration:" | cyan }}
{{ .ConfigText | yellow }}
`

var appStatusTemplate = `{{ "There are " | cyan }}{{ .Count | green }}{{ " available applications:" | cyan }}

{{ range $element := .AppsStatusItems }}
  {{- if (eq $element.IsValid true) }} ` +
	prompter.CheckMarkEmoji + ` {{ $element.Name }}
  {{- else }} ` +
	prompter.CrossMarkEmoji + ` {{ $element.Name }}
    {{- if (eq $element.MissingInstructionFile true) }}
    • {{ "Missing instructions.yaml file" | red }} 
    {{ end }}
    {{- if (eq $element.InvalidInstructionFormat true) }} 
    • {{ "Invalid instructions.yaml file format" | red }} 
    {{ end }}
  {{- end }}
{{ end }}
`
