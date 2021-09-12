package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
)

var CreateFakeParser = func() *fakeParserImpl {
	return &fakeParserImpl{}
}

type fakeParserImpl struct {
	Parser
	ParseAnchorFolderInfoMock func(text string) (*models.AnchorFolderInfo, error)
	ParseInstructionsMock     func(text string) (*models.InstructionsRoot, error)
}

func (p *fakeParserImpl) ParseAnchorFolderInfo(text string) (*models.AnchorFolderInfo, error) {
	return p.ParseAnchorFolderInfoMock(text)
}

func (p *fakeParserImpl) ParseInstructions(text string) (*models.InstructionsRoot, error) {
	return p.ParseInstructionsMock(text)
}
