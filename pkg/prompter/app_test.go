package prompter

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AppsPrompterShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "apps: append custom prompt options",
			Func: AppsAppendCustomPromptOptions,
		},
		{
			Name: "apps: set search prompt",
			Func: AppsSetSearchPrompt,
		},
		{
			Name: "apps: prepare template",
			Func: AppsPrepareTemplate,
		},
		{
			Name: "apps: prepare searcher",
			Func: AppsPrepareSearcher,
		},
		{
			Name: "apps: prepare full prompter",
			Func: AppsPrepareFullPrompter,
		},
	}
	harness.RunTests(t, tests)
}

var AppsAppendCustomPromptOptions = func(t *testing.T) {
	appTestData := stubs.GenerateApplicationTestData()
	result := appendAppsCustomOptions(appTestData)
	assert.EqualValues(t, CancelActionName, result[0].Name)
}

var AppsSetSearchPrompt = func(t *testing.T) {
	oldSearchPrompt := promptui.SearchPrompt
	setSearchAppPrompt()
	newSearchPrompt := promptui.SearchPrompt
	promptui.SearchPrompt = oldSearchPrompt
	assert.NotEmpty(t, newSearchPrompt)
	assert.Contains(t, newSearchPrompt, "Search:")
}

var AppsPrepareTemplate = func(t *testing.T) {
	template := prepareAppsTemplate()
	assert.NotNil(t, template)
	assert.Contains(t, template.Active, selectorEmoji)
	assert.Contains(t, template.Selected, selectorEmoji)
	assert.NotEmpty(t, template.Details, "expected details to exist")
}

var AppsPrepareSearcher = func(t *testing.T) {
	appTestData := stubs.GenerateApplicationTestData()
	searcherFunc := prepareAppsSearcher(appTestData)
	assert.NotNil(t, searcherFunc)
	found := searcherFunc("app-1", 0)
	assert.True(t, found)
	notFound := searcherFunc("123-app-1", 0)
	assert.False(t, notFound)
}

var AppsPrepareFullPrompter = func(t *testing.T) {
	appTestData := stubs.GenerateApplicationTestData()
	selector := preparePromptAppsItems(appTestData)
	assert.NotNil(t, selector)
	assert.Equal(t, selector.Label, "")
	assert.Equal(t, selector.Size, 10)
	assert.Equal(t, 3, len(selector.Items.([]*models.ApplicationInfo))) // + cancel button
	assert.Equal(t, selector.StartInSearchMode, true)
	assert.Equal(t, selector.HideSelected, true)
}
