package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
	"github.com/manifoldco/promptui"
	"strings"
)

type prompterImpl struct {
	Prompter
}

func New() Prompter {
	return &prompterImpl{}
}

func (p *prompterImpl) PromptApps(l locator.Locator) (*locator.AppContent, error) {
	setSearchAppPrompt()
	appsArr := l.Applications()
	appsSelector := prepareAppsItems(appsArr)
	i, _, err := appsSelector.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and run apps prompt selector. error: %s", err.Error())
	}
	logger.Debugf("selected app value. index: %d, name: %s", i+1, appsArr[i].Name)
	return appsArr[i], nil
}

func (p *prompterImpl) PromptInstructions(instructions *parser.Instructions) (*parser.PromptItem, error) {
	setSearchInstructionsPrompt()
	instSelector := prepareInstructionsItems(instructions)
	i, _, err := instSelector.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and run instruction prompt selector. error: %s", err.Error())
	}
	logger.Debugf("selected instruction value. index: %d, name: %s", i+1, instructions.Items[i].Id)
	return instructions.Items[i], nil
}

func prepareAppsItems(apps []*locator.AppContent) promptui.Select {
	apps = appendAppsCustomOptions(apps)
	appsTemplate := prepareAppsTemplate()
	appsSearcher := prepareAppsSearcher(apps)
	return prepareAppsSelector(apps, appsTemplate, appsSearcher)
}

func appendAppsCustomOptions(apps []*locator.AppContent) []*locator.AppContent {
	appDirs := make([]*locator.AppContent, 0, len(apps)+1)
	back := &locator.AppContent{
		Name: "Back",
	}
	appDirs = append(appDirs, back)
	appDirs = append(appDirs, apps...)
	return appDirs
}

func setSearchAppPrompt() {
	promptui.SearchPrompt = "Search App: "
}

func prepareAppsTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   promptui.IconSelect + " {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: promptui.IconSelect + " {{ .Name | red | cyan }}",
		Details:  appsPromptTemplateDetails,
	}
}

func prepareAppsSearcher(apps []*locator.AppContent) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := apps[index]
		name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareAppsSelector(
	apps []*locator.AppContent,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "Available Applications",
		Items:             apps,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
	}
}

func prepareInstructionsItems(instructions *parser.Instructions) promptui.Select {
	appendInstructionCustomOptions(instructions)
	instTemplate := prepareInstructionsTemplate()
	instSearcher := prepareInstructionsSearcher(instructions.Items)
	return prepareInstructionsSelector(instructions, instTemplate, instSearcher)
}

func setSearchInstructionsPrompt() {
	promptui.SearchPrompt = "Search Instruction: "
}

func appendInstructionCustomOptions(instructions *parser.Instructions) {
	instItems := make([]*parser.PromptItem, 0, len(instructions.Items)+1)
	back := &parser.PromptItem{
		Id: "Back",
	}
	instItems = append(instItems, back)
	instItems = append(instItems, instructions.Items...)
	instructions.Items = instItems
}

func prepareInstructionsTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   promptui.IconSelect + " {{ .Id | cyan }} ({{ .Title | red }})",
		Inactive: "  {{ .Id | cyan }} ({{ .Title | red }})",
		Selected: promptui.IconSelect + " {{ .Id | red | cyan }}",
		Details:  instructionsPromptTemplateDetails,
	}
}

func prepareInstructionsSearcher(items []*parser.PromptItem) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := items[index]
		name := strings.Replace(strings.ToLower(item.Id), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareInstructionsSelector(
	instructions *parser.Instructions,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "Available Instructions",
		Items:             instructions.Items,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
	}
}
