package parser

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Parser interface {
	Parse(text string) (*models.Instructions, error)
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
	locate := reg.Get(identifier).(Parser)
	if locate == nil {
		return nil, fmt.Errorf("failed to retrieve parser from registry")
	}
	return locate, nil
}
