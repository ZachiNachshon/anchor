package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/repository"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

type AnchorCollaborators struct {
	resolveConfigContext func(ctx common.Context, prmpt prompter.Prompter, s shell.Shell) error
	loadRepository       func(ctx common.Context) (string, error)
	scanAnchorfiles      func(ctx common.Context, repoPath string) error
}

var AnchorPreRunSequence = func() *AnchorCollaborators {
	return &AnchorCollaborators{
		resolveConfigContext: loadConfigContext,
		loadRepository:       loadRepository,
		scanAnchorfiles:      scanAnchorfilesRepositoryTree,
	}
}

func (c *AnchorCollaborators) Run(ctx common.Context) error {
	p, err := ctx.Registry().SafeGet(prompter.Identifier)
	if err != nil {
		return err
	}

	s, err := ctx.Registry().SafeGet(shell.Identifier)
	if err != nil {
		return err
	}

	err = c.resolveConfigContext(ctx, p.(prompter.Prompter), s.(shell.Shell))
	if err != nil {
		return err
	}

	repoPath, err := c.loadRepository(ctx)
	if err != nil {
		return err
	}

	err = c.scanAnchorfiles(ctx, repoPath)
	if err != nil {
		return err
	}

	return nil
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
	if repo, err := repository.GetRepositoryOriginByConfig(cfg.Config.ActiveContext.Context.Repository); err != nil {
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
		err := l.Scan(repoPath)
		if err != nil {
			errMsg := fmt.Sprintf("failed to scan anchorfiles repository. error: %s", err.Error())
			return fmt.Errorf(errMsg)
		}
	}
	return nil
}
