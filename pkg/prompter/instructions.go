package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/manifoldco/promptui"
	"strings"
)

func GenerateRunInstructionMessage(id string, title string) string {
	return fmt.Sprintf("Run instruction %s'%s'%s %s(%s)%s?",
		colors.Cyan, id, colors.Reset, colors.Purple, title, colors.Reset)
}

func preparePromptInstructionsItems(instructions *models.Instructions) promptui.Select {
	appendInstructionCustomOptions(instructions)
	instTemplate := prepareInstructionsTemplate(calculatePadding(instructions))
	instSearcher := prepareInstructionsSearcher(instructions.Items)
	return prepareInstructionsSelector(instructions, instTemplate, instSearcher)
}

func calculatePadding(instructions *models.Instructions) (string, string) {
	length := findLongestInstructionNameLength(instructions)
	return createPaddingString(length), createPaddingString(length + 2)
}

func findLongestInstructionNameLength(instructions *models.Instructions) int {
	maxNameLen := 0
	for _, v := range instructions.Items {
		instNameLen := len(v.Id)
		if instNameLen > maxNameLen {
			maxNameLen = instNameLen
		}
	}
	return maxNameLen
}

func setSearchInstructionsPrompt(appName string) {
	promptui.SearchPrompt = fmt.Sprintf("%sSearch%s %s%s :%s ", colors.Blue, colors.Green, appName, colors.Blue, colors.Reset)
}

func appendInstructionCustomOptions(instructions *models.Instructions) {
	instItems := make([]*models.InstructionItem, 0, len(instructions.Items)+1)
	back := &models.InstructionItem{
		Id: BackButtonName,
	}
	workflows := &models.InstructionItem{
		Id: "workflows...",
	}
	instItems = append(instItems, back)
	instItems = append(instItems, workflows)
	instItems = append(instItems, instructions.Items...)
	instructions.Items = instItems
}

func prepareInstructionsTemplate(activeSpacePadding string, inactiveSpacePadding string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   selectorEmoji + ` {{ printf ` + activeSpacePadding + ` .Id | cyan }} {{ if not (eq .Id "back") }} ({{ .Title | green }}) {{ end }}`,
		Inactive: ` {{ printf ` + inactiveSpacePadding + ` .Id | cyan }} {{ if not (eq .Id "back") }} ({{ .Title | faint }}) {{ end }}`,
		Selected: selectorEmoji + " {{ .Id | red | cyan }}",
		Details:  instructionsPromptTemplateDetails,
	}
}

func prepareInstructionsSearcher(items []*models.InstructionItem) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := items[index]
		name := strings.Replace(strings.ToLower(item.Id), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareInstructionsSelector(
	instructions *models.Instructions,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "",
		Items:             instructions.Items,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
	}
}
