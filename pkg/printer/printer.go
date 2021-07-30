package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/ZachiNachshon/anchor/pkg/utils/templates"
)

var configPrintTemplate = `{{ "Configuration Path:" | cyan }} 
{{ .ConfigFilePath | yellow }} 

{{ "Configuration:" | cyan }}
{{ .ConfigText | yellow }}
`

var appsPrintTemplate = `{{ "Available Applications (" | cyan }}{{ .Count | red }}{{ "):" | cyan }}{{ range $element := .AppsInfo }}
  • {{ $element.Name }}{{ end }}

`

type PrintConfigTemplateItems struct {
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

func (p *printerImpl) PrintApplications(apps []*models.ApplicationInfo) {
	data := struct {
		AppsInfo []*models.ApplicationInfo
		Count    int
	}{
		apps,
		len(apps),
	}
	if text, err := templates.TemplateToText(appsPrintTemplate, data); err != nil {
		logger.Error("Failed to prepare applications template string")
	} else {
		fmt.Print(text)
	}
}

func (p *printerImpl) PrintConfiguration(cfgFilePath string, cfgText string) {
	var items = PrintConfigTemplateItems{
		ConfigFilePath: cfgFilePath,
		ConfigText:     cfgText,
	}
	if text, err := templates.TemplateToText(configPrintTemplate, items); err != nil {
		logger.Error("Failed to prepare configuration template string")
	} else {
		fmt.Print(text)
	}
}
