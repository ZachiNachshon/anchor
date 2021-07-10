package input

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type UserInput interface {
	AskYesNoQuestion(question string) (bool, error)
	AskForNumber() (int, error)
	AskForNumberWithDefault() (int, error)
	PressAnyKeyToContinue() error
}

const (
	identifier string = "user-input"
)

func ToRegistry(reg *registry.InjectionsRegistry, in UserInput) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: in,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (UserInput, error) {
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(UserInput), nil
}
