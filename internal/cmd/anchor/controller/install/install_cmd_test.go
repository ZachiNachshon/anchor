package install

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

func Test_InstallCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start install action successfully",
			Func: StartInstallActionSuccessfully,
		},
		{
			Name: "fail install action",
			Func: FailInstallAction,
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

var StartInstallActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: install")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailInstallAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command := NewCommand(ctx, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: install")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, fun)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context) error {
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
		assert.Equal(t, "install", cmds[0].Use)
	})
}
