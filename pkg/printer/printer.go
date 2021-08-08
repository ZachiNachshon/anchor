package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/logger"

	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/ZachiNachshon/anchor/pkg/utils/templates"
)

const (
	Identifier string = "printer"
)

type Printer interface {
	PrintAnchorBanner()
	PrintAnchorVersion(version string)
	PrintApplications(appsStatus []*AppStatusTemplateItem)
	PrintConfiguration(cfgFilePath string, cfgText string)
}

func (as *AppStatusTemplateItem) CalculateValidity() {
	as.IsValid = !as.MissingInstructionFile && !as.InvalidInstructionFormat
}

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
