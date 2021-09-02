package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/repository"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

var AnchorPreRunSequence = func() *cmd.AnchorCollaborators {
	return &cmd.AnchorCollaborators{
		ResolveConfigContext: loadConfigContext,
		LoadRepository:       loadRepository,
		ScanAnchorfiles:      scanAnchorfilesRepositoryTree,
	}
}

func loadConfigContext(
	ctx common.Context,
	prmpt prompter.Prompter,
	s shell.Shell) error {

	resolved, err := ctx.Registry().SafeGet(config.Identifier)
	if err != nil {
		return err
	}
	cfgManager := resolved.(config.ConfigManager)

	cfg := config.FromContext(ctx)
	contextName := cfg.Config.CurrentContext
	if len(contextName) == 0 {
		if selectedCfgCtx, err := prmpt.PromptConfigContext(cfg.Config.Contexts); err != nil {
			return err
		} else if selectedCfgCtx.Name == prompter.CancelActionName {
			return fmt.Errorf("cannot proceed without selecting a configuration context, aborting")
		} else {
			_ = s.ClearScreen()
			contextName = selectedCfgCtx.Name
		}
	}
	return cfgManager.SwitchActiveConfigContextByName(cfg, contextName)
}

func loadRepository(ctx common.Context) (string, error) {
	cfg := config.FromContext(ctx)
	if repo, err := repository.GetRepositoryOriginByConfig(ctx, cfg.Config.ActiveContext.Context.Repository); err != nil {
		return "", err
	} else {
		if repoPath, err := repo.Load(ctx); err != nil {
			return "", err
		} else {
			ctx.(common.AnchorFilesPathSetter).SetAnchorFilesPath(repoPath)
			return repoPath, nil
		}
	}
}

func scanAnchorfilesRepositoryTree(ctx common.Context, repoPath string) error {
	var l locator.Locator
	if resolved, err := ctx.Registry().SafeGet(locator.Identifier); err != nil {
		return err
	} else {
		l = resolved.(locator.Locator)
		locatorErr := l.Scan(repoPath)
		if locatorErr != nil {
			errMsg := fmt.Sprintf("failed to scan anchorfiles repository. error: %s", locatorErr.GoError().Error())
			return fmt.Errorf(errMsg)
		}
	}
	return nil
}
