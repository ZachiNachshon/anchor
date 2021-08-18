package config_cmd

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/edit"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/set_context_entry"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/use_context"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/config_cmd/view"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/spf13/cobra"
)

type configCmd struct {
	cmd.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context

	addEditSubCmdFunc func(
		parent cmd.AnchorCommand,
		cfgManager config.ConfigManager,
		createCmd edit.NewCommandFunc) error

	addSetContextEntrySubCmdFunc func(
		parent cmd.AnchorCommand,
		cfgManager config.ConfigManager,
		createCmd set_context_entry.NewCommandFunc) error

	addUseContextSubCmdFunc func(
		parent cmd.AnchorCommand,
		cfgManager config.ConfigManager,
		createCmd use_context.NewCommandFunc) error

	addViewSubCmdFunc func(
		parent cmd.AnchorCommand,
		cfgManager config.ConfigManager,
		createCmd view.NewCommandFunc) error
}

var validArgs = []string{"edit", "set-context-entry", "use-context", "view"}

type NewCommandFunc func(ctx common.Context) *configCmd

func NewCommand(ctx common.Context) *configCmd {
	var cobraCmd = &cobra.Command{
		Use:       "config",
		Short:     "Configuration file management",
		Long:      `Configuration file management`,
		ValidArgs: validArgs,
	}

	return &configCmd{
		cobraCmd:                     cobraCmd,
		ctx:                          ctx,
		addEditSubCmdFunc:            edit.AddCommand,
		addSetContextEntrySubCmdFunc: set_context_entry.AddCommand,
		addUseContextSubCmdFunc:      use_context.AddCommand,
		addViewSubCmdFunc:            view.AddCommand,
	}
}

func (c *configCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *configCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(parent cmd.AnchorCommand, createCmd NewCommandFunc) error {
	if cfgManager, err := parent.GetContext().Registry().SafeGet(config.Identifier); err != nil {
		return err
	} else {
		newCmd := createCmd(parent.GetContext())
		manager := cfgManager.(config.ConfigManager)

		err = newCmd.addEditSubCmdFunc(newCmd, manager, edit.NewCommand)
		if err != nil {
			return err
		}

		err = newCmd.addSetContextEntrySubCmdFunc(newCmd, manager, set_context_entry.NewCommand)
		if err != nil {
			return err
		}

		err = newCmd.addUseContextSubCmdFunc(newCmd, manager, use_context.NewCommand)
		if err != nil {
			return err
		}

		err = newCmd.addViewSubCmdFunc(newCmd, manager, view.NewCommand)
		if err != nil {
			return err
		}
		parent.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
		return nil
	}
}
