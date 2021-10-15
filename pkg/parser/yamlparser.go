package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
)

type yamlParser struct{}

func New() Parser {
	return &yamlParser{}
}

func (p *yamlParser) ParseCommandFolderInfo(yamlText string) (*models.CommandFolderInfo, error) {
	items := &models.CommandFolderInfo{}
	if err := converters.UnmarshalYamlToObj(yamlText, items); err != nil {
		return nil, err
	}
	return items, nil
}

func (p *yamlParser) ParseInstructions(yamlText string) (*models.InstructionsRoot, error) {
	items := &models.InstructionsRoot{}
	if err := converters.UnmarshalYamlToObj(yamlText, items); err != nil {
		return nil, err
	}
	return items, nil
}
