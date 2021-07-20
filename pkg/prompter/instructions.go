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

func preparePromptInstructionsActions(instructions *models.Instructions) promptui.Select {
	appendInstructionCustomOptions(instructions)
	instTemplate := prepareInstructionsTemplate(calculatePadding(instructions))
	instSearcher := prepareInstructionsSearcher(instructions.Actions)
	return prepareInstructionsSelector(instructions, instTemplate, instSearcher)
}

func calculatePadding(instructions *models.Instructions) (string, string) {
	length := findLongestInstructionNameLength(instructions)
	return createPaddingString(length), createPaddingString(length + 2)
}

func findLongestInstructionNameLength(instructions *models.Instructions) int {
	maxNameLen := 0
	for _, v := range instructions.Actions {
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
	actions := instructions.Actions

	enrichedActionsList := make([]*models.Action, 0, len(actions)+2)
	backAction := &models.Action{
		Id: BackActionName,
	}
	enrichedActionsList = append(enrichedActionsList, backAction)

	if len(instructions.Workflows) > 0 {
		workflowsAction := &models.Action{
			Id: WorkflowsActionName,
		}
		enrichedActionsList = append(enrichedActionsList, workflowsAction)
	}

	enrichedActionsList = append(enrichedActionsList, actions...)
	instructions.Actions = enrichedActionsList
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

func prepareInstructionsSearcher(items []*models.Action) func(input string, index int) bool {
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
		Items:             instructions.Actions,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
	}
}
