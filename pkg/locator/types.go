package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Locator interface {
	Scan(anchorFilesLocalPath string) error
	Applications() []*models.ApplicationInfo
	ApplicationsAsMap() map[string]*models.ApplicationInfo
	Application(name string) *models.ApplicationInfo
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
