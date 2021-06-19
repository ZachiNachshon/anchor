package app

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/orchestrator"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InstallActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start the app instruction selection flow successfully",
			Func: StartTheAppInstructionSelectionFlowSuccessfully,
		},
		{
			Name: "fail installation due to invalid registry item",
			Func: FailDueToInvalidRegistryItem,
		},
	}
	harness.RunTests(t, tests)
}

var StartTheAppInstructionSelectionFlowSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			callCount := 0
			o := orchestrator.CreateFakeOrchestrator()
			o.OrchestrateAppInstructionSelectionMock = func() (*models.PromptItem, error) {
				callCount++
				return &models.PromptItem{
					Id:    "testItem",
					Title: "This is a test item",
					File:  "This is a test file path",
				}, nil
			}
			orchestrator.ToRegistry(ctx.Registry(), o)
			err := StartApplicationInstallFlow(ctx)
			assert.Nil(t, err, "expected cli action to have no errors")
			assert.Equal(t, 1, callCount, "expected instruction selection to be called exactly once")
		})
	})
}

var FailDueToInvalidRegistryItem = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			err := StartApplicationInstallFlow(ctx)
			assert.NotNil(t, err, "expected cli action to fail")
			assert.Contains(t, err.Error(), "failed to retrieve from registry")
		})
	})
}

//func populateRegistryWithFakes(ctx common.Context) {
//	registry := ctx.Registry()
//
//	l := locator.CreateFakeLocator(ctx.AnchorFilesPath())
//	locator.ToRegistry(registry, l)
//
//	pr := prompter.CreateFakePrompter()
//	prompter.ToRegistry(registry, pr)
//
//	e := extractor.CreateFakeExtractor()
//	extractor.ToRegistry(registry, e)
//
//	pa := parser.CreateFakeParser()
//	parser.ToRegistry(registry, pa)
//}
