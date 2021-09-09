package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InstructionsPrompterShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "instruction action: generate run message",
			Func: InstructionActionGenerateRunMessage,
		},
		{
			Name: "instructions action: set search prompt",
			Func: InstructionActionSetSearchPrompt,
		},
		{
			Name: "instructions action: calculate padding of longest action name",
			Func: InstructionActionCalculatePaddingOfLongestActionName,
		},
		{
			Name: "instructions action: find longest action name",
			Func: InstructionActionFindLongestActionName,
		},
		{
			Name: "instructions action: prepare template",
			Func: InstructionActionPrepareTemplate,
		},
		{
			Name: "instructions action: prepare searcher",
			Func: InstructionActionPrepareSearcher,
		},
		{
			Name: "instructions action: prepare full prompter",
			Func: InstructionActionPrepareFullPrompter,
		},
		{
			Name: "instructions workflow: set search prompt",
			Func: InstructionActionSetSearchPrompt,
		},
		{
			Name: "instructions workflow: calculate padding of longest action name",
			Func: InstructionWorkflowCalculatePaddingOfLongestActionName,
		},
		{
			Name: "instructions workflow: find longest action name",
			Func: InstructionWorkflowFindLongestActionName,
		},
		{
			Name: "instructions workflow: prepare template",
			Func: InstructionWorkflowPrepareTemplate,
		},
		{
			Name: "instructions workflow: prepare searcher",
			Func: InstructionWorkflowPrepareSearcher,
		},
		{
			Name: "instructions workflow: prepare full prompter",
			Func: InstructionWorkflowPrepareFullPrompter,
		},
	}
	harness.RunTests(t, tests)
}

var InstructionActionGenerateRunMessage = func(t *testing.T) {
	message := GenerateRunInstructionMessage("test-id", "action", "run instruction action from test")
	assert.NotEmpty(t, message)
	assert.Contains(t, message, "test-id", "expected id")
	assert.Contains(t, message, "action", "expected instruction type")
	assert.Contains(t, message, "run instruction action from test", "expected title")
}

var InstructionActionSetSearchPrompt = func(t *testing.T) {
	oldSearchPrompt := promptui.SearchPrompt
	setSearchInstructionsPrompt("test-app-search")
	newSearchPrompt := promptui.SearchPrompt
	promptui.SearchPrompt = oldSearchPrompt
	assert.NotEmpty(t, newSearchPrompt)
	assert.Contains(t, newSearchPrompt, "Search")
	assert.Contains(t, newSearchPrompt, "test-app-search")
}

var InstructionActionCalculatePaddingOfLongestActionName = func(t *testing.T) {
	actions := []*models.Action{
		{Id: "1234567890"}, {Id: "12345"},
	}
	left, right := calculateActionPadding(actions)
	assert.NotEmpty(t, left)
	assert.Contains(t, left, fmt.Sprintf("%vs", 10+selectorEmojiCharLength))
	assert.NotEmpty(t, right)
	assert.Contains(t, right, fmt.Sprintf("%vs", 10+selectorEmojiCharLength+2))
}

var InstructionActionFindLongestActionName = func(t *testing.T) {
	actions := []*models.Action{
		{Id: "1234567890"}, {Id: "12345"},
	}
	length := findLongestInstructionActionNameLength(actions)
	assert.NotEmpty(t, 10, length)
}

var InstructionActionPrepareTemplate = func(t *testing.T) {
	instructions := stubs.GenerateInstructionsTestData()
	activePadding, inactivePadding := calculateActionPadding(instructions.Instructions.Actions)
	template := prepareInstructionsActionTemplate(activePadding, inactivePadding)
	assert.NotNil(t, template)
	assert.Contains(t, template.Active, selectorEmoji)
	assert.Contains(t, template.Selected, selectorEmoji)
	assert.Contains(t, template.Active, BackActionName)
	assert.Contains(t, template.Inactive, BackActionName)
	assert.Contains(t, template.Active, activePadding)
	assert.Contains(t, template.Inactive, inactivePadding)
	assert.NotEmpty(t, template.Details, "expected details to exist")
}

