package parser

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type PromptItems struct {
	Items       []PromptItem `yaml:"promptItems"`
	AutoRun     []string     `yaml:"autoRun"`
	AutoCleanup []string     `yaml:"autoCleanup"`
}

type PromptItem struct {
	Id    string `yaml:"id"`
	Title string `yaml:"title"`
	File  string `yaml:"file"`
}

type Parser interface {
	Parse(yamlText string) (*PromptItems, error)
	//Find(text string) string
}

const (
	identifier string = "extractor"
)

func ToRegistry(reg *registry.InjectionsRegistry, locator Parser) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Parser, error) {
	locate := reg.Get(identifier).(Parser)
	if locate == nil {
		return nil, fmt.Errorf("failed to retrieve parser from registry")
	}
	return locate, nil
}
