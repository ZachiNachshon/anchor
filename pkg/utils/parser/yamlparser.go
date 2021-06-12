package parser

import (
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
)

type yamlParser struct{}

func New() Parser {
	return &yamlParser{}
}

func (p *yamlParser) Parse(yamlText string) (*models.Instructions, error) {
	items := &models.Instructions{}
	if err := converters.UnmarshalYamlToObj(yamlText, items); err != nil {
		return nil, err
	}
	return items, nil
}