var InstructionActionPrepareSearcher = func(t *testing.T) {
	instructions := stubs.GenerateInstructionsTestData()
	searcherFunc := prepareInstructionsActionSearcher(instructions.Instructions.Actions)
	assert.NotNil(t, searcherFunc)
	found := searcherFunc("app-1-a", 0)
	assert.True(t, found)
	notFound := searcherFunc("123-app-1-a", 0)
	assert.False(t, notFound)
}

var InstructionActionPrepareFullPrompter = func(t *testing.T) {
	instructions := stubs.GenerateInstructionsTestData()
	selector := preparePromptInstructionsActions(instructions.Instructions.Actions)
	assert.NotNil(t, selector)
	assert.Equal(t, selector.Label, "")
	assert.Equal(t, selector.Size, 15)
	assert.Equal(t, 2, len(selector.Items.([]*models.Action))) // + cancel button
	assert.Equal(t, selector.StartInSearchMode, true)
	assert.Equal(t, selector.HideSelected, true)
}

var InstructionWorkflowCalculatePaddingOfLongestActionName = func(t *testing.T) {
	workflows := []*models.Workflow{
		{Id: "1234567890"}, {Id: "12345"},
	}
	left, right := calculateWorkflowPadding(workflows)
	assert.NotEmpty(t, left)
	assert.Contains(t, left, fmt.Sprintf("%vs", 10+selectorEmojiCharLength))
	assert.NotEmpty(t, right)
	assert.Contains(t, right, fmt.Sprintf("%vs", 10+selectorEmojiCharLength+2))
}

var InstructionWorkflowFindLongestActionName = func(t *testing.T) {
	workflows := []*models.Workflow{
		{Id: "1234567890"}, {Id: "12345"},
	}
	length := findLongestInstructionWorkflowNameLength(workflows)
	assert.NotEmpty(t, 10, length)
}

var InstructionWorkflowPrepareTemplate = func(t *testing.T) {
	instructions := stubs.GenerateInstructionsTestData()
	activePadding, inactivePadding := calculateWorkflowPadding(instructions.Instructions.Workflows)
	template := prepareInstructionsWorkflowTemplate(activePadding, inactivePadding)
	assert.NotNil(t, template)
	assert.Contains(t, template.Active, selectorEmoji)
	assert.Contains(t, template.Selected, selectorEmoji)
	assert.Contains(t, template.Active, BackActionName)
	assert.Contains(t, template.Inactive, BackActionName)
	assert.Contains(t, template.Active, activePadding)
	assert.Contains(t, template.Inactive, inactivePadding)
	assert.NotEmpty(t, template.Details, "expected details to exist")
}

var InstructionWorkflowPrepareSearcher = func(t *testing.T) {
	instructions := stubs.GenerateInstructionsTestData()
	searcherFunc := prepareInstructionsWorkflowSearcher(instructions.Instructions.Workflows)
	assert.NotNil(t, searcherFunc)
	found := searcherFunc("app-1-w", 0)
	assert.True(t, found)
	notFound := searcherFunc("123-app-1-w", 0)
	assert.False(t, notFound)
}

var InstructionWorkflowPrepareFullPrompter = func(t *testing.T) {
	instructions := stubs.GenerateInstructionsTestData()
	selector := preparePromptInstructionsWorkflows(instructions.Instructions.Workflows)
	assert.NotNil(t, selector)
	assert.Equal(t, selector.Label, "")
	assert.Equal(t, selector.Size, 10)
	assert.Equal(t, 2, len(selector.Items.([]*models.Workflow))) // + cancel button
	assert.Equal(t, selector.StartInSearchMode, true)
	assert.Equal(t, selector.HideSelected, true)
}
