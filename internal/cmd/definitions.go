package cmd

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/prompter"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/spf13/cobra"
)

type AnchorCommand interface {
	GetCobraCmd() *cobra.Command
	GetContext() common.Context
}

type PreRunSequence func(ctx common.Context) error
type SetLoggerVerbosityFunc func(l logger.Logger, verbose bool) error

type New func(ctx common.Context)
type Init func(ctx common.Context)
type Append func(parent *cobra.Command)

type AnchorCollaborators struct {
	ResolveConfigContext func(ctx common.Context, prmpt prompter.Prompter, s shell.Shell) error
	LoadRepository       func(ctx common.Context) (string, error)
	ScanAnchorfiles      func(ctx common.Context, repoPath string) error
}

func (c *AnchorCollaborators) Run(ctx common.Context, p prompter.Prompter, s shell.Shell) error {
	err := c.ResolveConfigContext(ctx, p, s)
	if err != nil {
		return err
	}

	repoPath, err := c.LoadRepository(ctx)
	if err != nil {
		return err
	}

	err = c.ScanAnchorfiles(ctx, repoPath)
	if err != nil {
		return err
	}

	return nil
}
