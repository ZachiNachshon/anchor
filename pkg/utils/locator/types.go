package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Locator interface {
	Scan() error
	Applications() []*AppContent
	ApplicationsAsMap() map[string]*AppContent
	Application(name string) *AppContent
	//Print()
}

const (
	identifier string = "locator"
)

func ToRegistry(reg *registry.InjectionsRegistry, locator Locator) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Locator, error) {
	locate := reg.Get(identifier).(Locator)
	if locate == nil {
		return nil, fmt.Errorf("failed to retrieve locator from registry")
	}
	return locate, nil
}
