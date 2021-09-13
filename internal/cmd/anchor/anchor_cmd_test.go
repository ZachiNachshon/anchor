package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/spf13/cobra"
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
			Name: "fail to start pre run sequence",
			Func: FailToStartPreRunSequence,
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
			Name: "fail to run CLI root command due to missing logger from registry",
			Func: FailToRunCliRootCommandDueToMissingLoggerFromRegistry,
		},
		{
			Name: "fail to run CLI root command due to initialization error",
			Func: FailToRunCliRootCommandDueToInitializationError,
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
		{
			Name: "fail on all initialization flows",
			Func: FailOnAllInitializationFlows,
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

		_, err = drivers.CLI().RunCommand(command, fmt.Sprintf("--%s", globals.VerboseFlagName))
		assert.Nil(t, err, "expected command to succeed")
		assert.Equal(t, 1, setVerbosityCallCount, "expected func to be called exactly once")
		assert.True(t, verboseFlagValue)
	})
}

var InitFlagsAndSubCommandsUponInitialization = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				// Set default false to verbose flags since it might be true from another tests (global variable)
				verboseFlagValue = false
				fakeLoggerManager := logger.CreateFakeLoggerManager()
				command := NewCommand(ctx, fakeLoggerManager)

				stubs.GenerateAnchorFolderInfoTestData()
				command.startPreRunSequence = func(parent cmd.AnchorCommand, preRunSequence func(ctx common.Context) error) error {
					return nil
				}
				command.addDynamicSubCommandsFunc = func(parent cmd.AnchorCommand, createCmd dynamic.NewCommandsFunc) error {
					dynamicCmd1 := &cobra.Command{
						Use: "dynamic-cmd-1",
					}
					dynamicCmd2 := &cobra.Command{
						Use: "dynamic-cmd-2",
					}
					parent.GetCobraCmd().AddCommand(dynamicCmd1, dynamicCmd2)
					return nil
				}

				err := command.initialize()
				assert.Nil(t, err)

				assert.True(t, command.cobraCmd.HasPersistentFlags())
				flags := command.cobraCmd.PersistentFlags()
				verboseVal, err := flags.GetBool(globals.VerboseFlagName)
				assert.Nil(t, err)
				assert.Equal(t, false, verboseVal)

				assert.True(t, command.cobraCmd.HasSubCommands())
				cmds := command.cobraCmd.Commands()
				assert.Equal(t, 5, len(cmds))
				assert.Equal(t, "completion", cmds[0].Use)
				assert.Equal(t, "config", cmds[1].Use)
				assert.Equal(t, "dynamic-cmd-1", cmds[2].Use)
				assert.Equal(t, "dynamic-cmd-2", cmds[3].Use)
				assert.Equal(t, "version", cmds[4].Use)
			})
		})
	})
}

var FailToStartPreRunSequence = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			fakeLoggerManager := logger.CreateFakeLoggerManager()
			command := NewCommand(ctx, fakeLoggerManager)
			command.startPreRunSequence = func(parent cmd.AnchorCommand, preRunSequence func(ctx common.Context) error) error {
				return fmt.Errorf("failed to start pre run sequence")
			}
			err := command.initialize()
			assert.NotNil(t, err, "expected to fail on pre run sequence")
			assert.Equal(t, "failed to start pre run sequence", err.Error())
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

var FailToRunCliRootCommandDueToMissingLoggerFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
			err := RunCliRootCommand(ctx)
			assert.NotNil(t, err, "expected to fail")
			assert.Equal(t, "failed to retrieve from registry. name: logger-manager", err.Error())
		})
	})
}

var FailToRunCliRootCommandDueToInitializationError = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				err := RunCliRootCommand(ctx)
				assert.NotNil(t, err, "expected to fail")
			})
		})
	})
}

var RunCliRootCommandSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(lgr logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				// To pass local repo path validation
				cfg.Config.ActiveContext.Context.Repository.Local.Path = ctx.AnchorFilesPath()

				// Do not scan actual repo, use mocks
				fakeLocator := locator.CreateFakeLocator()
				fakeLocator.ScanMock = func(anchorfilesLocalPath string, e extractor.Extractor, pa parser.Parser) *errors.LocatorError {
					return nil
				}
				fakeLocator.AnchorFoldersMock = func() []*models.AnchorFolderInfo {
					return stubs.GenerateAnchorFolderInfoTestData()
				}
				reg := ctx.Registry()
				reg.Set(locator.Identifier, fakeLocator)
				reg.Set(prompter.Identifier, prompter.CreateFakePrompter())
				reg.Set(shell.Identifier, shell.CreateFakeShell())
				reg.Set(extractor.Identifier, extractor.CreateFakeExtractor())
				reg.Set(parser.Identifier, parser.CreateFakeParser())
				err := RunCliRootCommand(ctx)
				assert.Nil(t, err, "expected root command to run successfully")
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

var FailOnAllInitializationFlows = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		lgrMgr := logger.CreateFakeLoggerManager()
		newCmd := NewCommand(ctx, lgrMgr)

		// Fail on initializing flags
		newCmd.initFlagsFunc = func(o *anchorCmd) error {
			return fmt.Errorf("failed to init flags")
		}
		err := newCmd.initialize()
		assert.NotNil(t, err)
		assert.Equal(t, "failed to init flags", err.Error())
		newCmd.initFlagsFunc = func(o *anchorCmd) error {
			return nil
		}

		// Fail on pre run sequence
		newCmd.startPreRunSequence = func(parent cmd.AnchorCommand, preRunSequence func(ctx common.Context) error) error {
			return fmt.Errorf("failed on pre run sequence")
		}
		err = newCmd.initialize()
		assert.NotNil(t, err)
		assert.Equal(t, "failed on pre run sequence", err.Error())
		newCmd.startPreRunSequence = func(parent cmd.AnchorCommand, preRunSequence func(ctx common.Context) error) error {
			return nil
		}

		// Fail on dynamic commands
		newCmd.addDynamicSubCommandsFunc = func(parent cmd.AnchorCommand, createCmd dynamic.NewCommandsFunc) error {
			return fmt.Errorf("failed to add dynamic commands")
		}
		err = newCmd.initialize()
		assert.NotNil(t, err)
		assert.Equal(t, "failed to add dynamic commands", err.Error())
		newCmd.addDynamicSubCommandsFunc = func(parent cmd.AnchorCommand, createCmd dynamic.NewCommandsFunc) error {
			return nil
		}

		// Fail on config command
		newCmd.addConfigSubCmdFunc = func(parent cmd.AnchorCommand, createCmd config_cmd.NewCommandFunc) error {
			return fmt.Errorf("failed to add config subcommand")
		}
		err = newCmd.initialize()
		assert.NotNil(t, err)
		assert.Equal(t, "failed to add config subcommand", err.Error())
		newCmd.addConfigSubCmdFunc = func(parent cmd.AnchorCommand, createCmd config_cmd.NewCommandFunc) error {
			return nil
		}

		// Fail on version command
		newCmd.addVersionSubCmdFunc = func(parent cmd.AnchorCommand, createCmd version.NewCommandFunc) error {
			return fmt.Errorf("failed to add version subcommand")
		}
		err = newCmd.initialize()
		assert.NotNil(t, err)
		assert.Equal(t, "failed to add version subcommand", err.Error())
		newCmd.addVersionSubCmdFunc = func(parent cmd.AnchorCommand, createCmd version.NewCommandFunc) error {
			return nil
		}

		// Fail on completion command
		newCmd.addCompletionSubCmdFunc = func(root cmd.AnchorCommand, createCmd completion.NewCommandFunc) error {
			return fmt.Errorf("failed to add completion subcommand")
		}
		err = newCmd.initialize()
		assert.NotNil(t, err)
		assert.Equal(t, "failed to add completion subcommand", err.Error())
		newCmd.addCompletionSubCmdFunc = func(root cmd.AnchorCommand, createCmd completion.NewCommandFunc) error {
			return nil
		}

		// Succeed eventually when all init steps are valid
		err = newCmd.initialize()
		assert.Nil(t, err)
	})
}
