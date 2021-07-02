package parser

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Parser interface {
	ParseInstructions(text string) (*models.Instructions, error)
	//Find(text string) string
}

const (
	identifier string = "parser"
)

func ToRegistry(reg *registry.InjectionsRegistry, locator Parser) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Parser, error) {
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Parser), nil
}
