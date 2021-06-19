package prompter

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/manifoldco/promptui"
	"strings"
)

func preparePromptAppsItems(apps []*models.AppContent) promptui.Select {
	appsEnhanced := appendAppsCustomOptions(apps)
	appsTemplate := prepareAppsTemplate()
	appsSearcher := prepareAppsSearcher(appsEnhanced)
	return prepareAppsSelector(appsEnhanced, appsTemplate, appsSearcher)
}

func appendAppsCustomOptions(apps []*models.AppContent) []*models.AppContent {
	appDirs := make([]*models.AppContent, 0, len(apps)+1)
	cancel := &models.AppContent{
		Name: CancelButtonName,
	}
	appDirs = append(appDirs, cancel)
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

func prepareAppsSearcher(apps []*models.AppContent) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := apps[index]
		name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareAppsSelector(
	apps []*models.AppContent,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "Available Applications",
		Items:             apps,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
	}
}
