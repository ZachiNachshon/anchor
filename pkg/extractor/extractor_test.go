package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

const instructionsFileName string = "instructions.yaml"
const anchorfilesTestRelativePath string = "test/data/anchorfiles"

func prepareInstructionTestFilePath() string {
	return fmt.Sprintf("%s/%s/%s/%s", ioutils.GetRepositoryAbsoluteRootPath(),
		anchorfilesTestRelativePath, "app/first-app", instructionsFileName)
}

func Test_ExtractorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "extract actions from instructions successfully",
			Func: ExtractActionsFromInstructionsSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var ExtractActionsFromInstructionsSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config config.AnchorConfig) {
				// Given I prepare a valid path to anchorfiles test instructions
				path := prepareInstructionTestFilePath()
				// And I create a new extractor
				ext := New()
				// When I extract prompt items
				instRoot, err := ext.ExtractInstructions(path, parser.New())
				actions := instRoot.Instructions.Actions
				// Then I expect to extract 3 items successfully
				assert.Nil(t, err, "expected prompt item extraction to succeed")
				assert.Equal(t, 3, len(actions), "expected 3 instructions but found %v", len(actions))
				// And their names should match the items from the first-app application
				// TODO: Rename ids to alphanumeric characters to test ordering
				assert.Equal(t, "global-hello-world", actions[0].Id)
				assert.Equal(t, "goodbye-world", actions[1].Id)
				assert.Equal(t, "hello-world", actions[2].Id)
			})
		})
	})
}
