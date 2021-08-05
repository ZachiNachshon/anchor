package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/cmd/anchor/app"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cli"
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion"
	configCmd "github.com/ZachiNachshon/anchor/cmd/anchor/config"
	"github.com/ZachiNachshon/anchor/cmd/anchor/controller"
	"github.com/ZachiNachshon/anchor/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/resolver"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	version_pkg "github.com/ZachiNachshon/anchor/pkg/version"
	"github.com/spf13/cobra"
)

type anchorCmd struct {
	models.AnchorCommand
	cobraCmd *cobra.Command
	ctx      common.Context
}

var validArgs = []string{"app", "controller", "completion", "list", "version"}

func NewCommand(ctx common.Context) *anchorCmd {
	var rootCmd = &cobra.Command{
		Use:       "anchor",
		Short:     "Anchor your Ops environment into a version controlled repository",
		Long:      `Anchor your Ops environment into a version controlled repository`,
		ValidArgs: validArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := alignLoggerWithVerboseFlag(); err != nil {
				logger.Fatal(err.Error())
			}
		},
	}

	return &anchorCmd{
		cobraCmd: rootCmd,
		ctx:      ctx,
	}
}

func (cmd *anchorCmd) InitFlags() {
	cmd.cobraCmd.PersistentFlags().BoolVarP(
		&common.GlobalOptions.Verbose,
		"verbose",
		"v",
		common.GlobalOptions.Verbose,
		"anchor <command> -v")

	cmd.cobraCmd.PersistentFlags().SortFlags = false
}

func (cmd *anchorCmd) InitSubCommands() {

	//cobra.EnableCommandSorting = false

	// Apps Commands
	cmd.cobraCmd.AddCommand(app.NewCommand(cmd.ctx, loadRepoOrFail).GetCobraCmd())

	// CLI Commands
	cmd.cobraCmd.AddCommand(cli.NewCommand(cmd.ctx, loadRepoOrFail).GetCobraCmd())

	// Controller Commands
	cmd.cobraCmd.AddCommand(controller.NewCommand(cmd.ctx, loadRepoOrFail).GetCobraCmd())

	// Config Commands
	cmd.cobraCmd.AddCommand(configCmd.NewCommand(cmd.ctx).GetCobraCmd())

	// Version
	cmd.cobraCmd.AddCommand(version.NewCommand(cmd.ctx, version_pkg.DefineVersionActions()).GetCobraCmd())

	// Auto completion
	cmd.cobraCmd.AddCommand(completion.NewCommand(cmd.cobraCmd, cmd.ctx).GetCobraCmd())
}

func (cmd *anchorCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *anchorCmd) Execute() {
	// Cannot run on the command Run() method itself since we must initialize the logger
	// logger must be available at the PersistentPreRun() stage
	cmd.InitFlags()
	cmd.InitSubCommands()

	if err := cmd.cobraCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}

func loadConfigContextOrPrompt(ctx common.Context, cfg *config.AnchorConfig) error {
	contextName := cfg.Config.CurrentContext
	if len(contextName) == 0 {
		if prompt, err := prompter.FromRegistry(ctx.Registry()); err != nil {
			return err
		} else {
			if selectedCfgCtx, err := prompt.PromptConfigContext(cfg.Config.Contexts); err != nil {
				return err
			} else if selectedCfgCtx.Name == prompter.CancelActionName {
				return fmt.Errorf("cannot proceed without selecting a configuration context, aborting")
			} else {
				// Do not fail if screen cannot be cleared
				if s, err := shell.FromRegistry(ctx.Registry()); err == nil {
					_ = s.ClearScreen()
				}
				contextName = selectedCfgCtx.Name
			}
		}
	}
	return config.LoadActiveConfigByName(cfg, contextName)
}

func loadRepoOrFail(ctx common.Context) {
	cfg := ctx.Config().(config.AnchorConfig)

	err := loadConfigContextOrPrompt(ctx, &cfg)
	if err != nil {
		logger.Fatalf(err.Error())
		return
	}

	rslvr, err := resolver.GetResolverBasedOnConfig(cfg.Config.ActiveContext.Context.Repository)
	if err != nil {
		logger.Fatalf(err.Error())
		return
	}

	repoPath, err := rslvr.ResolveRepository(ctx)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(repoPath)
	scanAnchorfilesRepositoryTree(ctx, repoPath)
}

func scanAnchorfilesRepositoryTree(ctx common.Context, repoPath string) {
	l, _ := locator.FromRegistry(ctx.Registry())
	err := l.Scan(repoPath)
	if err != nil {
		logger.Fatalf("Failed to locate and scan anchorfiles repository content")
	}
}

func alignLoggerWithVerboseFlag() error {
	level := "info"
	if common.GlobalOptions.Verbose {
		level = "debug"
	}
	if err := logger.SetVerbosityLevel(level); err != nil {
		return err
	}
	return nil
}

func Main(ctx common.Context) {
	NewCommand(ctx).Execute()
}
