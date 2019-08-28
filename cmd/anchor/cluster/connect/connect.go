package connect

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/connect/node"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/connect/pod"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

type connectCmd struct {
	cobraCmd *cobra.Command
	opts     ConnectOptions
}

type ConnectOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewCommand(opts *common.CmdRootOptions) *connectCmd {
	var cobraCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a kubernetes [node, pod] by name",
		Long:  `Connect to a kubernetes [node, pod] by name`,
		Args:  cobra.NoArgs,
	}

	var connectCmd = new(connectCmd)
	connectCmd.cobraCmd = cobraCmd
	connectCmd.opts.CmdRootOptions = opts

	if err := connectCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := connectCmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return connectCmd
}

func (cmd *connectCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *connectCmd) initFlags() error {
	return nil
}

func (cmd *connectCmd) initSubCommands() error {
	cmd.cobraCmd.AddCommand(node.NewCommand().GetCobraCmd())
	cmd.cobraCmd.AddCommand(pod.NewCommand().GetCobraCmd())
	return nil
}
