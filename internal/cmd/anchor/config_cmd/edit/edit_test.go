package edit

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EditActionShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "complete runner method successfully",
			Func: CompleteRunnerMethodSuccessfully,
		},
		{
			Name: "fail runner due to preparation",
			Func: FailRunnerDueToPreparation,
		},
		{
			Name: "prepare registry components",
			Func: PrepareRegistryComponents,
		},
		{
			Name: "fail resolving registry components",
			Func: FailResolvingRegistryComponents,
		},
		{
			Name: "edit configuration file",
			Func: EditConfigurationFile,
		},
	}
	harness.RunTests(t, tests)
}

var CompleteRunnerMethodSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeCfgMgr := config.CreateFakeConfigManager()
			fakeO := NewOrchestrator(fakeCfgMgr)
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *editOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return nil
			}
			runCallCount := 0
			fakeO.runFunc = func(o *editOrchestrator, ctx common.Context) error {
				runCallCount++
				return nil
			}
			err := ConfigEdit(ctx, fakeO)
			assert.Nil(t, err, "expected not to fail")
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
			assert.Equal(t, 1, runCallCount, "expected func to be called exactly once")
		})
	})
}

var FailRunnerDueToPreparation = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			fakeCfgMgr := config.CreateFakeConfigManager()
			fakeO := NewOrchestrator(fakeCfgMgr)
			prepareCallCount := 0
			fakeO.prepareFunc = func(o *editOrchestrator, ctx common.Context) error {
				prepareCallCount++
				return fmt.Errorf("failed to prepare runner")
			}
			err := ConfigEdit(ctx, fakeO)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to prepare runner", err.Error())
			assert.Equal(t, 1, prepareCallCount, "expected func to be called exactly once")
		})
	})
}

var PrepareRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()

		fakeShell := shell.CreateFakeShell()
		reg.Set(shell.Identifier, fakeShell)

		fakeCfgMgr := config.CreateFakeConfigManager()
		fakeO := NewOrchestrator(fakeCfgMgr)
		err := fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
		assert.NotNil(t, fakeO.s)
	})
}

var FailResolvingRegistryComponents = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		fakeCfgMgr := config.CreateFakeConfigManager()
		fakeO := NewOrchestrator(fakeCfgMgr)

		err := fakeO.prepareFunc(fakeO, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", shell.Identifier))

		fakeShell := shell.CreateFakeShell()
		reg.Set(shell.Identifier, fakeShell)

		err = fakeO.prepareFunc(fakeO, ctx)
		assert.Nil(t, err)
	})
}

var EditConfigurationFile = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			configFolderPath := "/path/to/config"
			fakeCfgManager := config.CreateFakeConfigManager()
			fakeCfgManager.GetConfigFilePathMock = func() (string, error) {
				return configFolderPath, nil
			}

			fakeShell := shell.CreateFakeShell()
			executeCallCount := 0
			fakeShell.ExecuteTTYMock = func(script string) error {
				assert.Equal(t, fmt.Sprintf("vi %s/config.yaml", configFolderPath), script)
				executeCallCount++
				return nil
			}

			fakeO := NewOrchestrator(fakeCfgManager)
			fakeO.s = fakeShell

			err := run(fakeO, ctx)
			assert.Nil(t, err, "expected edit to succeed")
			assert.Equal(t, 1, executeCallCount, "expected to be called exactly once")
		})
	})
}
