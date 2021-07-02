package parser

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParserShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "parse prompt items successfully",
			Func: ParsePromptItemsSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var ParsePromptItemsSuccessfully = func(t *testing.T) {
	// Given I create a new YAML parser
	parser := New()
	// When I extract instructions items
	prompItems, err := parser.ParseInstructions(instructionsYamlText)
	// Then I expect to extract exactly the amount of prompt items
	assert.Nil(t, err, "expected parser to succeed")
	assert.Equal(t, 2, len(prompItems.Items), "expected 2 instructions but found %v", len(prompItems.Items))
	// And their names should match
	assert.Equal(t, "hello-world", prompItems.Items[0].Id)
	assert.Equal(t, "goodbye-world", prompItems.Items[1].Id)
	// And I expect a single auto run action
	assert.Equal(t, 1, len(prompItems.AutoRun), "expected 1 auto run action but found %v", len(prompItems.AutoRun))
	// And I expect the auto run action name to match
	assert.Equal(t, "hello-world", prompItems.AutoRun[0])
	// And I expect a single auto cleanup action
	assert.Equal(t, 1, len(prompItems.AutoCleanup), "expected 1 auto cleanup action but found %v", len(prompItems.AutoCleanup))
	// And I expect the auto cleanup action name to match
	assert.Equal(t, "goodbye-world", prompItems.AutoCleanup[0])
}
