package prompter

import "fmt"

var configContextPromptTemplateDetails = fmt.Sprintf(`{{ if not (eq .Name "%s") }}
{{ "Hint:" | blue }}
{{ "Lock the selected config context as the active one using:" | faint }}
  • {{ "anchor" | green }} config set-context {{ .Name | cyan }}

{{ else }}
{{ "Exit Application" | faint }}
{{ end }}`, CancelActionName)

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