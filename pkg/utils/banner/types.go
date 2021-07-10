package banner

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Banner interface {
	PrintAnchor()
}

const (
	identifier string = "banner"
)

func ToRegistry(reg *registry.InjectionsRegistry, locator Banner) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Banner, error) {
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Banner), nil
}
