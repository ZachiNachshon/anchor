package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
)

var CreateFakeParser = func() *fakeParserImpl {
	return &fakeParserImpl{}
}

type fakeParserImpl struct {
	Parser
	ParseMock func(text string) (*models.Instructions, error)
}

func (p *fakeParserImpl) ParseInstructions(text string) (*models.Instructions, error) {
	return p.ParseMock(text)
}
