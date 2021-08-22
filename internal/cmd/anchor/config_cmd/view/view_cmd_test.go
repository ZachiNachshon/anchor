package view

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

func Test_ViewCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start view action successfully",
			Func: StartViewActionSuccessfully,
		},
		{
			Name: "fail view action",
			Func: FailViewAction,
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

var StartViewActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				fakeCfgManager := config.CreateFakeConfigManager()
				callCount := 0
				fun := func(ctx common.Context, o *viewOrchestrator) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				out, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once")
				assert.Nil(t, err, "expected cli action to succeed")
				logger.Info(out)
			})
		})
	})
}

var FailViewAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				callCount := 0
				fakeCfgManager := config.CreateFakeConfigManager()
				fun := func(ctx common.Context, o *viewOrchestrator) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command := NewCommand(ctx, fakeCfgManager, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once")
				assert.NotNil(t, err, "expected action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgManager := config.CreateFakeConfigManager()
		var fun ConfigViewFunc = func(ctx common.Context, o *viewOrchestrator) error {
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
		var fun ConfigViewFunc = func(ctx common.Context, o *viewOrchestrator) error {
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
		assert.Equal(t, "view", cmds[0].Use)
	})
}
