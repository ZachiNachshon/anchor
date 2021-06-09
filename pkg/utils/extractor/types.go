package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
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
	locate := reg.Get(identifier).(Extractor)
	if locate == nil {
		return nil, fmt.Errorf("failed to retrieve extractor from registry")
	}
	return locate, nil
}

type Extractor interface {
	ExtractPromptItems(instructionsPath string, p parser.Parser) (*parser.Instructions, error)
}
