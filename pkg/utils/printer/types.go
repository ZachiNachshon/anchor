package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Printer interface {
	PrintApplications(apps []*models.AppContent)
}

const (
	identifier string = "printer"
)

func ToRegistry(reg *registry.InjectionsRegistry, locator Printer) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: locator,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Printer, error) {
	p := reg.Get(identifier).(Printer)
	if p == nil {
		return nil, fmt.Errorf("failed to retrieve printer from registry")
	}
	return p, nil
}
