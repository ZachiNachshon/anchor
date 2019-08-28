package cluster

import (
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/apply"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/backup"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/connect"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/create"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/dashboard"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/delete"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/deploy"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/destroy"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/expose"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/log"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/registry"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/status"
	"github.com/ZachiNachshon/anchor/cmd/anchor/cluster/token"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/installer"
	"github.com/spf13/cobra"
)

type clusterCmd struct {
	cobraCmd *cobra.Command
	opts     ClusterOptions
}

type ClusterOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

var validArgs = []string{"apply", "connect", "create", "dashboard", "delete", "deploy", "expose", "registry", "status", "token"}

func NewCommand(opts *common.CmdRootOptions) *clusterCmd {
	var cobraCmd = &cobra.Command{
		Use:       "cluster",
		Short:     "Cluster commands",
		Aliases:   []string{"c"},
		ValidArgs: validArgs,
	}

	var clusterCmd = new(clusterCmd)
	clusterCmd.cobraCmd = cobraCmd
	clusterCmd.opts.CmdRootOptions = opts

	if err := checkPrerequisites(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := clusterCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := clusterCmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return clusterCmd
}

func (cmd *clusterCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func checkPrerequisites() error {
	if err := installer.NewKindInstaller(common.ShellExec).Check(); err != nil {
		return err
	}
	if err := installer.NewEnvsubstInstaller(common.ShellExec).Check(); err != nil {
		return err
	}
	if err := installer.NewKubectlInstaller(common.ShellExec).Check(); err != nil {
		return err
	}
	if err := installer.NewHostessInstaller(common.ShellExec).Check(); err != nil {
		return err
	}
	//if err := shell.NewHelmlInstaller(common.ShellExec).Check(); err != nil {
	//	return err
	//}
	return nil
}

func (cmd *clusterCmd) initFlags() error {
	return nil
}

func (cmd *clusterCmd) initSubCommands() error {

	// Kind Commands
	cmd.initClusterCommands()

	return nil
}

func (cmd *clusterCmd) initClusterCommands() {
	opts := cmd.opts.CmdRootOptions

	cmd.cobraCmd.AddCommand(create.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(apply.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(dashboard.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(destroy.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(status.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(deploy.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(delete.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(backup.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(registry.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(expose.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(connect.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(token.NewCommand(opts).GetCobraCmd())
	cmd.cobraCmd.AddCommand(log.NewCommand(opts).GetCobraCmd())
}
