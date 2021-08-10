package completion

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CliCommandShould(t *testing.T) {
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				cmd, _ := NewCommand(nil, ctx)
				assert.True(t, cmd.cobraCmd.HasSubCommands())
				cmds := cmd.cobraCmd.Commands()
				assert.Equal(t, 2, len(cmds))
				assert.Equal(t, "bash", cmds[0].Use)
				assert.Equal(t, "zsh", cmds[1].Use)
			})
		})
	})
}
