package extractor

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
)

var CreateFakeExtractor = func() *fakeExtractorImpl {
	return &fakeExtractorImpl{}
}

type fakeExtractorImpl struct {
	Extractor
	ExtractCommandFolderInfoMock func(dirPath string, p parser.Parser) (*models.CommandFolderInfo, error)
	ExtractInstructionsMock      func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error)
}

func (e *fakeExtractorImpl) ExtractCommandFolderInfo(dirPath string, p parser.Parser) (*models.CommandFolderInfo, error) {
	return e.ExtractCommandFolderInfoMock(dirPath, p)
}

func (e *fakeExtractorImpl) ExtractInstructions(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
	return e.ExtractInstructionsMock(instructionsPath, p)
}
