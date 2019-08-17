package status

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	cobraCmd *cobra.Command
	opts     StatusOptions
}

type StatusOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *statusCmd {
	var cobraCmd = &cobra.Command{
		Use:   "status",
		Short: fmt.Sprintf("Print cluster [%v] status", common.GlobalOptions.KindClusterName),
		Long:  fmt.Sprintf(`Print cluster [%v] status`, common.GlobalOptions.KindClusterName),
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.ClusterHeadline, "Status")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if err := cluster.Status(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var statusCmd = new(statusCmd)
	statusCmd.cobraCmd = cobraCmd
	statusCmd.opts.CmdRootOptions = opts

	if err := statusCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return statusCmd
}

func (cmd *statusCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *statusCmd) initFlags() error {
	return nil
}
