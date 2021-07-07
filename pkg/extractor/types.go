package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

const (
	identifier string = "extractor"
)

func ToRegistry(reg *registry.InjectionsRegistry, locator Extractor) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Extractor, error) {
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Extractor), nil
}

type Extractor interface {
	ExtractPromptItems(instructionsPath string, p parser.Parser) (*models.Instructions, error)
}
