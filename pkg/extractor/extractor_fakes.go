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
	ExtractAnchorFolderInfoMock func(dirPath string, p parser.Parser) (*models.AnchorFolderInfo, error)
	ExtractInstructionsMock     func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error)
}

func (e *fakeExtractorImpl) ExtractAnchorFolderInfo(dirPath string, p parser.Parser) (*models.AnchorFolderInfo, error) {
	return e.ExtractAnchorFolderInfoMock(dirPath, p)
}

func (e *fakeExtractorImpl) ExtractInstructions(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
	return e.ExtractInstructionsMock(instructionsPath, p)
}
