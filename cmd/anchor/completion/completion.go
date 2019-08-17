package completion

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion/bash"
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion/zsh"
	"github.com/ZachiNachshon/anchor/cmd/anchor/version"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type CompletionCmd struct {
	cobraCmd *cobra.Command
	opts     version.VersionOptions
}

type CompletionOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

const longDescription = `
Outputs anchor shell completion for the given shell (bash or zsh)
This depends on the bash-completion binary.  Example installation instructions:
# for bash users
	$ anchor completion bash > ~/.anchor-completion
	$ source ~/.anchor-completion

# for zsh users
	% anchor completion zsh > /usr/local/share/zsh/site-functions/_anchor
	% autoload -U compinit && compinit

Additionally, you may want to output the completion to a file and source in your .bashrc
Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2
`

func NewCommand(root *cobra.Command, opts *common.CmdRootOptions) *CompletionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate auto completion script for bash/zsh",
		Long:  longDescription,
		Args:  cobra.NoArgs,
	}

	var completionCmd = new(CompletionCmd)
	completionCmd.cobraCmd = cobraCmd
	completionCmd.opts.CmdRootOptions = opts

	if err := completionCmd.initFlags(); err != nil {
		// TODO: log error
	}

	if err := completionCmd.initSubCommands(root); err != nil {
		logger.Fatal(err.Error())
	}

	return completionCmd
}

func (cmd *CompletionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *CompletionCmd) initFlags() error {
	return nil
}

func (cmd *CompletionCmd) initSubCommands(root *cobra.Command) error {
	cmd.cobraCmd.AddCommand(bash.NewCommand(root).GetCobraCmd())
	cmd.cobraCmd.AddCommand(zsh.NewCommand(root).GetCobraCmd())
	return nil
}
