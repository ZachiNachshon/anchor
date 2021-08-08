package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/internal/repository"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

var LoadRepoOrFail = func(ctx common.Context) error {
	cfg := config.FromContext(ctx)

	err := loadConfigContextOrPrompt(ctx, &cfg)
	if err != nil {
		logger.Fatalf(err.Error())
		return nil
	}

	repo, err := repository.GetRepositoryBasedOnConfig(cfg.Config.ActiveContext.Context.Repository)
	if err != nil {
		logger.Fatalf(err.Error())
		return nil
	}

	repoPath, err := repo.Load(ctx)
	if err != nil {
		logger.Fatal(err.Error())
		return nil
	}

	ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(repoPath)
	scanAnchorfilesRepositoryTree(ctx, repoPath)
	return nil
}

var SetLoggerVerbosity = func(l logger.Logger, verbose bool) error {
	level := "info"
	if verbose {
		level = "debug"
	}
	if err := l.SetVerbosityLevel(level); err != nil {
		return err
	}
	return nil
}

func loadConfigContextOrPrompt(ctx common.Context, cfg *config.AnchorConfig) error {
	contextName := cfg.Config.CurrentContext
	if len(contextName) == 0 {
		var prompt prompter.Prompter
		if resolved, err := ctx.Registry().SafeGet(prompter.Identifier); err != nil {
			return err
		} else {
			prompt = resolved.(prompter.Prompter)
			if selectedCfgCtx, err := prompt.PromptConfigContext(cfg.Config.Contexts); err != nil {
				return err
			} else if selectedCfgCtx.Name == prompter.CancelActionName {
				return fmt.Errorf("cannot proceed without selecting a configuration context, aborting")
			} else {
				var s shell.Shell
				if resolved, err := ctx.Registry().SafeGet(shell.Identifier); err == nil {
					// Do not fail if screen cannot be cleared
					s = resolved.(shell.Shell)
					_ = s.ClearScreen()
				}
				contextName = selectedCfgCtx.Name
			}
		}
	}
	return config.LoadActiveConfigByName(cfg, contextName)
}

func scanAnchorfilesRepositoryTree(ctx common.Context, repoPath string) {
	var l locator.Locator
	if resolved, err := ctx.Registry().SafeGet(locator.Identifier); err != nil {
		logger.Fatal(err.Error())
	} else {
		l = resolved.(locator.Locator)
		err := l.Scan(repoPath)
		if err != nil {
			logger.Fatalf("Failed to locate and scan anchorfiles repository content")
		}
	}
}
