package completion

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion/bash"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion/zsh"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CompletionCommandShould(t *testing.T) {
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
		rootCmd := NewCommand(ctx)
		err := AddCommand(rootCmd, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")

		// Parent
		assert.True(t, rootCmd.GetCobraCmd().HasSubCommands())
		parentCmds := rootCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(parentCmds))
		cmdInTest := parentCmds[0]
		assert.Equal(t, "completion", cmdInTest.Use)

		// App command
		assert.True(t, cmdInTest.HasSubCommands())
		subCmds := cmdInTest.Commands()
		assert.Equal(t, 2, len(subCmds))
		assert.Equal(t, "bash", subCmds[0].Use)
		assert.Equal(t, "zsh", subCmds[1].Use)
	})
}

var FailToAddSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx)

		err := AddCommand(parentCmd, func(ctx common.Context) *completionCmd {
			cmdInTest := NewCommand(ctx)
			cmdInTest.addBashSubCmdFunc = func(root cmd.AnchorCommand, parent cmd.AnchorCommand, createCmd bash.NewCommandFunc) error {
				return fmt.Errorf("failed to add sub command: bash")
			}
			return cmdInTest
		})
		assert.NotNil(t, err, "expected add command to fail on: bash")
		assert.Equal(t, "failed to add sub command: bash", err.Error())

		err = AddCommand(parentCmd, func(ctx common.Context) *completionCmd {
			cmdInTest := NewCommand(ctx)
			cmdInTest.addZshSubCmdFunc = func(root cmd.AnchorCommand, parent cmd.AnchorCommand, createCmd zsh.NewCommandFunc) error {
				return fmt.Errorf("failed to add sub command: zsh")
			}
			return cmdInTest
		})
		assert.NotNil(t, err, "expected add command to fail on: zsh")
		assert.Equal(t, "failed to add sub command: zsh", err.Error())
	})
}
