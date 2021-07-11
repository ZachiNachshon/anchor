package shell

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Shell interface {
	ExecuteScriptRealtimeWithOutput(dir string, relativeScriptPath string, args ...string) (string, error)
	ExecuteScript(dir string, relativeScriptPath string, args ...string) error
	ExecuteWithOutput(script string) (string, error)
	Execute(script string) error
	ExecuteSilently(script string) error
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
	item := reg.Get(identifier)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", identifier)
	}
	return item.(Shell), nil
}
