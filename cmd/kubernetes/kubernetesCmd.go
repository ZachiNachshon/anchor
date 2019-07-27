package kubernetes

import (
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/installer"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type kubernetesCmd struct {
	cobraCmd *cobra.Command
	opts     KubernetesCmdOptions
}

type KubernetesCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewKubernetesCmd(opts *common.CmdRootOptions) *kubernetesCmd {
	var cobraCmd = &cobra.Command{
		Use:     "kubernetes",
		Short:   "Kubernetes commands",
		Aliases: []string{"k"},
	}

	var k8sCmd = new(kubernetesCmd)
	k8sCmd.cobraCmd = cobraCmd
	k8sCmd.opts.CmdRootOptions = opts

	if err := checkPrerequisites(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := k8sCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := k8sCmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return k8sCmd
}

func (cmd *kubernetesCmd) GetCobraCmd() *cobra.Command {
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

func (k *kubernetesCmd) initFlags() error {
	return nil
}

func (k *kubernetesCmd) initSubCommands() error {

	// Kind Commands
	k.initClusterCommands()

	return nil
}

func (k *kubernetesCmd) initClusterCommands() {
	opts := k.opts.CmdRootOptions

	k.cobraCmd.AddCommand(NewCreateCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewDashboardCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewDestroyCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewStatusCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewDeployCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewRemoveCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewRegistryCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewExposeCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewConnectCmd(opts).GetCobraCmd())
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
