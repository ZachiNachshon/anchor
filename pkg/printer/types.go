package printer

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Printer interface {
	PrintApplications(apps []*models.AppContent)
	PrintConfiguration(ctx common.Context, cfgFilePath string, cfgText string)
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
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Printer), nil
}
