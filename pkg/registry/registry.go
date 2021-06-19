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

func New() *InjectionsRegistry {
	return &InjectionsRegistry{
		items: make(map[string]interface{}),
	}
}

func Initialize() *InjectionsRegistry {
	// Only a single instance of registry should be used within a process lifecycle
	initOnlyOnce.Do(func() {
		New()
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
