package printer

import (
	"github.com/manifoldco/promptui"
)

type AnchorFolderItemStatusTemplate struct {
	Name                     string
	IsValid                  bool
	MissingInstructionFile   bool
	InvalidInstructionFormat bool
}

func (as *AnchorFolderItemStatusTemplate) CalculateValidity() bool {
	as.IsValid = !as.MissingInstructionFile && !as.InvalidInstructionFormat
	return as.IsValid
}

var configViewTemplate = `{{ "Configuration Path:" | cyan }} 
{{ .ConfigFilePath | yellow }} 

{{ "Configuration:" | cyan }}
{{ .ConfigText | yellow }}
`

var appStatusTemplate = `{{ "There are " | cyan }}{{ .Count | green }}{{ " available items:" | cyan }}

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
