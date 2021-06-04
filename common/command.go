package common

import "github.com/spf13/cobra"

type CliCommand interface {
	GetCobraCmd() *cobra.Command
	InitFlags()
	InitSubCommands()
}
