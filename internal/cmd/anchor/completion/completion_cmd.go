package completion

import (
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion/bash"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion/zsh"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"

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

func NewCommand(root *cobra.Command, ctx common.Context) (*completionCmd, error) {
	var cobraCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate auto completion script for [bash, zsh]",
		Long:  longDescription,
		Args:  cobra.NoArgs,
	}

	var compCmd = &completionCmd{
		cobraCmd: cobraCmd,
		ctx:      ctx,
	}

	if err := compCmd.InitSubCommands(root); err != nil {
		logger.Fatal(err.Error())
	}

	return compCmd, nil
}

func (cmd *completionCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *completionCmd) InitFlags() error {
	return nil
}

func (cmd *completionCmd) InitSubCommands(root *cobra.Command) error {
	if bashCmd, err := bash.NewCommand(root); err != nil {
		return err
	} else {
		cmd.cobraCmd.AddCommand(bashCmd.GetCobraCmd())
	}

	if zshCmd, err := zsh.NewCommand(root); err != nil {
		return err
	} else {
		cmd.cobraCmd.AddCommand(zshCmd.GetCobraCmd())
	}
	return nil
}
