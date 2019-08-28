package registry

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var shouldDeleteRegistry = false

type registryCmd struct {
	cobraCmd *cobra.Command
	opts     RegistryOptions
}

type RegistryOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCommand(opts *common.CmdRootOptions) *registryCmd {
	var cobraCmd = &cobra.Command{
		Use:   "registry",
		Short: fmt.Sprintf("Create a private docker registry [%v]", common.GlobalOptions.DockerRegistryDns),
		Long:  fmt.Sprintf("Create a private docker registry [%v]", common.GlobalOptions.DockerRegistryDns),
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.ClusterHeadline, "Registry")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if shouldDeleteRegistry {
				// Remove registry
				if err := cluster.DeleteRegistry(); err != nil {
					logger.Fatal(err.Error())
				}
			} else {
				// Deploy registry
				if err := cluster.Registry(true); err != nil {
					logger.Fatal(err.Error())
				}
			}
			logger.PrintCompletion()
		},
	}

	var registryCmd = new(registryCmd)
	registryCmd.cobraCmd = cobraCmd
	registryCmd.opts.CmdRootOptions = opts

	if err := registryCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return registryCmd
}

func (cmd *registryCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *registryCmd) initFlags() error {
	cmd.cobraCmd.Flags().BoolVarP(
		&shouldDeleteRegistry,
		"delete",
		"d",
		shouldDeleteRegistry,
		"anchor cluster registry -d")
	return nil
}
