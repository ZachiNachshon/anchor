package zsh

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ZshCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "generate zsh completion successfully",
			Func: GenerateZshCompletionSuccessfully,
		},
		{
			Name: "contain cobra command",
			Func: ContainCobraCommand,
		},
		{
			Name: "not contain context",
			Func: NotContainContext,
		},
		{
			Name: "add itself to parent command",
			Func: AddItselfToParentCommand,
		},
	}
	harness.RunTests(t, tests)
}

var GenerateZshCompletionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				var rootCobraCmd = &cobra.Command{
					Use:   "anchor",
					Short: "root cmd",
					Long:  `root cmd`,
				}
				rootWrapper := NewCommand(nil)
				rootWrapper.cobraCmd = rootCobraCmd

				command := NewCommand(rootWrapper)
				_, err := drivers.CLI().RunCommand(command)
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		rootCobraCmd := &cobra.Command{}
		rootWrapper := NewCommand(nil)
		rootWrapper.cobraCmd = rootCobraCmd
		anchorCmd := NewCommand(rootWrapper)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var NotContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		rootCobraCmd := &cobra.Command{}
		rootWrapper := NewCommand(nil)
		rootWrapper.cobraCmd = rootCobraCmd
		rootWrapper.cobraCmd = rootCobraCmd
		anchorCmd := NewCommand(rootWrapper)
		cmdCtx := anchorCmd.GetContext()
		assert.Nil(t, cmdCtx, "expected context not to exist")
	})
}

var AddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		rootWrapper := NewCommand(nil)
		parentCmd := NewCommand(rootWrapper)
		err := AddCommand(rootWrapper, parentCmd, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		cmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(cmds))
		assert.Equal(t, "zsh", cmds[0].Use)
	})
}
