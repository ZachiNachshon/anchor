package docker

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/docker/build"
	"github.com/ZachiNachshon/anchor/cmd/anchor/docker/purge"
	"github.com/ZachiNachshon/anchor/cmd/anchor/docker/push"
	"github.com/ZachiNachshon/anchor/cmd/anchor/docker/remove"
	"github.com/ZachiNachshon/anchor/cmd/anchor/docker/run"
	"github.com/ZachiNachshon/anchor/cmd/anchor/docker/stop"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/installer"
	"github.com/spf13/cobra"
)

type dockerCmd struct {
	cobraCmd *cobra.Command
	opts     DockerOptions
}

type DockerOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

var validArgs = []string{"build", "purge", "push", "remove", "run", "stop"}

func NewCommand(opts *common.CmdRootOptions) *dockerCmd {
	var cobraCmd = &cobra.Command{
		Use:       "docker",
		Short:     "Docker commands",
		Aliases:   []string{"d"},
		ValidArgs: validArgs,
	}

	var dockerCmd = new(dockerCmd)
	dockerCmd.cobraCmd = cobraCmd
	dockerCmd.opts.CmdRootOptions = opts

	if err := checkPrerequisites(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := dockerCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := dockerCmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return dockerCmd
}

func (cmd *dockerCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func checkPrerequisites() error {
	if err := installer.NewDockerInstaller(common.ShellExec).Check(); err != nil {
		return err
	}

	return nil
}

func (cmd *dockerCmd) initFlags() error {
	return nil
}

func (cmd *dockerCmd) initSubCommands() error {

	// Docker Commands
	cmd.initDockerCommands()

	return nil
}

func (cmd *dockerCmd) initDockerCommands() {
	opts := cmd.opts.CmdRootOptions

	cmd.cobraCmd.AddCommand(build.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(remove.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(push.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(run.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(stop.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(purge.NewCommand(opts).GetCobraCmd())
}
