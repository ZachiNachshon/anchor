package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Locator interface {
	Scan(anchorFilesLocalPath string) error
	Applications() []*models.AppContent
	ApplicationsAsMap() map[string]*models.AppContent
	Application(name string) *models.AppContent
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
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Locator), nil
}
