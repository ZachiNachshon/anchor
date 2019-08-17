package node

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type nodeCmd struct {
	cobraCmd *cobra.Command
}

func NewCommand() *nodeCmd {
	var cobraCmd = &cobra.Command{
		Use:   "node",
		Short: "Connect to a kubernetes node (if ^M appear as Enter, run - stty sane)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.ClusterHeadline, "Connect Node")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if err := cluster.ConnectToNode(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var nodeCmd = new(nodeCmd)
	nodeCmd.cobraCmd = cobraCmd

	return nodeCmd
}

func (cmd *nodeCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}
