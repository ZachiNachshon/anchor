package config_cmd

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/edit"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/set_context_entry"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/use_context"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/view"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConfigCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
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

var ContainExpectedSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				cmd := NewCommand(ctx)
				assert.True(t, cmd.cobraCmd.HasSubCommands())
				cmds := cmd.cobraCmd.Commands()
				assert.Equal(t, 4, len(cmds))
				assert.Equal(t, "edit", cmds[0].Use)
				assert.Contains(t, cmds[1].Use, "set-context-entry")
				assert.Equal(t, "use-context", cmds[2].Use)
				assert.Equal(t, "view", cmds[3].Use)
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		anchorCmd := NewCommand(ctx)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		anchorCmd := NewCommand(ctx)
		cmdCtx := anchorCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}

var AddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		fakeCfgMgr := config.CreateFakeConfigManager()
		ctx.Registry().Set(config.Identifier, fakeCfgMgr)
		parentCmd := NewCommand(ctx)
		err := AddCommand(parentCmd, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")

		// Parent
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		parentCmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(parentCmds))
		cmdInTest := parentCmds[0]
		assert.Equal(t, "config", cmdInTest.Use)

		// App command
		assert.True(t, cmdInTest.HasSubCommands())
		subCmds := cmdInTest.Commands()
		assert.Equal(t, 4, len(subCmds))
		assert.Equal(t, "edit", subCmds[0].Use)
		assert.Contains(t, subCmds[1].Use, "set-context-entry")
		assert.Equal(t, "use-context", subCmds[2].Use)
		assert.Equal(t, "view", subCmds[3].Use)
	})
}

var FailToAddSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx)

		err := AddCommand(parentCmd, nil)
		assert.NotNil(t, err, "expected to fail on missing registry item")
		fakeCfgMgr := config.CreateFakeConfigManager()
		ctx.Registry().Set(config.Identifier, fakeCfgMgr)

		err = AddCommand(parentCmd, func(ctx common.Context) *configCmd {
			cmdInTest := NewCommand(ctx)
			cmdInTest.addEditSubCmdFunc = func(parent cmd.AnchorCommand, cfgManager config.ConfigManager, createCmd edit.NewCommandFunc) error {
				return fmt.Errorf("failed to add sub command: edit")
			}
			return cmdInTest
		})
		assert.NotNil(t, err, "expected add command to fail on: edit")
		assert.Equal(t, "failed to add sub command: edit", err.Error())

		err = AddCommand(parentCmd, func(ctx common.Context) *configCmd {
			cmdInTest := NewCommand(ctx)
			cmdInTest.addSetContextEntrySubCmdFunc = func(parent cmd.AnchorCommand, cfgManager config.ConfigManager, createCmd set_context_entry.NewCommandFunc) error {
				return fmt.Errorf("failed to add sub command: set-context-entry")
			}
			return cmdInTest
		})
		assert.NotNil(t, err, "expected add command to fail on: set-context-entry")
		assert.Equal(t, "failed to add sub command: set-context-entry", err.Error())

		err = AddCommand(parentCmd, func(ctx common.Context) *configCmd {
			cmdInTest := NewCommand(ctx)
			cmdInTest.addUseContextSubCmdFunc = func(parent cmd.AnchorCommand, cfgManager config.ConfigManager, createCmd use_context.NewCommandFunc) error {
				return fmt.Errorf("failed to add sub command: use-context")
			}
			return cmdInTest
		})
		assert.NotNil(t, err, "expected add command to fail on: use-context")
		assert.Equal(t, "failed to add sub command: use-context", err.Error())

		err = AddCommand(parentCmd, func(ctx common.Context) *configCmd {
			cmdInTest := NewCommand(ctx)
			cmdInTest.addViewSubCmdFunc = func(parent cmd.AnchorCommand, cfgManager config.ConfigManager, createCmd view.NewCommandFunc) error {
				return fmt.Errorf("failed to add sub command: view")
			}
			return cmdInTest
		})
		assert.NotNil(t, err, "expected add command to fail on: view")
		assert.Equal(t, "failed to add sub command: view", err.Error())
	})
}
