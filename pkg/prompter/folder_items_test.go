package prompter

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FolderItemsPrompterShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "anchor folder items: append custom prompt options",
			Func: AnchorFolderItemsAppendCustomPromptOptions,
		},
		{
			Name: "anchor folder items: set search prompt",
			Func: AnchorFolderItemsSetSearchPrompt,
		},
		{
			Name: "anchor folder items: prepare template",
			Func: AnchorFolderItemsPrepareTemplate,
		},
		{
			Name: "anchor folder items: prepare searcher",
			Func: AnchorFolderItemsPrepareSearcher,
		},
		{
			Name: "anchor folder items: prepare full prompter",
			Func: AnchorFolderItemsPrepareFullPrompter,
		},
	}
	harness.RunTests(t, tests)
}

var AnchorFolderItemsAppendCustomPromptOptions = func(t *testing.T) {
	folderItemsTestData := stubs.GenerateAnchorFolderItemsInfoTestData()
	result := appendFolderItemsCustomOptions(folderItemsTestData)
	assert.EqualValues(t, CancelActionName, result[0].Name)
}

var AnchorFolderItemsSetSearchPrompt = func(t *testing.T) {
	oldSearchPrompt := promptui.SearchPrompt
	setSearchFolderItemPrompt()
	newSearchPrompt := promptui.SearchPrompt
	promptui.SearchPrompt = oldSearchPrompt
	assert.NotEmpty(t, newSearchPrompt)
	assert.Contains(t, newSearchPrompt, "Search:")
}

var AnchorFolderItemsPrepareTemplate = func(t *testing.T) {
	template := prepareFolderItemsTemplate()
	assert.NotNil(t, template)
	assert.Contains(t, template.Active, selectorEmoji)
	assert.Contains(t, template.Selected, selectorEmoji)
	assert.NotEmpty(t, template.Details, "expected details to exist")
}

var AnchorFolderItemsPrepareSearcher = func(t *testing.T) {
	folderItemsTestData := stubs.GenerateAnchorFolderItemsInfoTestData()
	searcherFunc := prepareFolderItemsSearcher(folderItemsTestData)
	assert.NotNil(t, searcherFunc)
	found := searcherFunc("app-1", 0)
	assert.True(t, found)
	notFound := searcherFunc("123-app-1", 0)
	assert.False(t, notFound)
}

var AnchorFolderItemsPrepareFullPrompter = func(t *testing.T) {
	folderItemsTestData := stubs.GenerateAnchorFolderItemsInfoTestData()
	selector := preparePromptFolderItemItems(folderItemsTestData)
	assert.NotNil(t, selector)
	assert.Equal(t, selector.Label, "")
	assert.Equal(t, selector.Size, 15)
	assert.Equal(t, 3, len(selector.Items.([]*models.AnchorFolderItemInfo))) // + cancel button
	assert.Equal(t, selector.StartInSearchMode, true)
	assert.Equal(t, selector.HideSelected, true)
}
