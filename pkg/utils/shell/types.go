package shell

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Shell interface {
	ExecuteScript(dir string, relativeScriptPath string, args ...string) error
	ExecuteWithOutput(script string) (string, error)
	Execute(script string) error
	ExecuteTTY(script string) error
	ExecuteInBackground(script string) error
}

const (
	identifier string = "shell"
)

func ToRegistry(reg *registry.InjectionsRegistry, shell Shell) {
	reg.Register(registry.RegistryTuple{
		Name:  identifier,
		Value: shell,
	})
}

func FromRegistry(reg *registry.InjectionsRegistry) (Shell, error) {
	locate := reg.Get(identifier).(Shell)
	if locate == nil {
		return nil, fmt.Errorf("failed to retrieve shell executor from registry")
	}
	return locate, nil
}
