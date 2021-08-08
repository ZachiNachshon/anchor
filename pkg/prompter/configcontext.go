package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/manifoldco/promptui"
	"strings"
)

func generateConfigContextSelectionMessage() {
	fmt.Printf("\n%sPLEASE SELECT CONFIG CONTEXT:%s\n\n", colors.Red, colors.Reset)
}

func preparePromptConfigContextItems(cfgContexts []*config.Context) promptui.Select {
	appsEnhanced := appendConfigContextCustomOptions(cfgContexts)
	appsTemplate := prepareConfigContextTemplate()
	appsSearcher := prepareConfigContextSearcher(appsEnhanced)
	return prepareConfigContextSelector(appsEnhanced, appsTemplate, appsSearcher)
}

func appendConfigContextCustomOptions(cfgContexts []*config.Context) []*config.Context {
	cfgContextsOptions := make([]*config.Context, 0, len(cfgContexts)+1)
	cancel := &config.Context{
		Name: CancelActionName,
	}
	cfgContextsOptions = append(cfgContextsOptions, cancel)
	cfgContextsOptions = append(cfgContextsOptions, cfgContexts...)
	return cfgContextsOptions
}

func setSearchConfigContextPrompt() {
	promptui.SearchPrompt = colors.Blue + "Search: " + colors.Reset
}

func prepareConfigContextTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   selectorEmoji + " {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: selectorEmoji + " {{ .Name | red | cyan }}",
		Details:  configContextPromptTemplateDetails,
	}
}

func prepareConfigContextSearcher(cfgContexts []*config.Context) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := cfgContexts[index]
		name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareConfigContextSelector(
	cfgContexts []*config.Context,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "",
		Items:             cfgContexts,
		Templates:         templates,
		Size:              5,
		Searcher:          searcher,
		StartInSearchMode: false,
		HideSelected:      true,
	}
}
