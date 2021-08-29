package prompter

import (
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConfigContextPrompterShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "config context: append custom prompt options",
			Func: ConfigContextAppendCustomPromptOptions,
		},
		{
			Name: "config context: set search prompt",
			Func: ConfigContextSetSearchPrompt,
		},
		{
			Name: "config context: prepare template",
			Func: ConfigContextPrepareTemplate,
		},
		{
			Name: "config context: prepare searcher",
			Func: ConfigContextPrepareSearcher,
		},
		{
			Name: "config context: prepare full prompter",
			Func: ConfigContextPrepareFullPrompter,
		},
	}
	harness.RunTests(t, tests)
}

var ConfigContextAppendCustomPromptOptions = func(t *testing.T) {
	cfgYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
            url: /remote/url
`
	anchorCfg, _ := config.YamlToConfigObj(cfgYamlText)
	result := appendConfigContextCustomOptions(anchorCfg.Config.Contexts)
	assert.EqualValues(t, CancelActionName, result[0].Name)
}

var ConfigContextSetSearchPrompt = func(t *testing.T) {
	oldSearchPrompt := promptui.SearchPrompt
	setSearchConfigContextPrompt()
	newSearchPrompt := promptui.SearchPrompt
	promptui.SearchPrompt = oldSearchPrompt
	assert.NotEmpty(t, newSearchPrompt)
	assert.Contains(t, newSearchPrompt, "Search:")
}

var ConfigContextPrepareTemplate = func(t *testing.T) {
	template := prepareConfigContextTemplate()
	assert.NotNil(t, template)
	assert.Contains(t, template.Active, selectorEmoji)
	assert.Contains(t, template.Selected, selectorEmoji)
	assert.NotEmpty(t, template.Details, "expected details to exist")
}

var ConfigContextPrepareSearcher = func(t *testing.T) {
	cfgYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
            url: /remote/url
`
	anchorCfg, _ := config.YamlToConfigObj(cfgYamlText)
	searcherFunc := prepareConfigContextSearcher(anchorCfg.Config.Contexts)
	assert.NotNil(t, searcherFunc)
	found := searcherFunc("test-cfg", 0)
	assert.True(t, found)
	notFound := searcherFunc("123-test-cfg", 0)
	assert.False(t, notFound)
}

var ConfigContextPrepareFullPrompter = func(t *testing.T) {
	cfgYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
            url: /remote/url
`
	anchorCfg, _ := config.YamlToConfigObj(cfgYamlText)
	selector := preparePromptConfigContextItems(anchorCfg.Config.Contexts)
	assert.NotNil(t, selector)
	assert.Equal(t, selector.Label, "")
	assert.Equal(t, selector.Size, 5)
	assert.Equal(t, 2, len(selector.Items.([]*config.Context))) // + cancel button
	assert.Equal(t, selector.StartInSearchMode, false)
	assert.Equal(t, selector.HideSelected, true)
}
