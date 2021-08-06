package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AnchorCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "expect verbosity once flag is set",
			Func: ExpectVerbosityOnceFlagIsSet,
		},
		{
			Name: "contain expected sub commands",
			Func: ContainExpectedSubCommands,
		},
		{
			Name: "init flags and sub-command upon initialization",
			Func: InitFlagsAndSubCommandsUponInitialization,
		},
		{
			Name: "fail on invalid logger verbosity level",
			Func: FailOnInvalidLoggerVerbosityLevel,
		},
	}
	harness.RunTests(t, tests)
}

var ExpectVerbosityOnceFlagIsSet = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				// Initialize verbose flags which is a global variable to default
				verboseFlagValue = false
				cmd := NewCommand(ctx)
				cmd.InitFlags()
				if _, err := drivers.CLI().RunCommand(cmd, fmt.Sprintf("--%s", verboseFlagName)); err != nil {
					logger.Fatalf("expected test to succeed. error: %s", err.Error())
				} else {
					assert.True(t, verboseFlagValue)
				}
			})
		})
	})
}

var ContainExpectedSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				cmd := NewCommand(ctx)
				cmd.InitSubCommands()
				assert.True(t, cmd.cobraCmd.HasSubCommands())
				cmds := cmd.cobraCmd.Commands()
				assert.Equal(t, 6, len(cmds))
				assert.Equal(t, "app", cmds[0].Use)
				assert.Equal(t, "cli", cmds[1].Use)
				assert.Equal(t, "completion", cmds[2].Use)
				assert.Equal(t, "config", cmds[3].Use)
				assert.Equal(t, "controller", cmds[4].Use)
				assert.Equal(t, "version", cmds[5].Use)
			})
		})
	})
}

var InitFlagsAndSubCommandsUponInitialization = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				// Initialize verbose flags which is a global variable to default
				verboseFlagValue = false
				cmd := NewCommand(ctx)
				cmd.Initialize()

				assert.True(t, cmd.cobraCmd.HasPersistentFlags())
				flags := cmd.cobraCmd.PersistentFlags()
				verboseVal, err := flags.GetBool(verboseFlagName)
				assert.Nil(t, err)
				assert.Equal(t, false, verboseVal)

				assert.True(t, cmd.cobraCmd.HasSubCommands())
				cmds := cmd.cobraCmd.Commands()
				assert.Equal(t, 6, len(cmds))
				assert.Equal(t, "app", cmds[0].Use)
				assert.Equal(t, "cli", cmds[1].Use)
				assert.Equal(t, "completion", cmds[2].Use)
				assert.Equal(t, "config", cmds[3].Use)
				assert.Equal(t, "controller", cmds[4].Use)
				assert.Equal(t, "version", cmds[5].Use)
			})
		})
	})
}

var FailOnInvalidLoggerVerbosityLevel = func(t *testing.T) {
	verboseFlagValue = false
	alignLoggerWithVerboseFlag()
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config config.AnchorConfig) {
				// Initialize verbose flags which is a global variable to default
				verboseFlagValue = false
				cmd := NewCommand(ctx)
				cmd.InitFlags()
				if _, err := drivers.CLI().RunCommand(cmd, fmt.Sprintf("--%s", verboseFlagName)); err != nil {
					logger.Fatalf("expected test to succeed. error: %s", err.Error())
				} else {
					assert.True(t, verboseFlagValue)
				}
			})
		})
	})
}
