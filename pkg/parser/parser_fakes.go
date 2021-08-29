package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
)

var CreateFakeParser = func() *fakeParserImpl {
	return &fakeParserImpl{}
}

type fakeParserImpl struct {
	Parser
	ParseInstructionsMock func(text string) (*models.InstructionsRoot, error)
}

func (p *fakeParserImpl) ParseInstructions(text string) (*models.InstructionsRoot, error) {
	return p.ParseInstructionsMock(text)
}
