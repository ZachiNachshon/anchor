package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

const instructionsFileName string = "instructions.yaml"
const anchorfilesTestRelativePath string = "test/data/anchorfiles"

func prepareInstructionTestFilePath() string {
	path, _ := ioutils.GetWorkingDirectory()
	return fmt.Sprintf("%s/%s/%s/%s", ioutils.GetRepositoryAbsoluteRootPath(path),
		anchorfilesTestRelativePath, "app/first-app", instructionsFileName)
}

func Test_ExtractorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail to extract on invalid path",
			Func: FailedToExtractOnInvalidPath,
		},
		{
			Name: "fail to parse extracted instructions",
			Func: FailedToParseExtractedInstructions,
		},
		{
			Name: "extract actions from instructions successfully",
			Func: ExtractActionsFromInstructionsSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var FailedToExtractOnInvalidPath = func(t *testing.T) {
	ext := New()
	instRoot, err := ext.ExtractInstructions("/invalid/path", parser.New())
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "invalid instructions path")
	assert.Nil(t, instRoot)
}

var FailedToParseExtractedInstructions = func(t *testing.T) {
	path := prepareInstructionTestFilePath()
	ext := New()
	fakeParser := parser.CreateFakeParser()
	fakeParser.ParseInstructionsMock = func(text string) (*models.InstructionsRoot, error) {
		return nil, fmt.Errorf("failed to parse")
	}
	instRoot, err := ext.ExtractInstructions(path, fakeParser)
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "failed to parse")
	assert.Nil(t, instRoot)
}

var ExtractActionsFromInstructionsSuccessfully = func(t *testing.T) {
	path := prepareInstructionTestFilePath()
	ext := New()
	instRoot, err := ext.ExtractInstructions(path, parser.New())
	actions := instRoot.Instructions.Actions
	assert.Nil(t, err, "expected prompt item extraction to succeed")
	assert.Equal(t, 3, len(actions), "expected 3 instructions but found %v", len(actions))
	// TODO: Rename ids to alphanumeric characters to test ordering
	assert.Equal(t, "global-hello-world", actions[0].Id)
	assert.Equal(t, "goodbye-world", actions[1].Id)
	assert.Equal(t, "hello-world", actions[2].Id)
}
