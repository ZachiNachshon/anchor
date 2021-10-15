package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
)

var CreateFakeParser = func() *fakeParserImpl {
	return &fakeParserImpl{}
}

type fakeParserImpl struct {
	Parser
	ParseCommandFolderInfoMock func(text string) (*models.CommandFolderInfo, error)
	ParseInstructionsMock      func(text string) (*models.InstructionsRoot, error)
}

func (p *fakeParserImpl) ParseCommandFolderInfo(text string) (*models.CommandFolderInfo, error) {
	return p.ParseCommandFolderInfoMock(text)
}

func (p *fakeParserImpl) ParseInstructions(text string) (*models.InstructionsRoot, error) {
	return p.ParseInstructionsMock(text)
}
