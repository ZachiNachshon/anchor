package completion

import (
	"github.com/ZachiNachshon/anchor/internal/cmd"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion/bash"
	"github.com/ZachiNachshon/anchor/internal/cmd/anchor/completion/zsh"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/spf13/cobra"
)

type completionCmd struct {
	cobraCmd *cobra.Command
	ctx      common.Context

	addBashSubCmdFunc func(root cmd.AnchorCommand, parent cmd.AnchorCommand, createCmd bash.NewCommandFunc) error
	addZshSubCmdFunc  func(root cmd.AnchorCommand, parent cmd.AnchorCommand, createCmd zsh.NewCommandFunc) error
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

type NewCommandFunc func(ctx common.Context) *completionCmd

func NewCommand(ctx common.Context) *completionCmd {
	var cobraCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate auto completion script for [bash, zsh]",
		Long:  longDescription,
		Args:  cobra.NoArgs,
	}

	return &completionCmd{
		cobraCmd:          cobraCmd,
		ctx:               ctx,
		addBashSubCmdFunc: bash.AddCommand,
		addZshSubCmdFunc:  zsh.AddCommand,
	}
}

func (c *completionCmd) GetCobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *completionCmd) GetContext() common.Context {
	return c.ctx
}

func AddCommand(root cmd.AnchorCommand, createCmd NewCommandFunc) error {
	newCmd := createCmd(root.GetContext())

	err := newCmd.addBashSubCmdFunc(root, newCmd, bash.NewCommand)
	if err != nil {
		return err
	}

	err = newCmd.addZshSubCmdFunc(root, newCmd, zsh.NewCommand)
	if err != nil {
		return err
	}

	root.GetCobraCmd().AddCommand(newCmd.GetCobraCmd())
	return nil
}
