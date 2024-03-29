package parser

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParserShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "parse anchor folder info successfully",
			Func: ParseCommandFolderInfoSuccessfully,
		},
		{
			Name: "fail to parse anchor folder info ",
			Func: FailToParseCommandFolderInfo,
		},
		{
			Name: "parse instruction actions successfully",
			Func: ParseInstructionActionsSuccessfully,
		},
		{
			Name: "parse instruction workflows successfully",
			Func: ParseInstructionWorkflowsSuccessfully,
		},
		{
			Name: "fail to parse instruction",
			Func: FailToParseInstructions,
		},
	}
	harness.RunTests(t, tests)
}

var ParseCommandFolderInfoSuccessfully = func(t *testing.T) {
	parser := New()

	commandFolderInfo, err := parser.ParseCommandFolderInfo(commandFolderInfoYamlText)
	assert.Nil(t, err, "expected parser to succeed")
	assert.NotNil(t, commandFolderInfo, "expected a valid anchor folder info YAML")
	assert.NotEmpty(t, commandFolderInfo.Name, "expected valid attribute: name")
	assert.Empty(t, commandFolderInfo.DirPath, "expected emtpy attribute: dirPath")
	assert.Nil(t, commandFolderInfo.Items, "expected emtpy attribute: items")

	commandFolderCmd := commandFolderInfo.Command
	assert.Equal(t, "app", commandFolderCmd.Use)
	assert.Equal(t, "Application commands", commandFolderCmd.Short)
}

var FailToParseCommandFolderInfo = func(t *testing.T) {
	parser := New()
	invalidYamlText := "@#$%!@#<invalid> yaml: -anchor folder info"
	commandFolderInfo, err := parser.ParseCommandFolderInfo(invalidYamlText)
	assert.NotNil(t, err, "expected to fail")
	assert.Empty(t, commandFolderInfo)
}

var ParseInstructionActionsSuccessfully = func(t *testing.T) {
	parser := New()

	instRootTestData, err := parser.ParseInstructions(instructionsOnlyActionsYamlText)
	assert.Nil(t, err, "expected parser to succeed")
	assert.NotNil(t, instRootTestData, "expected a valid instruction root")
	assert.NotNil(t, instRootTestData.Instructions, "expected a valid instruction object")

	actions := instRootTestData.Instructions.Actions
	assert.Equal(t, 3, len(actions), "expected 3 instructions but found %v", len(actions))
	assert.Equal(t, "hello-world", actions[0].Id)
	assert.Equal(t, "goodbye-world", actions[1].Id)
}

var ParseInstructionWorkflowsSuccessfully = func(t *testing.T) {
	parser := New()

	instRootTestData, err := parser.ParseInstructions(instructionsWithWorkflowsYamlText)
	assert.Nil(t, err, "expected parser to succeed")
	assert.NotNil(t, instRootTestData, "expected a valid instruction root")
	assert.NotNil(t, instRootTestData.Instructions, "expected a valid instruction object")

	workflows := instRootTestData.Instructions.Workflows
	assert.Equal(t, 1, len(workflows), "expected 1 workflow but found %v", len(workflows))
	assert.Equal(t, "talk-to-the-world", workflows[0].Id)
}

var FailToParseInstructions = func(t *testing.T) {
	parser := New()
	invalidYamlText := "@#$%!@#<invalid> yaml: -instructions"
	instRootTestData, err := parser.ParseInstructions(invalidYamlText)
	assert.NotNil(t, err, "expected to fail")
	assert.Empty(t, instRootTestData)
}
