package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"

	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/manifoldco/promptui"
	"strings"
)

const kubernetesContextPreMessage = "Kubernetes context identified:"

func prepareFailedKubernetesContextMessage(err error) string {
	logger.Errorf("failed to generate kubernetes context prompt. error: %s", err.Error())
	return kubernetesContextPreMessage + fmt.Sprintf("\n\n%s%s%s\n\n",
		colors.Red, "Failed to retrieve Kubernetes cluster information !", colors.Reset)
}

func GenerateKubernetesRunInstructionMessage(s shell.Shell, id string, instType string, title string) string {
	k8sConfigPath, err := s.ExecuteReturnOutput("echo ${KUBECONFIG}")
	if err != nil {
		return prepareFailedKubernetesContextMessage(err) + GenerateRunInstructionMessage(id, instType, title)
	}

	k8sCurrCtx, err := s.ExecuteReturnOutput("kubectl config current-context")
	if err != nil {
		return prepareFailedKubernetesContextMessage(err) + GenerateRunInstructionMessage(id, instType, title)
	}

	// Use --minify flag to see only the configuration information associated with the current context
	k8sClusters, err := s.ExecuteReturnOutput("kubectl config view --minify -o jsonpath='{.clusters[*].name}'")
	if err != nil {
		return prepareFailedKubernetesContextMessage(err) + GenerateRunInstructionMessage(id, instType, title)
	}

	k8sUsers, err := s.ExecuteReturnOutput("kubectl config view --minify  -o jsonpath='{.users[*].name}'")
	if err != nil {
		return prepareFailedKubernetesContextMessage(err) + GenerateRunInstructionMessage(id, instType, title)
	}

	k8sServers, err := s.ExecuteReturnOutput("kubectl config view --minify -o jsonpath='{.clusters[*].cluster.server}'")
	if err != nil {
		return prepareFailedKubernetesContextMessage(err) + GenerateRunInstructionMessage(id, instType, title)
	}

	contextMsg := fmt.Sprintf(`%s

  Config Path......: %s
  Current Context..: %s
  Cluster..........: %s
  User.............: %s
  Server...........: %s

`, kubernetesContextPreMessage, strings.TrimSpace(k8sConfigPath), strings.TrimSpace(k8sCurrCtx),
		strings.TrimSpace(k8sClusters), strings.TrimSpace(k8sUsers), strings.TrimSpace(k8sServers))

	return contextMsg + GenerateRunInstructionMessage(id, instType, title)
}

func GenerateRunInstructionMessage(id string, instType string, title string) string {
	return fmt.Sprintf("Run instruction %s %s'%s'%s %s(%s)%s?",
		instType, colors.Cyan, id, colors.Reset, colors.Purple, title, colors.Reset)
}

func setSearchInstructionsPrompt(commandItemName string) {
	promptui.SearchPrompt = fmt.Sprintf("%sSearch%s %s%s :%s ", colors.Blue, colors.Green, commandItemName, colors.Blue, colors.Reset)
}

func preparePromptInstructionsActions(actions []*models.Action) promptui.Select {
	instTemplate := prepareInstructionsActionTemplate(calculateActionPadding(actions))
	instSearcher := prepareInstructionsActionSearcher(actions)
	return prepareInstructionsActionsSelector(actions, instTemplate, instSearcher)
}

func calculateActionPadding(actions []*models.Action) (string, string) {
	length := findLongestInstructionActionNameLength(actions)
	return createPaddingLeftString(length + selectorEmojiCharLength),
		createPaddingLeftString(length + selectorEmojiCharLength + 2)
}

func findLongestInstructionActionNameLength(actions []*models.Action) int {
	maxNameLen := 0
	for _, v := range actions {
		actionNameLen := len(v.Id)
		if actionNameLen > maxNameLen {
			maxNameLen = actionNameLen
		}
	}
	return maxNameLen
}

func prepareInstructionsActionTemplate(activeSpacePadding string, inactiveSpacePadding string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label: "{{ . }}",
		Active: selectorEmoji + ` {{ printf ` + activeSpacePadding + ` .Id | cyan }}` +
			`{{ if and (not ( eq .Id "` + BackActionName + `")) (not ( eq .Id "` + WorkflowsActionName + `")) }}` +
			`({{ .Title | green }})` +
			`{{ end }}`,
		Inactive: ` {{ printf ` + inactiveSpacePadding + ` .Id | cyan }}` +
			`{{ if and (not ( eq .Id "` + BackActionName + `")) (not ( eq .Id "` + WorkflowsActionName + `")) }}` +
			`({{ .Title | faint }})` +
			`{{ end }}`,
		Selected: selectorEmoji + " {{ .Id | red | cyan }}",
		Details:  instructionsActionPromptTemplateDetails,
	}
}

func prepareInstructionsActionSearcher(items []*models.Action) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := items[index]
		name := strings.Replace(strings.ToLower(item.Id), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareInstructionsActionsSelector(
	actions []*models.Action,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "",
		Items:             actions,
		Templates:         templates,
		Size:              15,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
		Stdout:            &bellSkipper{},
	}
}

func preparePromptInstructionsWorkflows(workflows []*models.Workflow) promptui.Select {
	instTemplate := prepareInstructionsWorkflowTemplate(calculateWorkflowPadding(workflows))
	instSearcher := prepareInstructionsWorkflowSearcher(workflows)
	return prepareInstructionsWorkflowsSelector(workflows, instTemplate, instSearcher)
}

func calculateWorkflowPadding(workflows []*models.Workflow) (string, string) {
	length := findLongestInstructionWorkflowNameLength(workflows)
	return createPaddingLeftString(length + selectorEmojiCharLength),
		createPaddingLeftString(length + selectorEmojiCharLength + 2)
}

func findLongestInstructionWorkflowNameLength(workflows []*models.Workflow) int {
	maxNameLen := 0
	for _, v := range workflows {
		workflowNameLen := len(v.Id)
		if workflowNameLen > maxNameLen {
			maxNameLen = workflowNameLen
		}
	}
	return maxNameLen
}

func prepareInstructionsWorkflowTemplate(activeSpacePadding string, inactiveSpacePadding string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label: "{{ . }}",
		Active: selectorEmoji + ` {{ printf ` + activeSpacePadding + ` .Id | cyan }}` +
			`{{ if and (not ( eq .Id "` + BackActionName + `")) }}` +
			`({{ .Title | green }})` +
			`{{ end }}`,
		Inactive: ` {{ printf ` + inactiveSpacePadding + ` .Id | cyan }}` +
			`{{ if and (not ( eq .Id "` + BackActionName + `")) }}` +
			`({{ .Title | faint }})` +
			`{{ end }}`,
		Selected: selectorEmoji + " {{ .Id | red | cyan }}",
		Details:  instructionsWorkflowPromptTemplateDetails,
	}
}

func prepareInstructionsWorkflowSearcher(items []*models.Workflow) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := items[index]
		name := strings.Replace(strings.ToLower(item.Id), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareInstructionsWorkflowsSelector(
	workflows []*models.Workflow,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "",
		Items:             workflows,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
		Stdout:            &bellSkipper{},
	}
}
