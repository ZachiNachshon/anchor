package prompter

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
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
			Name: "create a valid apps prompt selector",
			Func: CreateValidAppsPromptSelector,
		},
		{
			Name: "add back option to instructions prompt selector",
			Func: AddBackOptionToInstructionsPromptSelector,
		},
	}
	harness.RunTests(t, tests)
}

var AddBackOptionToAppsPromptSelector = func(t *testing.T) {
	appData := generateApplicationTestData()
	appsSelector := prepareAppsItems(appData)
	assert.EqualValues(t, appsSelector.Items.([]*locator.AppContent)[0].Name, "Back")
}

var CreateValidAppsPromptSelector = func(t *testing.T) {
	appData := generateApplicationTestData()
	appsSelector := prepareAppsItems(appData)
	assert.EqualValues(t, appsSelector.Label, "Available Applications")
	assert.EqualValues(t, appsSelector.Items.([]*locator.AppContent)[1].Name, "first-application")
	assert.EqualValues(t, appsSelector.Items.([]*locator.AppContent)[2].Name, "second-application")
	assert.EqualValues(t, appsSelector.Templates.Details, appsPromptTemplateDetails)
}

var AddBackOptionToInstructionsPromptSelector = func(t *testing.T) {
	instSelector := prepareInstructionsItems(FirstAppInstructions)
	assert.EqualValues(t, instSelector.Items.([]*parser.PromptItem)[0].Id, "Back")
}
