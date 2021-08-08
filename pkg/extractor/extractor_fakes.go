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
	ExtractInstructionsMock func(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error)
}

func (e *fakeExtractorImpl) ExtractInstructions(instructionsPath string, p parser.Parser) (*models.InstructionsRoot, error) {
	return e.ExtractInstructionsMock(instructionsPath, p)
}
