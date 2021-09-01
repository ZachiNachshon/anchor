package printer

import (
	"github.com/manifoldco/promptui"
)

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
	promptui.IconGood + ` {{ $element.Name }}
  {{- else }} ` +
	promptui.IconBad + ` {{ $element.Name }}
    {{- if (eq $element.MissingInstructionFile true) }}
    • {{ "Missing instructions.yaml file" | red }} 
    {{ end }}
    {{- if (eq $element.InvalidInstructionFormat true) }} 
    • {{ "Invalid instructions.yaml file format" | red }} 
    {{ end }}
  {{- end }}
{{ end }}
`
