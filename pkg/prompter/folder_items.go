package prompter

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
	"github.com/manifoldco/promptui"
	"strings"
)

func preparePromptFolderItemItems(folderItems []*models.AnchorFolderItemInfo) promptui.Select {
	folderItemsEnhanced := appendFolderItemsCustomOptions(folderItems)
	folderItemsTemplate := prepareFolderItemsTemplate()
	folderItemsSearcher := prepareFolderItemsSearcher(folderItemsEnhanced)
	return prepareFolderItemsSelector(folderItemsEnhanced, folderItemsTemplate, folderItemsSearcher)
}

func appendFolderItemsCustomOptions(folderItems []*models.AnchorFolderItemInfo) []*models.AnchorFolderItemInfo {
	folderItemsDirs := make([]*models.AnchorFolderItemInfo, 0, len(folderItems)+1)
	cancel := &models.AnchorFolderItemInfo{
		Name: CancelActionName,
	}
	folderItemsDirs = append(folderItemsDirs, cancel)
	folderItemsDirs = append(folderItemsDirs, folderItems...)
	return folderItemsDirs
}

func setSearchFolderItemPrompt() {
	promptui.SearchPrompt = colors.Blue + "Search: " + colors.Reset
}

func prepareFolderItemsTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   selectorEmoji + " {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: selectorEmoji + " {{ .Name | red | cyan }}",
		Details:  folderItemPromptTemplateDetails,
	}
}

func prepareFolderItemsSearcher(folderItems []*models.AnchorFolderItemInfo) func(input string, index int) bool {
	return func(input string, index int) bool {
		item := folderItems[index]
		name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
}

func prepareFolderItemsSelector(
	folderItems []*models.AnchorFolderItemInfo,
	templates *promptui.SelectTemplates,
	searcher func(input string, index int) bool) promptui.Select {

	return promptui.Select{
		Label:             "",
		Items:             folderItems,
		Templates:         templates,
		Size:              15,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
	}
}
