package parser

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
)

type yamlParser struct {
	yamlText string
}

func New(yamlText string) Parser {
	return &yamlParser{
		yamlText: yamlText,
	}
}

func (x *yamlParser) Parse(yamlText string) (*PromptItems, error) {
	items := &PromptItems{}
	if err := converters.UnmarshalYamlToObj(yamlText, items); err != nil {
		return nil, err
	}
	return items, nil
}
