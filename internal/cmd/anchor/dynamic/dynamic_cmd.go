package dynamic

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/run"
	_select "github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/select"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/dynamic/status"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/spf13/cobra"
)

type dynamicCmd struct {
	cmd.AnchorCommand
	cobraCmd    *cobra.Command
	ctx         common.Context
	commandName string

	addRunSubCmdFunc    func(parent cmd.AnchorCommand, parentFolderName string, createCmd run.NewCommandFunc) error
	addSelectSubCmdFunc func(parent cmd.AnchorCommand, parentFolderName string, createCmd _select.NewCommandFunc) error
	addStatusSubCmdFunc func(parent cmd.AnchorCommand, parentFolderName string, createCmd status.NewCommandFunc) error
}

type NewCommandsFunc func(ctx common.Context, anchorFolders []*models.AnchorFolderInfo) ([]*dynamicCmd, error)

func NewCommands(ctx common.Context, anchorFolders []*models.AnchorFolderInfo) ([]*dynamicCmd, error) {
	var cmds = make([]*dynamicCmd, len(anchorFolders))
	for i := 0; i < len(anchorFolders); i++ {
		aFolder := anchorFolders[i]
		cobraCmd := newCobraCommand(aFolder)
		newDynamicCmd := newCommand(ctx, cobraCmd, aFolder.Name)
		cmds[i] = newDynamicCmd
	}

	return cmds, nil
}

func newCommand(ctx common.Context, cobraCmd *cobra.Command, name string) *dynamicCmd {
	return &dynamicCmd{
		ctx:                 ctx,
		cobraCmd:            cobraCmd,
		commandName:         name,
		addRunSubCmdFunc:    run.AddCommand,
		addSelectSubCmdFunc: _select.AddCommand,
		addStatusSubCmdFunc: status.AddCommand,
	}
}

func newCobraCommand(anchorFolder *models.AnchorFolderInfo) *cobra.Command {
	return &cobra.Command{
		Use:   anchorFolder.Command.Use,
		Short: anchorFolder.Command.Short,
	}
}

func (c *dynamicCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *dynamicCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommands(parent cmd.AnchorCommand, anchorCollaborators *cmd.AnchorCollaborators, createCmds NewCommandsFunc) error {
	l, pr, s, err := resolveFromRegistry(parent.GetContext())
	if err != nil {
		return err
	}

	err = anchorCollaborators.Run(parent.GetContext(), pr, s)
	if err != nil {
		return err
	}

	anchorFolders := l.AnchorFolders()
	newCmds, err := createCmds(parent.GetContext(), anchorFolders)
	if err != nil {
		return err
	}

	for i := 0; i < len(newCmds); i++ {
		newCmd := newCmds[i]
		err = newCmd.addSelectSubCmdFunc(newCmd, newCmd.commandName, _select.NewCommand)
		if err != nil {
			return err
		}

		err = newCmd.addStatusSubCmdFunc(newCmd, newCmd.commandName, status.NewCommand)
		if err != nil {
			return err
		}

		err = newCmd.addRunSubCmdFunc(newCmd, newCmd.commandName, run.NewCommand)
		if err != nil {
			return err
		}

		parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	}

	return nil
}

func resolveFromRegistry(ctx common.Context) (locator.Locator, prompter.Prompter, shell.Shell, error) {
	l, err := ctx.Registry().SafeGet(locator.Identifier)
	if err != nil {
		return nil, nil, nil, err
	}

	pr, err := ctx.Registry().SafeGet(prompter.Identifier)
	if err != nil {
		return nil, nil, nil, err
	}

	s, err := ctx.Registry().SafeGet(shell.Identifier)
	if err != nil {
		return nil, nil, nil, err
	}

	return l.(locator.Locator), pr.(prompter.Prompter), s.(shell.Shell), nil
}
