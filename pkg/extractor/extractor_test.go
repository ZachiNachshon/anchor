package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	parser2 "github.com/ZachiNachshon/anchor/pkg/parser"
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
			Name: "extract prompt items from instructions successfully",
			Func: ExtractPromptItemsFromInstructions,
		},
	}
	harness.RunTests(t, tests)
}

var ExtractPromptItemsFromInstructions = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config config.AnchorConfig) {
				// Given I prepare a valid path to anchorfiles test instructions
				path := prepareInstructionTestFilePath()
				// And I create a new extractor
				ext := New()
				// When I extract prompt items
				prompItems, err := ext.ExtractPromptItems(path, parser2.New())
				// Then I expect to extract 3 items successfully
				assert.Nil(t, err, "expected prompt item extraction to succeed")
				assert.Equal(t, 3, len(prompItems.Items), "expected 3 instructions but found %v", len(prompItems.Items))
				// And their names should match the items from the first-app application
				assert.Equal(t, prompItems.Items[0].Id, "hello-world")
				assert.Equal(t, prompItems.Items[1].Id, "goodbye-world")
				assert.Equal(t, prompItems.Items[2].Id, "global-hello-world")
			})
		})
	})
}
