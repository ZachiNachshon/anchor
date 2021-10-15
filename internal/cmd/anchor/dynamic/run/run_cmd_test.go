package run

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/globals"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/drivers"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RunCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "start run action successfully",
			Func: StartRunActionSuccessfully,
		},
		{
			Name: "fail run action",
			Func: FailRunAction,
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
			Name: "add itself to parent command",
			Func: AddItselfToParentCommand,
		},
		{
			Name: "fail on mutual exclusive flags",
			Func: FailOnMutualExclusiveFlags,
		},
		{
			Name: "fail on missing mandatory flags",
			Func: FailOnMissingMandatoryFlags,
		},
		{
			Name: "run in verbose mode if verbose flag supplied",
			Func: RunInVerboseModeIfVerboseFlagSupplied,
		},
		{
			Name: "Run action by id",
			Func: RunActionById,
		},
		{
			Name: "Run workflow by id",
			Func: RunWorkflowById,
		},
		{
			Name: "Fail to initialize flags",
			Func: FailToInitializeFlags,
		},
	}
	harness.RunTests(t, tests)
}

var StartRunActionSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
					callCount++
					return nil
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command, "cmd-item-name", "--action=example")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: run")
				assert.Nil(t, err, "expected cli action to have no errors")
			})
		})
	})
}

var FailRunAction = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				callCount := 0
				var fun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
					callCount++
					return fmt.Errorf("an error occurred")
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command, "cmd-item-name", "--action=example")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once. name: run")
				assert.NotNil(t, err, "expected cli action to fail")
				assert.Contains(t, err.Error(), "an error occurred")
			})
		})
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, stubs.CommandFolder1Name, fun)
		cobraCmd := anchorCmd.GetCobraCmd()
		assert.NotNil(t, cobraCmd, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		var fun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
			return nil
		}
		anchorCmd := NewCommand(ctx, stubs.CommandFolder1Name, fun)
		cmdCtx := anchorCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}

var AddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx, stubs.CommandFolder1Name, nil)
		err := AddCommand(parentCmd, stubs.CommandFolder1Name, NewCommand)
		assert.Nil(t, err, "expected add command to succeed")
		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		cmds := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(cmds))
		assert.Equal(t, "run [COMMAND_ITEM_NAME] [--action=ACTION-ID or --workflow=WORKFLOW-ID]", cmds[0].Use)
	})
}

var FailOnMutualExclusiveFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				commandItemName := "test-command-item-name"
				command := NewCommand(ctx, stubs.CommandFolder1Name, nil)
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command,
					commandItemName,
					fmt.Sprintf("--%s=%s", actionNameFlagName, stubs.CommandFolder1Item1Action1Id),
					fmt.Sprintf("--%s=%s", workflowNameFlagName, stubs.CommandFolder1Item1Workflow1Id),
				)
				assert.NotNil(t, err, "expected command to fail")
				assert.Contains(t, err.Error(), "mutual exclusive flags")
			})
		})
	})
}

var FailOnMissingMandatoryFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				commandItemName := "test-command-item-name"
				command := NewCommand(ctx, stubs.CommandFolder1Name, nil)
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command,
					commandItemName,
				)
				assert.NotNil(t, err, "expected command to fail")
				assert.Contains(t, err.Error(), "missing mandatory flag(s)")
			})
		})
	})
}

var RunInVerboseModeIfVerboseFlagSupplied = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				commandItemName := "test-command-item-name"
				callCount := 0
				var fun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
					callCount++
					assert.True(t, o.verboseFlag, "expected verbose flag to exist")
					return nil
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				flagVal := true
				command.GetCobraCmd().PersistentFlags().BoolVar(&flagVal, globals.VerboseFlagName, true, "")
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command,
					fmt.Sprintf("--%s", globals.VerboseFlagName),
					commandItemName,
					fmt.Sprintf("--%s=%s", actionNameFlagName, stubs.CommandFolder1Item1Action1Id),
				)
				assert.Nil(t, err, "expected command to succeed")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once")
			})
		})
	})
}

var RunActionById = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				commandItemName := "test-command-item-name"
				callCount := 0
				var fun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
					callCount++
					assert.Equal(t, stubs.CommandFolder1Item1Action1Id, identifier)
					return nil
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command,
					commandItemName,
					fmt.Sprintf("--%s=%s", actionNameFlagName, stubs.CommandFolder1Item1Action1Id),
				)
				assert.Nil(t, err, "expected command to succeed")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once")
			})
		})
	})
}

var RunWorkflowById = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(cfg *config.AnchorConfig) {
				commandItemName := "test-command-item-name"
				callCount := 0
				var fun = func(ctx common.Context, o *runOrchestrator, identifier string) error {
					callCount++
					assert.Equal(t, stubs.CommandFolder2Item1Workflow1Id, identifier)
					return nil
				}
				command := NewCommand(ctx, stubs.CommandFolder1Name, fun)
				_ = command.initFlagsFunc(command)
				_, err := drivers.CLI().RunCommand(command,
					commandItemName,
					fmt.Sprintf("--%s=%s", workflowNameFlagName, stubs.CommandFolder2Item1Workflow1Id),
				)
				assert.Nil(t, err, "expected command to succeed")
				assert.Equal(t, 1, callCount, "expected action to be called exactly once")
			})
		})
	})
}

var FailToInitializeFlags = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		parentCmd := NewCommand(ctx, stubs.CommandFolder1Name, nil)
		var newCmdFunc NewCommandFunc = func(ctx common.Context, parentFolderName string, runFunc DynamicRunFunc) *runCmd {
			c := NewCommand(ctx, stubs.CommandFolder1Name, nil)
			c.initFlagsFunc = func(o *runCmd) error {
				return fmt.Errorf("failed to initialize flags")
			}
			return c
		}
		err := AddCommand(parentCmd, stubs.CommandFolder1Name, newCmdFunc)
		assert.NotNil(t, err, "expected add command to fail")
		assert.Equal(t, "failed to initialize flags", err.Error())
	})
}
