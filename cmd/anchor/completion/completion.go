package completion

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion/bash"
	"github.com/ZachiNachshon/anchor/cmd/anchor/completion/zsh"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/spf13/cobra"
)

type completionCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context
}

const longDescription = `
Outputs anchor shell completion for the given shell (bash or zsh)
This depends on the bash-completion binary. Example installation instructions:

# for bash users
	$ anchor completion bash > ~/.anchor-completion
	$ source ~/.anchor-completion

# for zsh users
	% anchor completion zsh > /usr/local/share/zsh/site-functions/_anchor
	% autoload -U compinit && compinit

Additionally, you may want to output the completion to a file and source in your .bashrc
Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2
`

func NewCommand(root *cobra.Command, ctx common.Context) *completionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate auto completion script for [bash, zsh]",
		Long:  longDescription,
		Args:  cobra.NoArgs,
	}

	var completionCmd = &completionCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	if err := completionCmd.initSubCommands(root); err != nil {
		logger.Fatal(err.Error())
	}

	return completionCmd
}

func (cmd *completionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *completionCmd) initSubCommands(root *cobra.Command) error {
	cmd.cobraCmd.AddCommand(bash.NewCommand(root).GetCobraCmd())
	cmd.cobraCmd.AddCommand(zsh.NewCommand(root).GetCobraCmd())
	return nil
}
