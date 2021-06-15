package prompter

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/kits"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PrompterShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "add back option to apps prompt selector",
			Func: AddBackOptionToAppsPromptSelector,
		},
		{
			Name: "create a valid apps prompt selector",
			Func: CreateValidAppsPromptSelector,
		},
		{
			Name: "add back option to instructions prompt selector",
			Func: AddBackOptionToInstructionsPromptSelector,
		},
		{
			Name: "create a valid instructions prompt selector",
			Func: CreateValidInstructionsPromptSelector,
		},
	}
	harness.RunTests(t, tests)
}

var AddBackOptionToAppsPromptSelector = func(t *testing.T) {
	appData := kits.GenerateApplicationTestData()
	appsSelector := preparePromptAppsItems(appData)
	assert.EqualValues(t, "cancel", appsSelector.Items.([]*models.AppContent)[0].Name)
}

var CreateValidAppsPromptSelector = func(t *testing.T) {
	appData := kits.GenerateApplicationTestData()
	appsSelector := preparePromptAppsItems(appData)
	assert.EqualValues(t, "Available Applications", appsSelector.Label)
	assert.EqualValues(t, "first-application", appsSelector.Items.([]*models.AppContent)[1].Name)
	assert.EqualValues(t, "second-application", appsSelector.Items.([]*models.AppContent)[2].Name)
	assert.EqualValues(t, appsSelector.Templates.Details, appsPromptTemplateDetails)
}

var AddBackOptionToInstructionsPromptSelector = func(t *testing.T) {
	instData := kits.GenerateInstructionsTestData()
	instSelector := preparePromptInstructionsItems(instData)
	assert.EqualValues(t, "back", instSelector.Items.([]*models.PromptItem)[0].Id)
}

var CreateValidInstructionsPromptSelector = func(t *testing.T) {
	instData := kits.GenerateInstructionsTestData()
	instSelector := preparePromptInstructionsItems(instData)
	assert.EqualValues(t, "Available Instructions", instSelector.Label)
	assert.EqualValues(t, "hello-world", instSelector.Items.([]*models.PromptItem)[1].Id)
	assert.EqualValues(t, "1st instruction", instSelector.Items.([]*models.PromptItem)[1].Title)
	assert.EqualValues(t, "/path/to/hello-world", instSelector.Items.([]*models.PromptItem)[1].File)
	assert.EqualValues(t, "goodbye-world", instSelector.Items.([]*models.PromptItem)[2].Id)
	assert.EqualValues(t, "2nd instruction", instSelector.Items.([]*models.PromptItem)[2].Title)
	assert.EqualValues(t, "/path/to/goodbye-world", instSelector.Items.([]*models.PromptItem)[2].File)
	assert.EqualValues(t, instructionsPromptTemplateDetails, instSelector.Templates.Details)
}
