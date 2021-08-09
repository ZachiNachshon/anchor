package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"

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
			Name: "have valid completion args as the sub-commands",
			Func: HaveValidCompletionArgsAsTheSubCommands,
		},
		{
			Name: "call set logger verbosity",
			Func: CallSetLoggerVerbosity,
		},
		{
			Name: "fail to set logger verbosity",
			Func: FailToSetLoggerVerbosity,
		},
		{
			Name: "run CLI root command successfully",
			Func: RunCliRootCommandSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var ExpectVerbosityOnceFlagIsSet = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		// Initialize verbose flags which is a global variable to default
		verboseFlagValue = false
		var fun = func(l logger.Logger, verbose bool) error {
			return nil
		}
		cmd := NewCommand(ctx, fun)
		cmd.InitFlags()
		if _, err := drivers.CLI().RunCommand(cmd, fmt.Sprintf("--%s", verboseFlagName)); err != nil {
			logger.Fatalf("expected test to succeed. error: %s", err.Error())
		} else {
			assert.True(t, verboseFlagValue)
		}
	})
}

var ContainExpectedSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(l logger.Logger, verbose bool) error {
			return nil
		}
		cmd := NewCommand(ctx, fun)
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
}

var InitFlagsAndSubCommandsUponInitialization = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		// Initialize verbose flags which is a global variable to default
		verboseFlagValue = false
		var fun = func(l logger.Logger, verbose bool) error {
			return nil
		}
		cmd := NewCommand(ctx, fun)
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
}

var HaveValidCompletionArgsAsTheSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(l logger.Logger, verbose bool) error {
			return nil
		}
		cmd := NewCommand(ctx, fun)
		cmd.InitSubCommands()
		assert.NotNil(t, cmd.cobraCmd.ValidArgs)
		assert.True(t, cmd.cobraCmd.HasSubCommands())
		assert.Equal(t, len(cmd.cobraCmd.Commands()), len(cmd.cobraCmd.ValidArgs))
		for _, subCmd := range cmd.cobraCmd.Commands() {
			assert.Contains(t, cmd.cobraCmd.ValidArgs, subCmd.Use)
		}
	})
}

var CallSetLoggerVerbosity = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		setLoggerCallCount := 0
		var fun = func(l logger.Logger, verbose bool) error {
			setLoggerCallCount++
			return nil
		}
		_, err := drivers.CLI().RunCommand(NewCommand(ctx, fun))
		assert.Equal(t, 1, setLoggerCallCount, "expected action to be called exactly once. name: set-logger-verbosity")
		assert.Nil(t, err, "expected cli action to have no errors")
	})
}

var FailToSetLoggerVerbosity = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		setLoggerCallCount := 0
		var fun = func(l logger.Logger, verbose bool) error {
			setLoggerCallCount++
			return fmt.Errorf("failed to set verbosity")
		}
		_, err := drivers.CLI().RunCommand(NewCommand(ctx, fun))
		assert.Equal(t, 1, setLoggerCallCount, "expected action to be called exactly once. name: set-logger-verbosity")
		assert.NotNil(t, err, "expected cli action to fail")
		assert.Equal(t, "failed to set verbosity", err.Error())
	})
}

var RunCliRootCommandSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				err := RunCliRootCommand(ctx)
				assert.Nil(t, err, "expected cli action to succeed")
			})
		})
	})
}
