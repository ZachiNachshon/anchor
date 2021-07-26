package shell

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Shell interface {
	ExecuteScriptFile(dir string, relativeScriptPath string, args ...string) error
	ExecuteScriptFileWithOutputToFile(
		workingDirectory string,
		relativeScriptPath string,
		outputFilePath string,
		args ...string) error

	Execute(script string) error
	ExecuteWithOutputToFile(script string, outputFilePath string) error

	ExecuteWithOutput(script string) (string, error)
	ExecuteSilently(script string) error
	ExecuteTTY(script string) error
	ExecuteInBackground(script string) error
	ClearScreen() error
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
