package delete

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var namespace = common.GlobalOptions.DockerImageNamespace

type deleteCmd struct {
	cobraCmd *cobra.Command
	opts     DeleteOptions
}

type DeleteOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *deleteCmd {
	var cobraCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a previously deployed container manifest",
		Long:  `Delete a previously deployed container manifest`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline(logger.ClusterHeadline, "Delete")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if _, err := cluster.Delete(args[0], namespace); err != nil {
				logger.Fatal(err.Error())
			}
			if err := cluster.DisablePortForwarding(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var removeCmd = new(deleteCmd)
	removeCmd.cobraCmd = cobraCmd
	removeCmd.opts.CmdRootOptions = opts

	if err := removeCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return removeCmd
}

func (cmd *deleteCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *deleteCmd) initFlags() error {
	cmd.cobraCmd.PersistentFlags().StringVarP(
		&namespace,
		"namespace",
		"n",
		namespace,
		"anchor cluster apply <name> -n <namespace>")
	return nil
}
