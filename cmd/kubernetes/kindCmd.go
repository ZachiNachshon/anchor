package kubernetes

import (
	"fmt"
	"github.com/kit/cmd/docker"
	"github.com/kit/cmd/types"

	"github.com/spf13/cobra"
)

type KindCmd struct {
	cobraCmd *cobra.Command
	opts     KindCmdOptions
}

type KindCmdOptions struct {
	*types.CmdRootOptions

	// Additional Build Params
}

func NewKindCmd(opts *types.CmdRootOptions) *KindCmd {
	var cobraCmd = &cobra.Command{
		Use:   "build",
		Short: "Builds a Dockerfile",
		Long:  `Builds a docker image from the DOCKER_FILES repository.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Build a Dockerfile...")
		},
	}

	var kindCmd = new(KindCmd)
	kindCmd.cobraCmd = cobraCmd
	kindCmd.opts.CmdRootOptions = opts

	if err := kindCmd.initFlags(); err != nil {
		// TODO: log error
	}

	return kindCmd
}

func (cmd *KindCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *KindCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}

func (root *CmdRoot) initSubCommands() error {

	// Docker Commands
	root.initDockerCommands()

	// Kubernetes Commands
	root.initKubernetesCommands()

	// Admin
	root.cobraCmd.AddCommand(NewVersionCmd(&root.opts).GetCobraCmd())

	return nil
}

func (root *CmdRoot) initDockerCommands() {
	root.cobraCmd.AddCommand(docker.NewBuildCmd(&root.opts).GetCobraCmd())
	root.cobraCmd.AddCommand(docker.NewCleanCmd(&root.opts).GetCobraCmd())
	root.cobraCmd.AddCommand(docker.NewListCmd(&root.opts).GetCobraCmd())
	root.cobraCmd.AddCommand(docker.NewPushCmd(&root.opts).GetCobraCmd())
	root.cobraCmd.AddCommand(docker.NewRunCmd(&root.opts).GetCobraCmd())
	root.cobraCmd.AddCommand(docker.NewStopCmd(&root.opts).GetCobraCmd())
}
