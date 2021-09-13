package anchor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/repository"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/locator"
	"github.com/ZachiNachshon/anchor/pkg/parser"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func NewAnchorCollaborators() *AnchorCollaborators {
	return &AnchorCollaborators{
		prepareRegistryItemsFunc: prepareRegistryItems,
		resolveConfigContextFunc: resolveConfigContext,
		loadRepositoryFunc:       loadRepository,
		scanAnchorfilesFunc:      scanAnchorfilesRepositoryTree,
	}
}

type AnchorCollaborators struct {
	prmptr prompter.Prompter
	s      shell.Shell
	l      locator.Locator
	e      extractor.Extractor
	prsr   parser.Parser

	prepareRegistryItemsFunc func(c *AnchorCollaborators, ctx common.Context) error
	resolveConfigContextFunc func(c *AnchorCollaborators, ctx common.Context) error
	loadRepositoryFunc       func(c *AnchorCollaborators, ctx common.Context) (string, error)
	scanAnchorfilesFunc      func(c *AnchorCollaborators, ctx common.Context, repoPath string) error
}

func (c *AnchorCollaborators) Run(ctx common.Context) error {
	err := c.prepareRegistryItemsFunc(c, ctx)
	if err != nil {
		return err
	}

	err = c.resolveConfigContextFunc(c, ctx)
	if err != nil {
		return err
	}

	repoPath, err := c.loadRepositoryFunc(c, ctx)
	if err != nil {
		return err
	}

	err = c.scanAnchorfilesFunc(c, ctx, repoPath)
	if err != nil {
		return err
	}

	return nil
}

func prepareRegistryItems(c *AnchorCollaborators, ctx common.Context) error {
	reg := ctx.Registry()
	if s, err := reg.SafeGet(shell.Identifier); err != nil {
		return err
	} else {
		c.s = s.(shell.Shell)
	}
	if prmptr, err := reg.SafeGet(prompter.Identifier); err != nil {
		return err
	} else {
		c.prmptr = prmptr.(prompter.Prompter)
	}
	if l, err := reg.SafeGet(locator.Identifier); err != nil {
		return err
	} else {
		c.l = l.(locator.Locator)
	}
	if e, err := reg.SafeGet(extractor.Identifier); err != nil {
		return err
	} else {
		c.e = e.(extractor.Extractor)
	}
	if prsr, err := reg.SafeGet(parser.Identifier); err != nil {
		return err
	} else {
		c.prsr = prsr.(parser.Parser)
	}
	return nil
}

func resolveConfigContext(c *AnchorCollaborators, ctx common.Context) error {
	resolved, err := ctx.Registry().SafeGet(config.Identifier)
	if err != nil {
		return err
	}
	cfgManager := resolved.(config.ConfigManager)

	cfg := config.FromContext(ctx)
	contextName := cfg.Config.CurrentContext
	if len(contextName) == 0 {
		if selectedCfgCtx, err := c.prmptr.PromptConfigContext(cfg.Config.Contexts); err != nil {
			return err
		} else if selectedCfgCtx.Name == prompter.CancelActionName {
			return fmt.Errorf("cannot proceed without selecting a configuration context, aborting")
		} else {
			_ = c.s.ClearScreen()
			contextName = selectedCfgCtx.Name
		}
	}
	return cfgManager.SwitchActiveConfigContextByName(cfg, contextName)
}

func loadRepository(c *AnchorCollaborators, ctx common.Context) (string, error) {
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

func scanAnchorfilesRepositoryTree(c *AnchorCollaborators, ctx common.Context, repoPath string) error {
	locatorErr := c.l.Scan(repoPath, c.e, c.prsr)
	if locatorErr != nil {
		errMsg := fmt.Sprintf("failed to scan anchorfiles repository. error: %s", locatorErr.GoError().Error())
		return fmt.Errorf(errMsg)
	}
	return nil
}
