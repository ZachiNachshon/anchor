package dynamic

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/run"
	_select "github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/select"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/status"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/data/stubs"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AppCommandShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "create new dynamic commands from anchor folders",
			Func: CreateNewDynamicCommandsFromCommandFolders,
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
			Name: "add commands: fail to resolve from registry",
			Func: AddCommandsFailToResolveFromRegistry,
		},
		{
			Name: "add commands: fail on create dynamic commands",
			Func: AddCommandsFailOnCreateDynamicCommands,
		},
		{
			Name: "add commands: fail to add sub commands",
			Func: AddCommandsFailToAddSubCommands,
		},
		{
			Name: "add commands: add itself to parent command",
			Func: AddCommandAddItselfToParentCommand,
		},
	}
	harness.RunTests(t, tests)
}

var CreateNewDynamicCommandsFromCommandFolders = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		commandFolder := stubs.GenerateCommandFolderInfoTestData()
		commands, err := NewCommands(ctx, commandFolder)
		assert.Nil(t, err, "expected to create dynamic commands successfully")
		assert.Equal(t, 2, len(commands))
		firstCmd := commands[0]
		assert.NotNil(t, firstCmd.GetCobraCmd())
		assert.NotEmpty(t, firstCmd.GetCobraCmd().Use)
		assert.NotEmpty(t, firstCmd.GetCobraCmd().Short)
		assert.NotNil(t, firstCmd.GetContext())
		assert.NotNil(t, firstCmd.commandName)
		assert.NotNil(t, firstCmd.addStatusSubCmdFunc)
		assert.NotNil(t, firstCmd.addSelectSubCmdFunc)
	})
}

var ContainCobraCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		commandFolder := stubs.GenerateCommandFolderInfoTestData()
		cobraCmd := newCobraCommand(commandFolder[0])
		anchorCmd := newCommand(ctx, cobraCmd, "test-name")
		result := anchorCmd.GetCobraCmd()
		assert.NotNil(t, result, "expected cobra command to exist")
	})
}

var ContainContext = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		anchorCmd := newCommand(ctx, nil, "test-name")
		cmdCtx := anchorCmd.GetContext()
		assert.NotNil(t, cmdCtx, "expected context to exist")
		assert.Equal(t, ctx, cmdCtx)
	})
}

var AddCommandsFailToResolveFromRegistry = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		parentCmd := newCommand(ctx, nil, "")

		err := AddCommands(parentCmd, nil)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("failed to retrieve from registry. name: %s", locator.Identifier))
		fakeLocator := locator.CreateFakeLocator()
		fakeLocator.CommandFoldersMock = func() []*models.CommandFolderInfo {
			return nil
		}
		reg.Set(locator.Identifier, fakeLocator)

		createCmdFunc := func(ctx common.Context, commandFolders []*models.CommandFolderInfo) ([]*dynamicCmd, error) {
			return []*dynamicCmd{}, nil
		}
		err = AddCommands(parentCmd, createCmdFunc)
		assert.Nil(t, err)
	})
}

var AddCommandsFailOnCreateDynamicCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		fakeLocator := locator.CreateFakeLocator()
		fakeLocator.CommandFoldersMock = func() []*models.CommandFolderInfo {
			return nil
		}
		reg.Set(locator.Identifier, fakeLocator)
		reg.Set(prompter.Identifier, prompter.CreateFakePrompter())
		reg.Set(shell.Identifier, shell.CreateFakeShell())

		parentCmd := newCommand(ctx, nil, "")

		createCmdFunc := func(ctx common.Context, commandFolders []*models.CommandFolderInfo) ([]*dynamicCmd, error) {
			return nil, fmt.Errorf("fail to create dynamic commands")
		}
		err := AddCommands(parentCmd, createCmdFunc)
		assert.NotNil(t, err, "expected to fail on dynamic commands creation")
		assert.Equal(t, "fail to create dynamic commands", err.Error())
	})
}

