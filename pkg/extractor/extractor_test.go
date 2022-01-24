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

func prepareCommandFolderTestPath() string {
	path, _ := ioutils.GetWorkingDirectory()
	return fmt.Sprintf("%s/%s/%s", ioutils.GetRepositoryAbsoluteRootPath(path),
		anchorfilesTestRelativePath, "app")
}

func prepareInstructionTestFilePath(appName string) string {
	path, _ := ioutils.GetWorkingDirectory()
	return fmt.Sprintf("%s/%s/app/%s/%s", ioutils.GetRepositoryAbsoluteRootPath(path),
		anchorfilesTestRelativePath, appName, globals.InstructionsFileName)
}

func Test_ExtractorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "anchor folder info: fail to extract on invalid path",
			Func: CommandFolderInfoFailedToExtractOnInvalidPath,
		},
		{
			Name: "anchor folder info: fail to parse extracted command",
			Func: CommandFolderInfoFailedToParseExtractedCommand,
		},
		{
			Name: "anchor folder info: extract command successfully",
			Func: CommandFolderInfoExtractCommandSuccessfully,
		},
		{
			Name: "anchor folder info: missing command attributes",
			Func: CommandFolderInfoMissingCommandAttributes,
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
		{
			Name: "instructions: extract mixed id and displayName instructions successfully",
			Func: InstructionsExtractMixedIdAndDisplayNameInstructionsSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var CommandFolderInfoFailedToExtractOnInvalidPath = func(t *testing.T) {
	ext := New()
	instRoot, err := ext.ExtractCommandFolderInfo("/invalid/path", parser.New())
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "invalid anchor folder info path")
	assert.Nil(t, instRoot)
}

var CommandFolderInfoFailedToParseExtractedCommand = func(t *testing.T) {
	folderPath := prepareCommandFolderTestPath()
	ext := New()
	fakeParser := parser.CreateFakeParser()
	fakeParser.ParseCommandFolderInfoMock = func(text string) (*models.CommandFolderInfo, error) {
		return nil, fmt.Errorf("failed to parse")
	}
	instRoot, err := ext.ExtractCommandFolderInfo(folderPath, fakeParser)
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "failed to parse")
	assert.Nil(t, instRoot)
}

var CommandFolderInfoExtractCommandSuccessfully = func(t *testing.T) {
	folderPath := prepareCommandFolderTestPath()
	ext := New()
	commandFolders, err := ext.ExtractCommandFolderInfo(folderPath, parser.New())
	assert.Nil(t, err, "expected item extraction to succeed")
	assert.NotNil(t, commandFolders)
	assert.NotEmpty(t, commandFolders.Name)
	assert.NotEmpty(t, commandFolders.Command.Use)
	assert.NotEmpty(t, commandFolders.Command.Short)
	assert.NotEmpty(t, commandFolders.DirPath)
	assert.NotEmpty(t, commandFolders.Description)
	assert.Nil(t, commandFolders.Items)
}

var CommandFolderInfoMissingCommandAttributes = func(t *testing.T) {
	folderPath := prepareCommandFolderTestPath()
	ext := New()
	fakeParser := parser.CreateFakeParser()
	fakeParser.ParseCommandFolderInfoMock = func(text string) (*models.CommandFolderInfo, error) {
		return &models.CommandFolderInfo{
			Command: &models.CommandFolderCommand{
				Use:   "",
				Short: "",
			},
		}, nil
	}
	commandFolders, err := ext.ExtractCommandFolderInfo(folderPath, fakeParser)
	assert.NotNil(t, err, "expected item extraction to fail")
	assert.Equal(t, "bad anchor folder command file structure", err.Error())
	assert.Nil(t, commandFolders)
}

var InstructionsFailedToExtractOnInvalidPath = func(t *testing.T) {
	ext := New()
	instRoot, err := ext.ExtractInstructions("/invalid/path", parser.New())
	assert.NotNil(t, err, "expected extractor to fail")
	assert.Contains(t, err.Error(), "invalid instructions path")
	assert.Nil(t, instRoot)
}

var InstructionsFailedToParseExtractedInstructions = func(t *testing.T) {
	filePath := prepareInstructionTestFilePath("first-app")
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
	filePath := prepareInstructionTestFilePath("second-app")
	ext := New()
	instRoot, err := ext.ExtractInstructions(filePath, parser.New())
	actions := instRoot.Instructions.Actions
	workflows := instRoot.Instructions.Workflows
	assert.Nil(t, err, "expected prompt item extraction to succeed")
	assert.Equal(t, 3, len(actions), "expected 3 actions but found %v", len(actions))
	assert.Equal(t, 2, len(workflows), "expected 2 workflows but found %v", len(workflows))
	assert.Equal(t, "global-hello-universe", actions[0].Id)
	assert.Equal(t, "goodbye-universe", actions[1].Id)
	assert.Equal(t, "hello-universe", actions[2].Id)
	assert.Equal(t, "only-hello", workflows[0].Id)
	assert.Equal(t, "talk-to-the-universe", workflows[1].Id)
}

var InstructionsExtractMixedIdAndDisplayNameInstructionsSuccessfully = func(t *testing.T) {
	filePath := prepareInstructionTestFilePath("mixed-app")
	ext := New()
	instRoot, err := ext.ExtractInstructions(filePath, parser.New())
	actions := instRoot.Instructions.Actions
	workflows := instRoot.Instructions.Workflows
	assert.Nil(t, err, "expected prompt item extraction to succeed")
	assert.Equal(t, 4, len(actions), "expected 4 actions but found %v", len(actions))
	assert.Equal(t, 4, len(workflows), "expected 4 workflows but found %v", len(workflows))
}
