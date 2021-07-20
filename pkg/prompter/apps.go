package prompter

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/manifoldco/promptui"
	"strings"
)

func preparePromptAppsItems(apps []*models.ApplicationInfo) promptui.Select {
	appsEnhanced := appendAppsCustomOptions(apps)
	appsTemplate := prepareAppsTemplate()
	appsSearcher := prepareAppsSearcher(appsEnhanced)
	return prepareAppsSelector(appsEnhanced, appsTemplate, appsSearcher)
}

func appendAppsCustomOptions(apps []*models.ApplicationInfo) []*models.ApplicationInfo {
	appDirs := make([]*models.ApplicationInfo, 0, len(apps)+1)
	cancel := &models.ApplicationInfo{
		Name: CancelActionName,
	}
	appDirs = append(appDirs, cancel)
	appDirs = append(appDirs, apps...)
	return appDirs
}

func setSearchAppPrompt() {
	promptui.SearchPrompt = colors.Blue + "Search: " + colors.Reset
}

func prepareAppsTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   selectorEmoji + " {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: selectorEmoji + " {{ .Name | red | cyan }}",
		Details:  appsPromptTemplateDetails,
	}
}

func prepareAppsSearcher(apps []*models.ApplicationInfo) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := apps[index]
		name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareAppsSelector(
	apps []*models.ApplicationInfo,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "",
		Items:             apps,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
	}
}
