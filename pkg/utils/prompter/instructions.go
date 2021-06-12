package prompter

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
	"github.com/manifoldco/promptui"
	"strings"
)

func preparePromptInstructionsItems(instructions *parser.Instructions) promptui.Select {
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
		Id: backButtonName,
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
