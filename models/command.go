package models

import "github.com/spf13/cobra"

type AnchorCommand interface {
	GetCobraCmd() *cobra.Command
	InitFlags()
	InitSubCommands()
}
