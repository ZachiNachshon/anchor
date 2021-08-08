package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
)

const (
	Identifier string = "parser"
)

type Parser interface {
	ParseInstructions(text string) (*models.InstructionsRoot, error)
}

type yamlParser struct{}

func New() Parser {
	return &yamlParser{}
}

func (p *yamlParser) ParseInstructions(yamlText string) (*models.InstructionsRoot, error) {
	items := &models.InstructionsRoot{}
	if err := converters.UnmarshalYamlToObj(yamlText, items); err != nil {
		return nil, err
	}
	return items, nil
}
