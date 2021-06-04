package registry

import (
	"sync"
)

type RegistryTuple struct {
	Name  string
	Value interface{}
}

type InjectionsRegistry struct {
	items map[string]interface{}
}

var (
	initOnlyOnce sync.Once
	registry     *InjectionsRegistry
)

func Initialize() *InjectionsRegistry {
	initOnlyOnce.Do(func() {
		registry = &InjectionsRegistry{
			items: make(map[string]interface{}),
		}
	})
	return registry
}

func (r *InjectionsRegistry) Get(name string) interface{} {
	if value, exists := r.items[name]; exists {
		return value
	}
	return nil
}

func (r *InjectionsRegistry) Register(tuple RegistryTuple) *InjectionsRegistry {
	r.items[tuple.Name] = tuple.Value
	return r
}

//registry.Shell = reg.Shell
//registry.DirLocator = reg.DirLocator
//registry.CmdExtractor = reg.CmdExtractor
//registry.Clipboard = reg.Clipboard
