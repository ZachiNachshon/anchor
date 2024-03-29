package _select

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SelectCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start select action successfully",
			Func: StartSelectActionSuccessfully,
		},
		{
			Name: "fail select action",
			Func: FailSelectAction,
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
			Name: "enable command verbosity",
			Func: EnableCommandVerbosity,
		},
	}
	harness.RunTests(t, tests)
}

var StartSelectActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun DynamicSelectFunc = func(ctx common.Context, o *selectOrchestrator) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: select")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailSelectAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context, o *selectOrchestrator) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: select")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context, o *selectOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, stubs.CommandFolder1Name, fun)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context, o *selectOrchestrator) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, stubs.CommandFolder1Name, fun)
		cmdCtx := anchorCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}

var AddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx, stubs.CommandFolder1Name, nil)
		err := AddCommand(parentCmd, stubs.CommandFolder1Name, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		cmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(cmds))
		assert.Equal(t, "select", cmds[0].Use)
	})
}

var EnableCommandVerbosity = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun DynamicSelectFunc = func(ctx common.Context, o *selectOrchestrator) error {
					callCount++
					assert.True(t, o.verboseFlag, "expected verbose flag to exist")
					return nil
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				flagVal := false
				command.GetCobraCmd().PersistentFlags().BoolVar(&flagVal, globals.VerboseFlagName, true, "")
				_, err := drivers.CLI().RunCommand(command, fmt.Sprintf("--%s", globals.VerboseFlagName))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: select")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}
