package extractor

import (
	"github.com/ZachiNachshon/anchor/models"
	parser2 "github.com/ZachiNachshon/anchor/pkg/parser"
)

var CreateFakeExtractor = func() *fakeExtractorImpl {
	return &fakeExtractorImpl{}
}

type fakeExtractorImpl struct {
	Extractor
	ExtractPromptItemsMock func(instructionsPath string, p parser2.Parser) (*models.Instructions, error)
}

func (e *fakeExtractorImpl) ExtractPromptItems(instructionsPath string, p parser2.Parser) (*models.Instructions, error) {
	return e.ExtractPromptItemsMock(instructionsPath, p)
}
