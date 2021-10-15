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
			Func: CommandFolderItemsAppendCustomPromptOptions,
		},
		{
			Name: "anchor folder items: set search prompt",
			Func: CommandFolderItemsSetSearchPrompt,
		},
		{
			Name: "anchor folder items: prepare template",
			Func: CommandFolderItemsPrepareTemplate,
		},
		{
			Name: "anchor folder items: prepare searcher",
			Func: CommandFolderItemsPrepareSearcher,
		},
		{
			Name: "anchor folder items: prepare full prompter",
			Func: CommandFolderItemsPrepareFullPrompter,
		},
	}
	harness.RunTests(t, tests)
}

var CommandFolderItemsAppendCustomPromptOptions = func(t *testing.T) {
	folderItemsTestData := stubs.GenerateCommandFolderItemsInfoTestData()
	result := appendFolderItemsCustomOptions(folderItemsTestData)
	assert.EqualValues(t, CancelActionName, result[0].Name)
}

var CommandFolderItemsSetSearchPrompt = func(t *testing.T) {
	oldSearchPrompt := promptui.SearchPrompt
	setSearchFolderItemPrompt()
	newSearchPrompt := promptui.SearchPrompt
	promptui.SearchPrompt = oldSearchPrompt
	assert.NotEmpty(t, newSearchPrompt)
	assert.Contains(t, newSearchPrompt, "Search:")
}

var CommandFolderItemsPrepareTemplate = func(t *testing.T) {
	template := prepareFolderItemsTemplate()
	assert.NotNil(t, template)
	assert.Contains(t, template.Active, selectorEmoji)
	assert.Contains(t, template.Selected, selectorEmoji)
	assert.NotEmpty(t, template.Details, "expected details to exist")
}

var CommandFolderItemsPrepareSearcher = func(t *testing.T) {
	folderItemsTestData := stubs.GenerateCommandFolderItemsInfoTestData()
	searcherFunc := prepareFolderItemsSearcher(folderItemsTestData)
	assert.NotNil(t, searcherFunc)
	found := searcherFunc("app-1", 0)
	assert.True(t, found)
	notFound := searcherFunc("123-app-1", 0)
	assert.False(t, notFound)
}

var CommandFolderItemsPrepareFullPrompter = func(t *testing.T) {
	folderItemsTestData := stubs.GenerateCommandFolderItemsInfoTestData()
	selector := preparePromptFolderItemItems(folderItemsTestData)
	assert.NotNil(t, selector)
	assert.Equal(t, selector.Label, "")
	assert.Equal(t, selector.Size, 15)
	assert.Equal(t, 3, len(selector.Items.([]*models.CommandFolderItemInfo))) // + cancel button
	assert.Equal(t, selector.StartInSearchMode, true)
	assert.Equal(t, selector.HideSelected, true)
}
