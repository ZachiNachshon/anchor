package extractor

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
)

var CreateFakeExtractor = func() *fakeExtractorImpl {
	return &fakeExtractorImpl{}
}

type fakeExtractorImpl struct {
	Extractor
	ExtractPromptItemsMock func(instructionsPath string, p parser.Parser) (*models.Instructions, error)
}

func (e *fakeExtractorImpl) ExtractPromptItems(instructionsPath string, p parser.Parser) (*models.Instructions, error) {
	return e.ExtractPromptItemsMock(instructionsPath, p)
}
