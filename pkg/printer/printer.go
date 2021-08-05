package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/ZachiNachshon/anchor/pkg/utils/templates"
)

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
{{ end }}`

type ConfigViewTemplateItem struct {
	ConfigFilePath string
	ConfigText     string
}

type printerImpl struct {
	Printer
}

func New() Printer {
	return &printerImpl{}
}

func (p *printerImpl) PrintAnchorBanner() {
	fmt.Printf(colors.Blue + `
     \                  |                  
    _ \    __ \    __|  __ \    _ \    __| 
   ___ \   |   |  (     | | |  (   |  |    
 _/    _\ _|  _| \___| _| |_| \___/  _|

` + colors.Reset)
}

func (p *printerImpl) PrintAnchorVersion(version string) {
	fmt.Println(version)
}

func (p *printerImpl) PrintApplications(appsStatus []*AppStatusTemplateItem) {
	data := struct {
		AppsStatusItems []*AppStatusTemplateItem
		Count           int
	}{
		appsStatus,
		len(appsStatus),
	}
	if text, err := templates.TemplateToText(appStatusTemplate, data); err != nil {
		logger.Error("Failed to prepare applications template string")
	} else {
		fmt.Print(text)
	}
}

func (p *printerImpl) PrintConfiguration(cfgFilePath string, cfgText string) {
	var items = ConfigViewTemplateItem{
		ConfigFilePath: cfgFilePath,
		ConfigText:     cfgText,
	}
	if text, err := templates.TemplateToText(configViewTemplate, items); err != nil {
		logger.Error("Failed to prepare configuration template string")
	} else {
		fmt.Print(text)
	}
}
