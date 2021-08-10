package app

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"

	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AppCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "load repo or fail successfully",
			Func: LoadRepoOrFailSuccessfully,
		},
		{
			Name: "contain expected sub commands",
			Func: ContainExpectedSubCommands,
		},
	}
	harness.RunTests(t, tests)
}

var LoadRepoOrFailSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context) error {
					callCount++
					return nil
				}
				command, err := NewCommand(ctx, fun)
				_, err = drivers.CLI().RunCommand(command)
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: loadRepoOrFail")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var ContainExpectedSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				var fun = func(ctx common.Context) error {
					return nil
				}
				command, _ := NewCommand(ctx, fun)
				assert.True(t, command.cobraCmd.HasSubCommands())
				cmds := command.cobraCmd.Commands()
				assert.Equal(t, 2, len(cmds))
				assert.Equal(t, "select", cmds[0].Use)
				assert.Equal(t, "status", cmds[1].Use)
			})
		})
	})
}
