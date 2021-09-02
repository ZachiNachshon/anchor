package status

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

func Test_StatusCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start status action successfully",
			Func: StartStatusActionSuccessfully,
		},
		{
			Name: "fail status action",
			Func: FailStatusAction,
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
		{
			Name: "fail on mutual exclusive flags",
			Func: FailOnMutualExclusiveFlags,
		},
		{
			Name: "start status command with all flags",
			Func: StartStatusCommandWithAllFlags,
		},
		{
			Name: "fail to initialize flags",
			Func: FailToInitializeFlags,
		},
	}
	harness.RunTests(t, tests)
}

var StartStatusActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context, o *statusOrchestrator) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: status")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailStatusAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context, o *statusOrchestrator) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command := NewCommand(ctx, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: status")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context, o *statusOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, fun)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context, o *statusOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, fun)
		cmdCtx := anchorCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}

var AddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx, nil)
		err := AddCommand(parentCmd, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		cmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(cmds))
		assert.Equal(t, "status", cmds[0].Use)
	})
}

var FailOnMutualExclusiveFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		command := NewCommand(ctx, nil)
		_ = command.initFlagsFunc(command)
		_, err := drivers.CLI().RunCommand(command,
			fmt.Sprintf("--%s", validStatusOnlyFlagName),
			fmt.Sprintf("--%s", invalidStatusOnlyFlagName))

		assert.NotNil(t, err, "expected cli action to fail")
		assert.Contains(t, err.Error(), "mutual exclusive flags")
	})
}

var StartStatusCommandWithAllFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		command := NewCommand(ctx, nil)
		err := command.initFlagsFunc(command)
		assert.Nil(t, err, "expected flags init to have no errors")
		assert.NotNil(t, command.GetCobraCmd().Flag(validStatusOnlyFlagName), "expected flag to exist")
		assert.NotNil(t, command.GetCobraCmd().Flag(invalidStatusOnlyFlagName), "expected flag to exist")
	})
}

var FailToInitializeFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx, nil)

		var newCmdFunc NewCommandFunc = func(ctx common.Context, statusFunc AppStatusFunc) *statusCmd {
			c := NewCommand(ctx, nil)
			c.initFlagsFunc = func(o *statusCmd) error {
				return fmt.Errorf("failed to initialize flags")
			}
			return c
		}

		err := AddCommand(parentCmd, newCmdFunc)
		assert.NotNil(t, err, "expected add command to fail")
		assert.Equal(t, "failed to initialize flags", err.Error())
	})
}
