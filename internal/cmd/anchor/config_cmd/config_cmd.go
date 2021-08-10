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
}

var validArgs = []string{"edit", "set-context-entry", "use-context", "view"}

func NewCommand(ctx common.Context) (*configCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:       "config",
		Short:     "Configuration file management",
		Long:      `Configuration file management`,
		ValidArgs: validArgs,
	}

	var cmd = &configCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	err := cmd.InitSubCommands()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (cmd *configCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *configCmd) InitFlags() error {
	return nil
}

func (cmd *configCmd) InitSubCommands() error {
	if cfgManager, err := cmd.ctx.Registry().SafeGet(config.Identifier); err != nil {
		return err
	} else {
		manager := cfgManager.(config.ConfigManager)
		if viewCmd, err := view.NewCommand(cmd.ctx, manager, view.ConfigView); err != nil {
			return err
		} else {
			cmd.cobraCmd.AddCommand(viewCmd.GetCobraCmd())
		}

		if editCmd, err := edit.NewCommand(cmd.ctx, manager, edit.ConfigEdit); err != nil {
			return err
		} else {
			cmd.cobraCmd.AddCommand(editCmd.GetCobraCmd())
		}

		if useCtxCmd, err := use_context.NewCommand(cmd.ctx, manager, use_context.ConfigUseContext); err != nil {
			return err
		} else {
			cmd.cobraCmd.AddCommand(useCtxCmd.GetCobraCmd())
		}

		if setCtxEntryCmd, err := set_context_entry.NewCommand(cmd.ctx, manager, set_context_entry.ConfigSetContextEntry); err != nil {
			return err
		} else {
			cmd.cobraCmd.AddCommand(setCtxEntryCmd.GetCobraCmd())
		}
		return nil
	}
}
