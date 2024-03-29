package common

import (
	"context"
	"github.com/ZachiNachshon/anchor/internal/registry"
)

type Context interface {
	GoContext() context.Context
	Config() interface{}
	Registry() *registry.InjectionsRegistry
	Logger() interface{}
	NonCmdScopedFlags() NonCmdScopedFlags
	AnchorFilesPath() string
}

type ConfigSetter interface {
	SetConfig(cfg interface{})
}

type LoggerSetter interface {
	SetLogger(logger interface{})
}

type AnchorFilesPathSetter interface {
	SetAnchorFilesPath(path string)
}

type NonCmdScopedFlagsSetter interface {
	SetNonCmdScopedFlags(flags NonCmdScopedFlags)
}

type NonCmdScopedFlags struct {
	NoAutoUpdate   bool
	NoAnchorOutput bool
}

type anchorContext struct {
	goContext                context.Context
	config                   interface{}
	registry                 *registry.InjectionsRegistry
	logger                   interface{}
	anchorFilesRepoLocalPath string
	nonCmdScopedFlags        NonCmdScopedFlags
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

func (a *anchorContext) Logger() interface{} {
	return a.logger
}

func (a *anchorContext) SetLogger(logger interface{}) {
	a.logger = logger
}

func (a *anchorContext) AnchorFilesPath() string {
	return a.anchorFilesRepoLocalPath
}

func (a *anchorContext) SetAnchorFilesPath(path string) {
	a.anchorFilesRepoLocalPath = path
}

func (a *anchorContext) NonCmdScopedFlags() NonCmdScopedFlags {
	return a.nonCmdScopedFlags
}

func (a *anchorContext) SetNonCmdScopedFlags(flags NonCmdScopedFlags) {
	a.nonCmdScopedFlags = flags
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
		nonCmdScopedFlags: NonCmdScopedFlags{
			// Default values
			NoAutoUpdate:   false,
			NoAnchorOutput: false,
		},
	}
}
