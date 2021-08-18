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
			Name: "init flags and sub-command upon initialization",
			Func: InitFlagsAndSubCommandsUponInitialization,
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
		{
			Name: "contain cobra command",
			Func: ContainCobraCommand,
		},
		{
			Name: "contain context",
			Func: ContainContext,
		},
	}
	harness.RunTests(t, tests)
}

var ExpectVerbosityOnceFlagIsSet = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		// Append verbose flags which is a global variable to default
		verboseFlagValue = false
		fakeLoggerManager := logger.CreateFakeLoggerManager()
		setVerbosityCallCount := 0
		fakeLoggerManager.SetVerbosityLevelMock = func(level string) error {
			setVerbosityCallCount++
			return nil
		}
		command := NewCommand(ctx, fakeLoggerManager)
		err := command.initFlagsFunc(command)
		assert.Nil(t, err, "expected init flags to succeed")

		_, err = drivers.CLI().RunCommand(command, fmt.Sprintf("--%s", verboseFlagName))
		assert.Nil(t, err, "expected command to succeed")
		assert.Equal(t, 1, setVerbosityCallCount, "expected func to be called exactly once")
		assert.True(t, verboseFlagValue)
	})
}

var InitFlagsAndSubCommandsUponInitialization = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				// Append verbose flags which is a global variable to default
				verboseFlagValue = false
				fakeLoggerManager := logger.CreateFakeLoggerManager()
				command := NewCommand(ctx, fakeLoggerManager)
				err := command.initialize()
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
		command := NewCommand(ctx, fakeLogManager)
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
		command := NewCommand(ctx, fakeLogManager)
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

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		lgrMgr := logger.CreateFakeLoggerManager()
		newCmd := NewCommand(ctx, lgrMgr)
		cobraCmd := newCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		lgrMgr := logger.CreateFakeLoggerManager()
		newCmd := NewCommand(ctx, lgrMgr)
		cmdCtx := newCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}
