package destroy

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type destroyCmd struct {
	cobraCmd *cobra.Command
	opts     DestroyOptions
}

type DestroyOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *destroyCmd {
	var cobraCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy local Kubernetes cluster",
		Long:  `Destroy local Kubernetes cluster`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cluster.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.ClusterHeadline, "Destroy")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if err := cluster.Destroy(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var destroyCmd = new(destroyCmd)
	destroyCmd.cobraCmd = cobraCmd
	destroyCmd.opts.CmdRootOptions = opts

	if err := destroyCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return destroyCmd
}

func (cmd *destroyCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *destroyCmd) initFlags() error {
	return nil
}
