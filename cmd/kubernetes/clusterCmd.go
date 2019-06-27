package kubernetes

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/shell"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type clusterCmd struct {
	cobraCmd *cobra.Command
	opts     ClusterCmdOptions
}

type ClusterCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewKindCmd(opts *common.CmdRootOptions) *clusterCmd {
	var cobraCmd = &cobra.Command{
		Use:   "cluster",
		Short: "Cluster related commands (based on kind k8s cluster)",
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
	if err := shell.NewKindInstaller(common.ShellExec).Check(); err != nil {
		return err
	}
	if err := shell.NewKubectlInstaller(common.ShellExec).Check(); err != nil {
		return err
	}
	//if err := shell.NewHelmlInstaller(common.ShellExec).Check(); err != nil {
	//	return err
	//}
	return nil
}

func (k *clusterCmd) initFlags() error {
	return nil
}

func (k *clusterCmd) initSubCommands() error {

	// Kind Commands
	k.initClusterCommands()

	return nil
}

func (k *clusterCmd) initClusterCommands() {
	opts := k.opts.CmdRootOptions

	k.cobraCmd.AddCommand(NewCreateCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewDashboardCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewDeleteCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewStatusCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewDeployCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewListCmd(opts).GetCobraCmd())
}

func loadKubeConfig() error {
	// Export k8s configuration
	loadCmd := "kind get kubeconfig-path --name=" + common.GlobalOptions.KindClusterName
	if out, err := common.ShellExec.ExecuteWithOutput(loadCmd); err != nil {
		return err
	} else {
		out = strings.TrimSuffix(out, "\n")
		return os.Setenv("KUBECONFIG", out)
	}
}

func checkForActiveCluster(name string) (bool, error) {
	getClustersCmd := "kind get clusters"
	if out, err := common.ShellExec.ExecuteWithOutput(getClustersCmd); err != nil {
		return false, err
	} else {
		contains := strings.Contains(out, name)
		return contains, nil
	}
}
