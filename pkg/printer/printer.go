package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/ZachiNachshon/anchor/pkg/utils/templates"
	"github.com/manifoldco/promptui"
)

const (
	Identifier string = "printer"
)

type Printer interface {
	PrintAnchorBanner()
	PrintAnchorVersion(version string)
	PrintApplicationsStatus(appsStatus []*AppStatusTemplateItem)
	PrintConfiguration(cfgFilePath string, cfgText string)
	PrintMissingInstructions()

	PrepareRunActionPlainer(actionId string) PrinterPlainer
	PrepareRunActionSpinner(actionId string, scriptOutputPath string) PrinterSpinner
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

func (p *printerImpl) PrintApplicationsStatus(appsStatus []*AppStatusTemplateItem) {
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

func (p *printerImpl) PrintMissingInstructions() {
	fmt.Printf("%sMissing file: %s%s\n\n", colors.Red, globals.InstructionsFileName, colors.Reset)
}

func (p *printerImpl) PrepareRunActionPlainer(actionId string) PrinterPlainer {
	return NewPlainer(
		getPlainerRunActionMessage(actionId),
		getPlainerSuccessActionMessage(actionId),
		getPlainerFailureActionMessageFormat(actionId))
}

func (p *printerImpl) PrepareRunActionSpinner(actionId string, scriptOutputPath string) PrinterSpinner {
	return NewSpinner(
		getSpinnerRunActionMessage(actionId),
		getSpinnerSuccessActionMessage(actionId),
		getSpinnerFailureActionMessageFormat(actionId, scriptOutputPath))
}

func getPlainerRunActionMessage(actionId string) string {
	return fmt.Sprintf(`==> Running %s...

Output:`, actionId)
}

func getPlainerSuccessActionMessage(actionId string) string {
	return fmt.Sprintf(`
Result:
%s Action %s%s%s completed successfully

`,
		promptui.IconGood, colors.Cyan, actionId, colors.Reset)
}

func getPlainerFailureActionMessageFormat(actionId string) string {
	return fmt.Sprintf(`
Result:
%s Action %s%s%s failed

`, promptui.IconBad, colors.Cyan, actionId, colors.Reset)
}

func getSpinnerRunActionMessage(actionId string) string {
	return fmt.Sprintf(" Running %s...", actionId)
}

func getSpinnerSuccessActionMessage(actionId string) string {
	return fmt.Sprintf("%s Action %s%s%s completed successfully\n\n",
		promptui.IconGood, colors.Cyan, actionId, colors.Reset)
}

func getSpinnerFailureActionMessageFormat(actionId string, scriptOutputPath string) string {
	return fmt.Sprintf(`%s Action %s%s%s failed
    
    Reason: %%s
    Output: %s

`, promptui.IconBad, colors.Cyan, actionId, colors.Reset, scriptOutputPath)
}
