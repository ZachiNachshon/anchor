package prompter

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
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
			Name: "create valid apps prompt selector",
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
	appData := stubs.GenerateApplicationTestData()
	appsSelector := preparePromptAppsItems(appData)
	assert.EqualValues(t, CancelButtonName, appsSelector.Items.([]*models.AppContent)[0].Name)
}

var CreateValidAppsPromptSelector = func(t *testing.T) {
	appData := stubs.GenerateApplicationTestData()
	appsSelector := preparePromptAppsItems(appData)
	assert.EqualValues(t, "", appsSelector.Label)
	assert.EqualValues(t, stubs.App1Name, appsSelector.Items.([]*models.AppContent)[1].Name)
	assert.EqualValues(t, stubs.App2Name, appsSelector.Items.([]*models.AppContent)[2].Name)
	assert.EqualValues(t, appsSelector.Templates.Details, appsPromptTemplateDetails)
}

var AddBackOptionToInstructionsPromptSelector = func(t *testing.T) {
	instData := stubs.GenerateInstructionsTestData()
	instSelector := preparePromptInstructionsItems(instData)
	assert.EqualValues(t, BackButtonName, instSelector.Items.([]*models.PromptItem)[0].Id)
}

var CreateValidInstructionsPromptSelector = func(t *testing.T) {
	instData := stubs.GenerateInstructionsTestData()
	instSelector := preparePromptInstructionsItems(instData)
	assert.EqualValues(t, "", instSelector.Label)
	assert.EqualValues(t, stubs.App1InstructionsItem1Id, instSelector.Items.([]*models.PromptItem)[1].Id)
	assert.EqualValues(t, stubs.App1InstructionsItem1Title, instSelector.Items.([]*models.PromptItem)[1].Title)
	assert.EqualValues(t, stubs.App1InstructionsItem1File, instSelector.Items.([]*models.PromptItem)[1].File)
	assert.EqualValues(t, stubs.App1InstructionsItem2Id, instSelector.Items.([]*models.PromptItem)[2].Id)
	assert.EqualValues(t, stubs.App1InstructionsItem2Title, instSelector.Items.([]*models.PromptItem)[2].Title)
	assert.EqualValues(t, stubs.App1InstructionsItem2File, instSelector.Items.([]*models.PromptItem)[2].File)
	assert.EqualValues(t, instructionsPromptTemplateDetails, instSelector.Templates.Details)
}
