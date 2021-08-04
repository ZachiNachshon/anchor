package config

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConfigCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "contain expected sub commands",
			Func: ContainExpectedSubCommands,
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
				assert.Equal(t, len(cmds), 4)
				assert.Equal(t, cmds[0].Use, "edit")
				assert.Contains(t, cmds[1].Use, "set-context-entry")
				assert.Equal(t, cmds[2].Use, "use-context")
				assert.Equal(t, cmds[3].Use, "view")
			})
		})
	})
}
