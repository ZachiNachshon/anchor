package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/installer"
	"github.com/spf13/cobra"
	"strings"
)

type dockerCmd struct {
	cobraCmd *cobra.Command
	opts     DockerCmdOptions
}

type DockerCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

var validArgs = []string{"build", "clean", "list", "purge", "push", "run", "stop"}

func NewDockerCmd(opts *common.CmdRootOptions) *dockerCmd {
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

func (d *dockerCmd) initSubCommands() error {

	// Docker Commands
	d.initDockerCommands()

	return nil
}

func (d *dockerCmd) initDockerCommands() {
	opts := d.opts.CmdRootOptions

	d.cobraCmd.AddCommand(NewBuildCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewRemoveCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewPushCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewRunCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewStopCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewPurgeCmd(opts).GetCobraCmd())
}

func RemoveNamespaceFromImageName(name string) string {
	if len(name) == 0 {
		return name
	}

	var result = name
	backslashIdx := strings.LastIndex(name, "/")
	if backslashIdx > 0 {
		result = result[backslashIdx+1:]
	}
	tagIdx := strings.LastIndex(name, ":")
	if tagIdx > 0 {
		result = result[:backslashIdx+1]
	}

	return result
}

func ComposeDockerContainerIdentifier(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v:%v", common.GlobalOptions.DockerImageNamespace, dirname, common.GlobalOptions.DockerImageTag)
	return imageIdentifier
}

func ComposeDockerContainerIdentifierNoTag(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v-%v", common.GlobalOptions.DockerImageNamespace, dirname)
	return imageIdentifier
}

func ComposeDockerImageIdentifierNoTag(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v", common.GlobalOptions.DockerImageNamespace, dirname)
	return imageIdentifier
}
