package dashboard

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var shouldDeleteDashboard = false

type dashboardCmd struct {
	cobraCmd *cobra.Command
	opts     DashboardOptions
}

type DashboardOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *dashboardCmd {
	var cobraCmd = &cobra.Command{
		Use:   "dashboard",
		Short: "Deploy a Kubernetes dashboard",
		Long:  `Deploy a Kubernetes dashboard`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.ClusterHeadline, "Dashboard")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if shouldDeleteDashboard {
				if err := cluster.UninstallDashboard(); err != nil {
					logger.Fatal(err.Error())
				}
			} else {
				if err := cluster.Dashboard(); err != nil {
					logger.Fatal(err.Error())
				}
			}
			logger.PrintCompletion()
		},
	}

	var dashboardCmd = new(dashboardCmd)
	dashboardCmd.cobraCmd = cobraCmd
	dashboardCmd.opts.CmdRootOptions = opts

	if err := dashboardCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return dashboardCmd
}

func (cmd *dashboardCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *dashboardCmd) initFlags() error {
	// TODO: Allow force creation by flag even if dashboard exists ?
	cmd.cobraCmd.Flags().BoolVarP(
		&shouldDeleteDashboard,
		"Delete Kubernetes dashboard",
		"d",
		shouldDeleteDashboard,
		"anchor cluster dashboard -d")
	return nil
}
