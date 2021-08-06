package controller

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/root"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ControllerCommandShould(t *testing.T) {
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
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				callCount := 0
				actions := &root.RootCommandActions{
					LoadRepoOrFail: func(ctx common.Context) {
						callCount++
					},
				}
				_, err := drivers.CLI().RunCommand(NewCommand(ctx, actions))
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: loadRepoOrFail")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var ContainExpectedSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				actions := &root.RootCommandActions{}
				cmd := NewCommand(ctx, actions)
				assert.True(t, cmd.cobraCmd.HasSubCommands())
				cmds := cmd.cobraCmd.Commands()
				assert.Equal(t, 1, len(cmds))
				assert.Equal(t, "install", cmds[0].Use)
			})
		})
	})
}
