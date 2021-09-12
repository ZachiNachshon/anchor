package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
)

type yamlParser struct{}

func New() Parser {
	return &yamlParser{}
}

func (p *yamlParser) ParseAnchorFolderInfo(yamlText string) (*models.AnchorFolderInfo, error) {
	items := &models.AnchorFolderInfo{}
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