var AddCommandsFailToAddSubCommands = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		folders := stubs.GenerateCommandFolderInfoTestData()
		parentCmd := newCommand(ctx, &cobra.Command{}, "")

		fakeLocator := locator.CreateFakeLocator()
		fakeLocator.CommandFoldersMock = func() []*models.CommandFolderInfo {
			return folders
		}
		reg.Set(locator.Identifier, fakeLocator)
		reg.Set(prompter.Identifier, prompter.CreateFakePrompter())
		reg.Set(shell.Identifier, shell.CreateFakeShell())

		newCmdFailSelect := &dynamicCmd{
			addSelectSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd _select.NewCommandFunc) error {
				return fmt.Errorf("fail to add sub command: select ")
			},
		}
		newCmdFailStatus := &dynamicCmd{
			addSelectSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd _select.NewCommandFunc) error {
				return nil
			},
			addStatusSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd status.NewCommandFunc) error {
				return fmt.Errorf("fail to add sub command: status")
			},
		}
		newCmdFailRun := &dynamicCmd{
			addSelectSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd _select.NewCommandFunc) error {
				return nil
			},
			addStatusSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd status.NewCommandFunc) error {
				return nil
			},
			addRunSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd run.NewCommandFunc) error {
				return fmt.Errorf("fail to add sub command: run")
			},
		}

		createCmdsFunc := func(ctx common.Context, commandFolders []*models.CommandFolderInfo) ([]*dynamicCmd, error) {
			assert.Equal(t, folders, commandFolders)
			return []*dynamicCmd{newCmdFailSelect}, nil
		}
		err := AddCommands(parentCmd, createCmdsFunc)
		assert.NotNil(t, err, "expected to fail adding sub command")
		assert.Equal(t, "fail to add sub command: select ", err.Error())

		createCmdsFunc = func(ctx common.Context, commandFolders []*models.CommandFolderInfo) ([]*dynamicCmd, error) {
			assert.Equal(t, folders, commandFolders)
			return []*dynamicCmd{newCmdFailStatus}, nil
		}
		err = AddCommands(parentCmd, createCmdsFunc)
		assert.NotNil(t, err, "expected to fail adding sub command")
		assert.Equal(t, "fail to add sub command: status", err.Error())

		createCmdsFunc = func(ctx common.Context, commandFolders []*models.CommandFolderInfo) ([]*dynamicCmd, error) {
			assert.Equal(t, folders, commandFolders)
			return []*dynamicCmd{newCmdFailRun}, nil
		}
		err = AddCommands(parentCmd, createCmdsFunc)
		assert.NotNil(t, err, "expected to fail adding sub command")
		assert.Equal(t, "fail to add sub command: run", err.Error())
	})
}

var AddCommandAddItselfToParentCommand = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		reg := ctx.Registry()
		folders := stubs.GenerateCommandFolderInfoTestData()
		parentCmd := newCommand(ctx, &cobra.Command{}, "root")

		fakeLocator := locator.CreateFakeLocator()
		fakeLocator.CommandFoldersMock = func() []*models.CommandFolderInfo {
			return folders
		}
		reg.Set(locator.Identifier, fakeLocator)
		reg.Set(prompter.Identifier, prompter.CreateFakePrompter())
		reg.Set(shell.Identifier, shell.CreateFakeShell())

		newCmdSuccess := &dynamicCmd{
			cobraCmd: &cobra.Command{
				Use: "test-use",
			},
			addSelectSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd _select.NewCommandFunc) error {
				return nil
			},
			addStatusSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd status.NewCommandFunc) error {
				return nil
			},
			addRunSubCmdFunc: func(parent cmd.AnchorCommand, commandFolderName string, createCmd run.NewCommandFunc) error {
				return nil
			},
		}
		createCmdsFunc := func(ctx common.Context, commandFolders []*models.CommandFolderInfo) ([]*dynamicCmd, error) {
			assert.Equal(t, folders, commandFolders)
			return []*dynamicCmd{newCmdSuccess}, nil
		}
		err := AddCommands(parentCmd, createCmdsFunc)
		assert.Nil(t, err, "expected to succeed fail adding sub commands")

		assert.True(t, parentCmd.GetCobraCmd().HasSubCommands())
		parentCommands := parentCmd.GetCobraCmd().Commands()
		assert.Equal(t, 1, len(parentCommands))

		commandFolder1 := parentCommands[0]
		assert.Equal(t, "test-use", commandFolder1.Use)
	})
}
