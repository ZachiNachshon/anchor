package registry

import (
	"fmt"
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
		registry = New()
	})
	return registry
}

func (r *InjectionsRegistry) Set(name string, item interface{}) {
	r.register(RegistryTuple{
		Name:  name,
		Value: item,
	})
}

func (r *InjectionsRegistry) Get(name string) interface{} {
	if value, exists := r.items[name]; exists {
		return value
	}
	return nil
}

func (r *InjectionsRegistry) SafeGet(name string) (interface{}, error) {
	item := r.Get(name)
	if item == nil {
		return nil, fmt.Errorf("failed to retrieve from registry. name: %s", name)
	}
	return item, nil
}

func (r *InjectionsRegistry) register(tuple RegistryTuple) *InjectionsRegistry {
	r.items[tuple.Name] = tuple.Value
	return r
}
