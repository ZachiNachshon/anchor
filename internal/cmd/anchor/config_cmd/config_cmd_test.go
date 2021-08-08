package config_cmd

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"

	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_ConfigCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "contain expected sub commands",
			Func: ContainExpectedSubCommands,
		},
		{
			Name: "have valid completion commands as the sub-commands",
			Func: HaveValidCompletionCommandsAsTheSubCommands,
		},
	}
	harness.RunTests(t, tests)
}

var ContainExpectedSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
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

var HaveValidCompletionCommandsAsTheSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				cmd := NewCommand(ctx)
				assert.NotNil(t, cmd.cobraCmd.ValidArgs)
				assert.True(t, cmd.cobraCmd.HasSubCommands())
				assert.Equal(t, len(cmd.cobraCmd.Commands()), len(cmd.cobraCmd.ValidArgs))
				for _, subCmd := range cmd.cobraCmd.Commands() {
					if strings.Contains(subCmd.Use, "set-context-entry") {
						assert.Contains(t, cmd.cobraCmd.ValidArgs, "set-context-entry")
					} else {
						assert.Contains(t, cmd.cobraCmd.ValidArgs, subCmd.Use)
					}
				}
			})
		})
	})
}
