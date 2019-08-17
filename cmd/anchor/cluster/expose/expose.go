package expose

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type exposeCmd struct {
	cobraCmd *cobra.Command
	opts     ExposeOptions
}

type ExposeOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *exposeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "expose",
		Short: "Expose a container port to the host instance",
		Long:  `Expose a container port to the host instance`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.ClusterHeadline, "Expose")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if err := cluster.Expose(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var exposeCmd = new(exposeCmd)
	exposeCmd.cobraCmd = cobraCmd
	exposeCmd.opts.CmdRootOptions = opts

	if err := exposeCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}
	return exposeCmd
}

func (cmd *exposeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *exposeCmd) initFlags() error {
	return nil
}
