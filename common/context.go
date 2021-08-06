package common

import (
	"context"
	"github.com/ZachiNachshon/anchor/pkg/registry"
)

type Context interface {
	GoContext() context.Context
	Config() interface{}
	Registry() *registry.InjectionsRegistry
	AnchorFilesPath() string
}

type ConfigSetter interface {
	SetConfig(cfg interface{})
}

type RegistryResolver interface {
	Register(registry.InjectionsRegistry)
}

type AnchorFilesPathSetter interface {
	SetAnchorFilesPath(path string)
}

type anchorContext struct {
	goContext                context.Context
	config                   interface{}
	anchorFilesRepoLocalPath string
	registry                 *registry.InjectionsRegistry
}

func (a *anchorContext) GoContext() context.Context {
	return a.goContext
}

func (a *anchorContext) Config() interface{} {
	return a.config
}

func (a *anchorContext) SetConfig(cfg interface{}) {
	a.config = cfg
}

func (a *anchorContext) AnchorFilesPath() string {
	return a.anchorFilesRepoLocalPath
}

func (a *anchorContext) SetAnchorFilesPath(path string) {
	a.anchorFilesRepoLocalPath = path
}

func (a *anchorContext) Registry() *registry.InjectionsRegistry {
	return a.registry
}

func createGoContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return ctx
}

func EmptyAnchorContext(reg *registry.InjectionsRegistry) Context {
	goCtx := createGoContext()
	return &anchorContext{
		goContext: goCtx,
		registry:  reg,
	}
}
