package common

import (
	"github.com/ZachiNachshon/anchor/internal/registry"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ContextShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "create valid empty context",
			Func: CreateValidEmptyContext,
		},
		{
			Name: "set config successfully",
			Func: SetConfigSuccessfully,
		},
		{
			Name: "set logger successfully",
			Func: SetLoggerSuccessfully,
		},
		{
			Name: "set anchorfiles path successfully",
			Func: SetAnchorFilesPathSuccessfully,
		},
		{
			Name: "set non command flags default values as expected",
			Func: SetNonCommandFlagsDefaultValuesAsExpected,
		},
	}
	harness.RunTests(t, tests)
}

var CreateValidEmptyContext = func(t *testing.T) {
	reg := registry.New()
	ctx := EmptyAnchorContext(reg)
	assert.Equal(t, reg, ctx.Registry())
	assert.NotNil(t, ctx.GoContext())
}

var SetConfigSuccessfully = func(t *testing.T) {
	reg := registry.New()
	ctx := EmptyAnchorContext(reg)
	dummyConfig := struct {
		Author string
	}{
		Author: "Zachi Nachshon",
	}
	ctx.(ConfigSetter).SetConfig(dummyConfig)
	assert.Equal(t, dummyConfig, ctx.Config())
}

var SetLoggerSuccessfully = func(t *testing.T) {
	reg := registry.New()
	ctx := EmptyAnchorContext(reg)
	dummyLogger := struct {
		Logger string
	}{
		Logger: "fastest logger in the world",
	}
	ctx.(LoggerSetter).SetLogger(dummyLogger)
	assert.Equal(t, dummyLogger, ctx.Logger())
}

var SetAnchorFilesPathSuccessfully = func(t *testing.T) {
	reg := registry.New()
	ctx := EmptyAnchorContext(reg)
	path := "/path/to/anchorfiles"
	ctx.(AnchorFilesPathSetter).SetAnchorFilesPath(path)
	assert.Equal(t, path, ctx.AnchorFilesPath())
}

var SetNonCommandFlagsDefaultValuesAsExpected = func(t *testing.T) {
	reg := registry.New()
	ctx := EmptyAnchorContext(reg)
	assert.False(t, ctx.NonCmdScopedFlags().NoAutoUpdate)
}
