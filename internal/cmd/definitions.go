package cmd

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/spf13/cobra"
)

type AnchorCommand interface {
	GetCobraCmd() *cobra.Command
	InitFlags()
	InitSubCommands()
}

type LoadRepoOrFailFunc func(ctx common.Context) error
type SetLoggerVerbosityFunc func(l logger.Logger, verbose bool) error
