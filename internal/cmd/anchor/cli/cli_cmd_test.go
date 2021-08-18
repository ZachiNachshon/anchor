package cli

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/cli/versions"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CliCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "run pre run sequence successfully",
			Func: RunPreRunSequenceSuccessfully,
		},
		{
			Name: "fail to run command",
			Func: FailToRunCommand,
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
			Name: "fail to add sub commands",
			Func: FailToAddSubCommands,
		},
	}
	harness.RunTests(t, tests)
}

var RunPreRunSequenceSuccessfully = func(t *testing.T) {
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
				assert.Equal(t, 1, callCount, "expected action to be called exactly once")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailToRunCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return fmt.Errorf("failed to run command")
				}
				command := NewCommand(ctx, fun)
				_, err := drivers.CLI().RunCommand(command)
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Equal(t, "failed to run command", err.Error())
				assert.Equal(t, 1, callCount, "expected action to be called exactly once")
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
		err := AddCommand(parentCmd, &cmd.AnchorCollaborators{}, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")

		// Parent
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		parentCmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(parentCmds))
		cmdInTest := parentCmds[0]
		assert.Equal(t, "cli", cmdInTest.Use)

		// App command
		assert.True(t, cmdInTest.HasSubCommands())
		subCmds := cmdInTest.Commands()
		assert.Equal(t, 1, len(subCmds))
		assert.Equal(t, "versions", subCmds[0].Use)
	})
}

var FailToAddSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx, nil)

		err := AddCommand(parentCmd, &cmd.AnchorCollaborators{}, func(ctx common.Context, preRunSequence cmd.PreRunSequence) *cliCmd {
			cmdInTest := NewCommand(ctx, preRunSequence)
			cmdInTest.addVersionSubCmdFunc = func(parent cmd.AnchorCommand, createCmd versions.NewCommandFunc) error {
				return fmt.Errorf("failed to add sub command: versions")
			}
			return cmdInTest
		})
		assert.NotNil(t, err, "expected add command to fail on: versions")
		assert.Equal(t, "failed to add sub command: versions", err.Error())
	})
}
