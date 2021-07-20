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
			Name: "add cancel option to apps prompt selector",
			Func: AddCancelOptionToAppsPromptSelector,
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
			Name: "add workflow option when instructions contain workflows",
			Func: AddWorkflowOptionWhenInstructionsContainWorkflows,
		},
		{
			Name: "do not add workflow option when instructions missing workflows",
			Func: DoNotAddWorkflowOptionWhenInstructionsMissingWorkflows,
		},
		{
			Name: "create a valid instructions actions prompt selector",
			Func: CreateValidInstructionsActionsPromptSelector,
		},
	}
	harness.RunTests(t, tests)
}

var AddCancelOptionToAppsPromptSelector = func(t *testing.T) {
	appTestData := stubs.GenerateApplicationTestData()
	appsSelector := preparePromptAppsItems(appTestData)
	assert.EqualValues(t, CancelActionName, appsSelector.Items.([]*models.ApplicationInfo)[0].Name)
}

var CreateValidAppsPromptSelector = func(t *testing.T) {
	appTestData := stubs.GenerateApplicationTestData()
	appsSelector := preparePromptAppsItems(appTestData)
	assert.EqualValues(t, "", appsSelector.Label)
	assert.EqualValues(t, stubs.App1Name, appsSelector.Items.([]*models.ApplicationInfo)[1].Name)
	assert.EqualValues(t, stubs.App2Name, appsSelector.Items.([]*models.ApplicationInfo)[2].Name)
	assert.EqualValues(t, appsSelector.Templates.Details, appsPromptTemplateDetails)
}

var AddBackOptionToInstructionsPromptSelector = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	instSelector := preparePromptInstructionsActions(instRootTestData.Instructions)
	assert.EqualValues(t, BackActionName, instSelector.Items.([]*models.Action)[0].Id)
}

var AddWorkflowOptionWhenInstructionsContainWorkflows = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestData()
	instSelector := preparePromptInstructionsActions(instRootTestData.Instructions)
	assert.EqualValues(t, WorkflowsActionName, instSelector.Items.([]*models.Action)[1].Id)
}

var DoNotAddWorkflowOptionWhenInstructionsMissingWorkflows = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestDataWithoutWorkflows()
	instSelector := preparePromptInstructionsActions(instRootTestData.Instructions)
	assert.NotContains(t, WorkflowsActionName, instSelector.Items)
}

var CreateValidInstructionsActionsPromptSelector = func(t *testing.T) {
	instRootTestData := stubs.GenerateInstructionsTestDataWithoutWorkflows()
	instSelector := preparePromptInstructionsActions(instRootTestData.Instructions)
	assert.EqualValues(t, "", instSelector.Label)
	assert.EqualValues(t, stubs.App1Action1Id, instSelector.Items.([]*models.Action)[1].Id)
	assert.EqualValues(t, stubs.App1Action1Title, instSelector.Items.([]*models.Action)[1].Title)
	assert.EqualValues(t, stubs.App1Action1ScriptFile, instSelector.Items.([]*models.Action)[1].File)
	assert.EqualValues(t, stubs.App1Action2Id, instSelector.Items.([]*models.Action)[2].Id)
	assert.EqualValues(t, stubs.App1Action2Title, instSelector.Items.([]*models.Action)[2].Title)
	assert.EqualValues(t, stubs.App1Action2File, instSelector.Items.([]*models.Action)[2].File)
	assert.EqualValues(t, instructionsPromptTemplateDetails, instSelector.Templates.Details)
}
