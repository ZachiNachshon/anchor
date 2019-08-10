package cmd

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
	"os"
)

type completionCmd struct {
	cobraCmd *cobra.Command
	opts     VersionCmdOptions
}

type CompletionCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCompletionCmd(root *RootCmd, opts *common.CmdRootOptions) *completionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate auto completion script for bash/zsh",
		Long:  `Generate auto completion script for bash/zsh`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if common.GlobalOptions.AutoCompleteGenerateFile {
				logger.PrintHeadline(fmt.Sprintf("Generating Anchor Auto Completion: %v", common.GlobalOptions.AutoCompletionDefaultShell))
				generateAutoCompletionFile(root)
				logger.PrintCompletion()
			} else {
				generateAutoCompletionStdOut(root)
			}
		},
	}

	var completionCmd = new(completionCmd)
	completionCmd.cobraCmd = cobraCmd
	completionCmd.opts.CmdRootOptions = opts

	if err := completionCmd.initFlags(); err != nil {
		// TODO: log error
	}

	return completionCmd
}

func (cmd *completionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *completionCmd) initFlags() error {
	cmd.cobraCmd.Flags().StringVarP(
		&common.GlobalOptions.AutoCompletionDefaultShell,
		"shell",
		"s",
		common.GlobalOptions.AutoCompletionDefaultShell,
		"anchor completion -s bash/zsh")

	_ = cmd.cobraCmd.MarkFlagRequired("shell")

	cmd.cobraCmd.Flags().BoolVarP(
		&common.GlobalOptions.AutoCompleteGenerateFile,
		"file",
		"f",
		common.GlobalOptions.AutoCompleteGenerateFile,
		"anchor completion -f")

	//_ = cmd.cobraCmd.MarkFlagRequired("file")
	return nil
}

func generateAutoCompletionFile(root *RootCmd) {
	var filename = ""
	switch common.GlobalOptions.AutoCompletionDefaultShell {
	case "bash":
		{
			filename = common.GlobalOptions.AutoCompletionDefaultFilePrefix + ".sh"
			_ = root.cobraCmd.GenZshCompletionFile(filename)
			break
		}
	case "zsh":
		{
			filename = common.GlobalOptions.AutoCompletionDefaultFilePrefix + ".zsh"
			_ = root.cobraCmd.GenZshCompletionFile(filename)
			break
		}
	default:
		{
			logger.Fatalf("%v is not supported for auto completion", common.GlobalOptions.AutoCompletionDefaultShell)
		}
	}

	if dir, err := os.Getwd(); err == nil {
		compInfo := fmt.Sprintf(`
  Successfully created auto completion file for [%v].
  Please add %v/%v to your $PATH.`, common.GlobalOptions.AutoCompletionDefaultShell, dir, filename)
		logger.Info(compInfo)
	}
}

func generateAutoCompletionStdOut(root *RootCmd) {
	compInfo := fmt.Sprintf(`
#
# Add on of the following command to your ~/.bash_profile or ~/.bashrc:
#   - source <(anchor completion -s bash)
#   - source <(anchor completion -s zsh)
#
# Auto completion script is as follows:
#
`)
	logger.Info(compInfo)

	switch common.GlobalOptions.AutoCompletionDefaultShell {
	case "bash":
		{
			_ = root.cobraCmd.GenZshCompletion(os.Stdout)
			break
		}
	case "zsh":
		{
			_ = root.cobraCmd.GenZshCompletion(os.Stdout)
			break
		}
	default:
		{
			logger.Fatalf("%v is not supported for auto completion", common.GlobalOptions.AutoCompletionDefaultShell)
		}
	}
}
