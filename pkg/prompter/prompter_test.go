package prompter

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PrompterShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "config context: fail to prompt",
			Func: ConfigContextFailToPrompt,
		},
		{
			Name: "config context: prompt successfully",
			Func: ConfigContextPromptSuccessfully,
		},
		{
			Name: "config context: cancel selection",
			Func: ConfigContextCancelSelection,
		},
		{
			Name: "anchor folder item: fail to prompt",
			Func: AnchorFolderItemFailToPrompt,
		},
		{
			Name: "anchor folder item: prompt successfully",
			Func: AnchorFolderItemPromptSuccessfully,
		},
		{
			Name: "anchor folder item: cancel selection",
			Func: AnchorFolderItemCancelSelection,
		},
		{
			Name: "actions: fail to prompt",
			Func: ActionsFailToPrompt,
		},
		{
			Name: "actions: prompt successfully",
			Func: ActionsPromptSuccessfully,
		},
		{
			Name: "workflows: fail to prompt",
			Func: WorkflowsFailToPrompt,
		},
		{
			Name: "workflows: prompt successfully",
			Func: WorkflowsPromptSuccessfully,
		},
		{
			Name: "create custom space for text padding",
			Func: CreateCustomSpaceForTextPadding,
		},
	}
	harness.RunTests(t, tests)
}

var ConfigContextFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cfgYamlText := `
config:
 currentContext: test-cfg-ctx-1
 contexts:
   - name: test-cfg-ctx-1
     context:
       repository:
         remote:
           url: /remote/url/one
   - name: test-cfg-ctx-2
     context:
       repository:
         remote:
           url: /remote/url/two
`
			anchorCfg, _ := config.YamlToConfigObj(cfgYamlText)
			prompter := New()
			prompter.runConfigCtxSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", fmt.Errorf("failed to select cfg ctx")
			}
			_, err := prompter.PromptConfigContext(anchorCfg.Config.Contexts)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to select cfg ctx", err.Error())
		})
	})
}

var ConfigContextPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cfgYamlText := `
config:
 currentContext: test-cfg-ctx-1
 contexts:
   - name: test-cfg-ctx-1
     context:
       repository:
         remote:
           url: /remote/url/one
   - name: test-cfg-ctx-2
     context:
       repository:
         remote:
           url: /remote/url/two
`
			anchorCfg, _ := config.YamlToConfigObj(cfgYamlText)
			cfgCtx := config.TryGetConfigContext(anchorCfg.Config.Contexts, "test-cfg-ctx-1")
			prompter := New()
			prompter.runConfigCtxSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 1, "", nil
			}
			selectedCtx, err := prompter.PromptConfigContext(anchorCfg.Config.Contexts)
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, cfgCtx.Name, selectedCtx.Name)
		})
	})
}

var ConfigContextCancelSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			cfgYamlText := `
config:
 currentContext: test-cfg-ctx-1
 contexts:
   - name: test-cfg-ctx-1
     context:
       repository:
         remote:
           url: /remote/url/one
   - name: test-cfg-ctx-2
     context:
       repository:
         remote:
           url: /remote/url/two
`
			anchorCfg, _ := config.YamlToConfigObj(cfgYamlText)
			prompter := New()
			prompter.runConfigCtxSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", nil
			}
			selectedCtx, err := prompter.PromptConfigContext(anchorCfg.Config.Contexts)
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, CancelActionName, selectedCtx.Name)
		})
	})
}

var AnchorFolderItemFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			anchorFolderItems := stubs.GenerateAnchorFolderItemsInfoTestData()
			prompter := New()
			prompter.runAnchorFolderItemSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", fmt.Errorf("failed to select anchor folder item")
			}
			_, err := prompter.PromptAnchorFolderItemSelection(anchorFolderItems)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to select anchor folder item", err.Error())
		})
	})
}

var AnchorFolderItemPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			anchorFolderItems := stubs.GenerateAnchorFolderItemsInfoTestData()
			anchorFolderItems1 := stubs.GetAnchorFolderItemByName(anchorFolderItems, stubs.AnchorFolder1Item1Name)
			prompter := New()
			prompter.runAnchorFolderItemSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 1, "", nil
			}
			selectedAnchorFolderItem, err := prompter.PromptAnchorFolderItemSelection(anchorFolderItems)
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, anchorFolderItems1.Name, selectedAnchorFolderItem.Name)
		})
	})
}

var AnchorFolderItemCancelSelection = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			anchorFolderItems := stubs.GenerateAnchorFolderItemsInfoTestData()
			prompter := New()
			prompter.runAnchorFolderItemSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", nil
			}
			selectedAnchorFolderItem, err := prompter.PromptAnchorFolderItemSelection(anchorFolderItems)
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, CancelActionName, selectedAnchorFolderItem.Name)
		})
	})
}

var ActionsFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			anchorFolderItems := stubs.GenerateAnchorFolderItemsInfoTestData()
			anchorFolderItems1 := stubs.GetAnchorFolderItemByName(anchorFolderItems, stubs.AnchorFolder1Item1Name)
			instData := stubs.GenerateInstructionsTestData()
			prompter := New()
			prompter.runActionSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", fmt.Errorf("failed to select action")
			}
			_, err := prompter.PromptInstructionActions(anchorFolderItems1.Name, instData.Instructions.Actions)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to select action", err.Error())
		})
	})
}

var ActionsPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			anchorFolderItems := stubs.GenerateAnchorFolderItemsInfoTestData()
			anchorFolderItems1 := stubs.GetAnchorFolderItemByName(anchorFolderItems, stubs.AnchorFolder1Item1Name)
			instData := stubs.GenerateInstructionsTestData()
			action1 := stubs.GetInstructionActionById(instData.Instructions, stubs.AnchorFolder1Item1Action1Id)
			prompter := New()
			prompter.runActionSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", nil
			}
			selectedAction, err := prompter.PromptInstructionActions(anchorFolderItems1.Name, instData.Instructions.Actions)
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, action1.Id, selectedAction.Id)
		})
	})
}

var WorkflowsFailToPrompt = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			anchorFolderItems := stubs.GenerateAnchorFolderItemsInfoTestData()
			anchorFolderItems1 := stubs.GetAnchorFolderItemByName(anchorFolderItems, stubs.AnchorFolder1Item1Name)
			instData := stubs.GenerateInstructionsTestData()
			prompter := New()
			prompter.runWorkflowSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", fmt.Errorf("failed to select workflow")
			}
			_, err := prompter.PromptInstructionWorkflows(anchorFolderItems1.Name, instData.Instructions.Workflows)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to select workflow", err.Error())
		})
	})
}

var WorkflowsPromptSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			anchorFolderItems := stubs.GenerateAnchorFolderItemsInfoTestData()
			anchorFolderItems1 := stubs.GetAnchorFolderItemByName(anchorFolderItems, stubs.AnchorFolder1Item1Name)
			instData := stubs.GenerateInstructionsTestData()
			workflow1 := stubs.GetInstructionWorkflowById(instData.Instructions, stubs.AnchorFolder1Item1Workflow1Id)
			prompter := New()
			prompter.runWorkflowSelectorFunc = func(p promptui.Select) (int, string, error) {
				return 0, "", nil
			}
			selectedWorkflow, err := prompter.PromptInstructionWorkflows(anchorFolderItems1.Name, instData.Instructions.Workflows)
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, workflow1.Id, selectedWorkflow.Id)
		})
	})
}

var CreateCustomSpaceForTextPadding = func(t *testing.T) {
	spacesString := createCustomSpacesString(10)
	assert.Equal(t, 10, len(spacesString))
}
