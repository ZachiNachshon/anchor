package use_context

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UseContextCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start use context action successfully",
			Func: StartUseContextActionSuccessfully,
		},
		{
			Name: "fail due to missing config context name",
			Func: FailDueToMissingConfigContextName,
		},
		{
			Name: "fail use context action",
			Func: FailUseContextAction,
		},
		{
			Name: "contain cobra command",
			Func: ContainCobraCommand,
		},
		{
			Name: "contain context",
			Func: ContainContext,
		},
		{
			Name: "add itself to parent command",
			Func: AddItselfToParentCommand,
		},
	}
	harness.RunTests(t, tests)
}

var StartUseContextActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "test-cfg-context"
				fakeCfgManager := config.CreateFakeConfigManager()
				callCount := 0
				fun := func(ctx common.Context, o *UseContextOrchestrator) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				_, err := drivers.CLI().RunCommand(command, configContextName)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: use-context")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailDueToMissingConfigContextName = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, o *UseContextOrchestrator) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 0, callCount, "expected action not to be called. name: use-context")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, "accepts 1 arg(s), received 0", err.Error())
			})
		})
	})
}

var FailUseContextAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				configContextName := "test-cfg-context"
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, o *UseContextOrchestrator) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				_, err := drivers.CLI().RunCommand(command, configContextName)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: use-context")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		var fun ConfigUseContextFunc = func(ctx common.Context, o *UseContextOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, fakeCfgManager, fun)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		var fun ConfigUseContextFunc = func(ctx common.Context, o *UseContextOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, fakeCfgManager, fun)
		cmdCtx := anchorCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}

var AddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		parentCmd := NewCommand(ctx, fakeCfgManager, nil)
		err := AddCommand(parentCmd, fakeCfgManager, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		cmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(cmds))
		assert.Contains(t, "use-context", cmds[0].Use)
	})
}
