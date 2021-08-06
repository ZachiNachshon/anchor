package root

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/config/resolver"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func StartRootCommandLoadRepoOrFailFlow(ctx common.Context) {
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

func scanAnchorfilesRepositoryTree(ctx common.Context, repoPath string) {
	l, _ := locator.FromRegistry(ctx.Registry())
	err := l.Scan(repoPath)
	if err != nil {
		logger.Fatalf("Failed to locate and scan anchorfiles repository content")
	}
}
