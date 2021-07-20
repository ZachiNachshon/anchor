package parser

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParserShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "parse instruction actions successfully",
			Func: ParseInstructionActionsSuccessfully,
		},
		{
			Name: "parse instruction workflows successfully",
			Func: ParseInstructionWorkflowsSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var ParseInstructionActionsSuccessfully = func(t *testing.T) {
	parser := New()

	instRootTestData, err := parser.ParseInstructions(instructionsOnlyActionsYamlText)
	assert.NotNil(t, instRootTestData, "expected a valid instruction root")
	assert.NotNil(t, instRootTestData.Instructions, "expected a valid instruction object")

	actions := instRootTestData.Instructions.Actions
	assert.Nil(t, err, "expected parser to succeed")
	assert.Equal(t, 3, len(actions), "expected 3 instructions but found %v", len(actions))
	assert.Equal(t, "hello-world", actions[0].Id)
	assert.Equal(t, "goodbye-world", actions[1].Id)
}

var ParseInstructionWorkflowsSuccessfully = func(t *testing.T) {
	parser := New()

	instRootTestData, err := parser.ParseInstructions(instructionsWithWorkflowsYamlText)
	assert.NotNil(t, instRootTestData, "expected a valid instruction root")
	assert.NotNil(t, instRootTestData.Instructions, "expected a valid instruction object")

	workflows := instRootTestData.Instructions.Workflows
	assert.Nil(t, err, "expected parser to succeed")
	assert.Equal(t, 1, len(workflows), "expected 1 workflow but found %v", len(workflows))
	assert.Equal(t, "talk-to-the-world", workflows[0].Id)
}
