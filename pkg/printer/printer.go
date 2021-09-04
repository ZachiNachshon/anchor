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
	PrintEmptyLines(count int)
	PrintSuccess(message string)
	PrintWarning(message string)

	PrepareRunActionPlainer(actionId string) PrinterPlainer
	PrepareRunActionSpinner(actionId string, scriptOutputPath string) PrinterSpinner
	PrepareReadRemoteHeadCommitHashSpinner(url string, branch string) PrinterSpinner
	PrepareCloneRepositorySpinner(url string, branch string) PrinterSpinner
	PrepareResetToRevisionSpinner(revision string) PrinterSpinner
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

func (p *printerImpl) PrintEmptyLines(count int) {
	for i := 0; i < count; i++ {
		fmt.Println()
	}
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

func (p *printerImpl) PrintSuccess(message string) {
	fmt.Printf("%s %s", promptui.IconGood, message)
	fmt.Println()
}

func (p *printerImpl) PrintWarning(message string) {
	fmt.Printf("%s %s", promptui.IconWarn, message)
	fmt.Println()
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

func (p *printerImpl) PrepareReadRemoteHeadCommitHashSpinner(url string, branch string) PrinterSpinner {
	return NewSpinner(
		getSpinnerReadRemoteHeadCommitHashMessage(url, branch),
		getSpinnerReadRemoteHeadCommitHashSuccessMessage(),
		getSpinnerReadRemoteHeadCommitHashFailureMessageFormat())
}

func (p *printerImpl) PrepareCloneRepositorySpinner(url string, branch string) PrinterSpinner {
	return NewSpinner(
		getCloneRepositoryMessage(url, branch),
		getCloneRepositorySuccessMessage(),
		getCloneRepositoryFailureMessageFormat())
}

func (p *printerImpl) PrepareResetToRevisionSpinner(revision string) PrinterSpinner {
	return NewSpinner(
		getResetToRevisionMessage(revision),
		getResetToRevisionSuccessMessage(revision),
		getResetToRevisionFailureMessageFormat(revision))
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
	return fmt.Sprintf("%s Action %s%s%s completed successfully",
		promptui.IconGood, colors.Cyan, actionId, colors.Reset)
}

func getSpinnerFailureActionMessageFormat(actionId string, scriptOutputPath string) string {
	return fmt.Sprintf(`%s Action %s%s%s failed
    
    Reason: %%s
    Output: %s
`, promptui.IconBad, colors.Cyan, actionId, colors.Reset, scriptOutputPath)
}

func getSpinnerReadRemoteHeadCommitHashMessage(url string, branch string) string {
	return fmt.Sprintf(" Reading HEAD commit-hash from remote repository (url: %s, branch: %s)...", url, branch)
}

func getSpinnerReadRemoteHeadCommitHashSuccessMessage() string {
	return fmt.Sprintf("%s Read remote repository HEAD commmit-hash", promptui.IconGood)
}

func getSpinnerReadRemoteHeadCommitHashFailureMessageFormat() string {
	return fmt.Sprintf(`%s Failed to read remote repository HEAD commit-hash.
    
    Reason: %%s

`, promptui.IconBad)
}

func getCloneRepositoryMessage(url string, branch string) string {
	return fmt.Sprintf(" Cloning repository (url: %s, branch: %s)...", url, branch)
}

func getCloneRepositorySuccessMessage() string {
	return fmt.Sprintf("%s Repository cloned successfully", promptui.IconGood)
}

func getCloneRepositoryFailureMessageFormat() string {
	return fmt.Sprintf(`%s Cloning remote repository failed.
    
    Reason: %%s

`, promptui.IconBad)
}

func getResetToRevisionMessage(revision string) string {
	return fmt.Sprintf(" Resetting to revision %s...", revision)
}

func getResetToRevisionSuccessMessage(revision string) string {
	return fmt.Sprintf("%s Reset to revision %s", promptui.IconGood, revision)
}

func getResetToRevisionFailureMessageFormat(revision string) string {
	return fmt.Sprintf(`%s Resetting to revision %s failed.
    
    Reason: %%s

`, promptui.IconBad, revision)
}
