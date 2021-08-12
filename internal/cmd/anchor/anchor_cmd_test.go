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
		fakeLoggerManager := logger.CreateFakeLoggerManager()
		setVerbosityCallCount := 0
		fakeLoggerManager.SetVerbosityLevelMock = func(level string) error {
			setVerbosityCallCount++
			return nil
		}
		command, _ := NewCommand(ctx, fakeLoggerManager)
		err := command.InitFlags()
		assert.Nil(t, err, "expected init flags to succeed")

		_, err = drivers.CLI().RunCommand(command, fmt.Sprintf("--%s", verboseFlagName))
		assert.Nil(t, err, "expected command to succeed")
		assert.Equal(t, 1, setVerbosityCallCount, "expected func to be called exactly once")
		assert.True(t, verboseFlagValue)
	})
}

var ContainExpectedSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				fakeLoggerManager := logger.CreateFakeLoggerManager()
				command, _ := NewCommand(ctx, fakeLoggerManager)
				err := command.InitSubCommands()
				assert.Nil(t, err)
				assert.True(t, command.cobraCmd.HasSubCommands())
				cmds := command.cobraCmd.Commands()
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
		with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
			// Initialize verbose flags which is a global variable to default
			verboseFlagValue = false
			fakeLoggerManager := logger.CreateFakeLoggerManager()
			command, _ := NewCommand(ctx, fakeLoggerManager)
			_, err := command.Initialize()
			assert.Nil(t, err)

			assert.True(t, command.cobraCmd.HasPersistentFlags())
			flags := command.cobraCmd.PersistentFlags()
			verboseVal, err := flags.GetBool(verboseFlagName)
			assert.Nil(t, err)
			assert.Equal(t, false, verboseVal)

			assert.True(t, command.cobraCmd.HasSubCommands())
			cmds := command.cobraCmd.Commands()
			assert.Equal(t, 6, len(cmds))
			assert.Equal(t, "app", cmds[0].Use)
			assert.Equal(t, "cli", cmds[1].Use)
			assert.Equal(t, "completion", cmds[2].Use)
			assert.Equal(t, "config", cmds[3].Use)
			assert.Equal(t, "controller", cmds[4].Use)
			assert.Equal(t, "version", cmds[5].Use)
		})
	})
}

var HaveValidCompletionArgsAsTheSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
			fakeLoggerManager := logger.CreateFakeLoggerManager()
			command, _ := NewCommand(ctx, fakeLoggerManager)
			err := command.InitSubCommands()
			assert.Nil(t, err)
			assert.NotNil(t, command.cobraCmd.ValidArgs)
			assert.True(t, command.cobraCmd.HasSubCommands())
			assert.Equal(t, len(command.cobraCmd.Commands()), len(command.cobraCmd.ValidArgs))
			for _, subCmd := range command.cobraCmd.Commands() {
				assert.Contains(t, command.cobraCmd.ValidArgs, subCmd.Use)
			}
		})
	})
}

var CallSetLoggerVerbosity = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		setVerbosityCallCount := 0
		fakeLogManager := logger.CreateFakeLoggerManager()
		fakeLogManager.SetVerbosityLevelMock = func(level string) error {
			setVerbosityCallCount++
			return nil
		}
		command, _ := NewCommand(ctx, fakeLogManager)
		_, err := drivers.CLI().RunCommand(command)
		assert.Equal(t, 1, setVerbosityCallCount, "expected action to be called exactly once. name: set-logger-verbosity")
		assert.Nil(t, err, "expected cli action to have no errors")
	})
}

var FailToSetLoggerVerbosity = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		setVerbosityCallCount := 0
		fakeLogManager := logger.CreateFakeLoggerManager()
		fakeLogManager.SetVerbosityLevelMock = func(level string) error {
			setVerbosityCallCount++
			return fmt.Errorf("failed to set verbosity")
		}
		command, _ := NewCommand(ctx, fakeLogManager)
		_, err := drivers.CLI().RunCommand(command)
		assert.Equal(t, 1, setVerbosityCallCount, "expected action to be called exactly once. name: set-logger-verbosity")
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
