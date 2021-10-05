package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

const anchorfilesTestRelativePath string = "test/data/anchorfiles"

func prepareAnchorFolderTestPath() string {
	path, _ := ioutils.GetWorkingDirectory()
	return fmt.Sprintf("%s/%s/%s", ioutils.GetRepositoryAbsoluteRootPath(path),
		anchorfilesTestRelativePath, "app")
}
func prepareInstructionTestFilePath() string {
	path, _ := ioutils.GetWorkingDirectory()
	return fmt.Sprintf("%s/%s/%s/%s", ioutils.GetRepositoryAbsoluteRootPath(path),
		anchorfilesTestRelativePath, "app/first-app", globals.InstructionsFileName)
}

func Test_ExtractorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "anchor folder info: fail to extract on invalid path",
			Func: AnchorFolderInfoFailedToExtractOnInvalidPath,
		},
		{
			Name: "anchor folder info: fail to parse extracted command",
			Func: AnchorFolderInfoFailedToParseExtractedCommand,
		},
		{
			Name: "anchor folder info: extract command successfully",
			Func: AnchorFolderInfoExtractCommandSuccessfully,
		},
		{
			Name: "anchor folder info: missing command attributes",
			Func: AnchorFolderInfoMissingCommandAttributes,
		},
		{
			Name: "instructions: fail to extract on invalid path",
			Func: InstructionsFailedToExtractOnInvalidPath,
		},
		{
			Name: "instructions: fail to parse extracted instructions",
			Func: InstructionsFailedToParseExtractedInstructions,
		},
		{
			Name: "instructions: extract instructions successfully",
			Func: InstructionsExtractInstructionsSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var AnchorFolderInfoFailedToExtractOnInvalidPath = func(t *testing.T) {
	ext := New()
	instRoot, err := ext.ExtractAnchorFolderInfo("/invalid/path", parser.New())
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "invalid anchor folder info path")
	assert.Nil(t, instRoot)
}

var AnchorFolderInfoFailedToParseExtractedCommand = func(t *testing.T) {
	folderPath := prepareAnchorFolderTestPath()
	ext := New()
	fakeParser := parser.CreateFakeParser()
	fakeParser.ParseAnchorFolderInfoMock = func(text string) (*models.AnchorFolderInfo, error) {
		return nil, fmt.Errorf("failed to parse")
	}
	instRoot, err := ext.ExtractAnchorFolderInfo(folderPath, fakeParser)
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "failed to parse")
	assert.Nil(t, instRoot)
}

var AnchorFolderInfoExtractCommandSuccessfully = func(t *testing.T) {
	folderPath := prepareAnchorFolderTestPath()
	ext := New()
	anchorFolders, err := ext.ExtractAnchorFolderInfo(folderPath, parser.New())
	assert.Nil(t, err, "expected item extraction to succeed")
	assert.NotNil(t, anchorFolders)
	assert.NotEmpty(t, anchorFolders.Name)
	assert.NotEmpty(t, anchorFolders.Command.Use)
	assert.NotEmpty(t, anchorFolders.Command.Short)
	assert.NotEmpty(t, anchorFolders.DirPath)
	assert.NotEmpty(t, anchorFolders.Description)
	assert.Nil(t, anchorFolders.Items)
}

var AnchorFolderInfoMissingCommandAttributes = func(t *testing.T) {
	folderPath := prepareAnchorFolderTestPath()
	ext := New()
	fakeParser := parser.CreateFakeParser()
	fakeParser.ParseAnchorFolderInfoMock = func(text string) (*models.AnchorFolderInfo, error) {
		return &models.AnchorFolderInfo{
			Command: &models.AnchorFolderCommand{
				Use:   "",
				Short: "",
			},
		}, nil
	}
	anchorFolders, err := ext.ExtractAnchorFolderInfo(folderPath, fakeParser)
	assert.NotNil(t, err, "expected item extraction to fail")
	assert.Equal(t, "bad anchor folder command file structure", err.Error())
	assert.Nil(t, anchorFolders)
}

var InstructionsFailedToExtractOnInvalidPath = func(t *testing.T) {
	ext := New()
	instRoot, err := ext.ExtractInstructions("/invalid/path", parser.New())
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "invalid instructions path")
	assert.Nil(t, instRoot)
}

var InstructionsFailedToParseExtractedInstructions = func(t *testing.T) {
	filePath := prepareInstructionTestFilePath()
	ext := New()
	fakeParser := parser.CreateFakeParser()
	fakeParser.ParseInstructionsMock = func(text string) (*models.InstructionsRoot, error) {
		return nil, fmt.Errorf("failed to parse")
	}
	instRoot, err := ext.ExtractInstructions(filePath, fakeParser)
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "failed to parse")
	assert.Nil(t, instRoot)
}

var InstructionsExtractInstructionsSuccessfully = func(t *testing.T) {
	filePath := prepareInstructionTestFilePath()
	ext := New()
	instRoot, err := ext.ExtractInstructions(filePath, parser.New())
	actions := instRoot.Instructions.Actions
	workflows := instRoot.Instructions.Workflows
	assert.Nil(t, err, "expected prompt item extraction to succeed")
	assert.Equal(t, 3, len(actions), "expected 3 actions but found %v", len(actions))
	assert.Equal(t, 2, len(workflows), "expected 2 workflows but found %v", len(workflows))
	assert.Equal(t, "global-hello-world", actions[0].Id)
	assert.Equal(t, "goodbye-world", actions[1].Id)
	assert.Equal(t, "hello-world", actions[2].Id)
	assert.Equal(t, "talk-to-the-world", workflows[0].Id)
}
