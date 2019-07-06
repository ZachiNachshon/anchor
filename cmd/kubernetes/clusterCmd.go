package kubernetes

import (
	"github.com/anchor/config"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/installer"
	"github.com/anchor/pkg/utils/locator"
	"github.com/anchor/pkg/utils/parser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

type ManifestCommand string

const (
	ManifestCommandPortForward ManifestCommand = "kubectl port-forward"
	ManifestCommandWait        ManifestCommand = "kubectl wait"
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
		Use:     "cluster",
		Short:   "Cluster commands",
		Aliases: []string{"c"},
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
	if err := installer.NewKubectlInstaller(common.ShellExec).Check(); err != nil {
		return err
	}
	if err := installer.NewEnvsubstInstaller(common.ShellExec).Check(); err != nil {
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
	k.cobraCmd.AddCommand(NewRemoveCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewListCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewRegistryCmd(opts).GetCobraCmd())
	k.cobraCmd.AddCommand(NewExposeCmd(opts).GetCobraCmd())
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

func extractManifestCmd(manifestFilePath string, manifestCommand ManifestCommand) (string, error) {
	var result = ""
	if contentByte, err := ioutil.ReadFile(manifestFilePath); err != nil {
		return "", err
	} else {

		l := locator.NewLocator()
		dirPath := l.GetRootFromManifestFile(manifestFilePath)

		// Load .env files from DOCKER_FILES location at root directory and then override from child
		config.LoadEnvVars(dirPath)

		var dockerfileContent = string(contentByte)

		p := parser.NewHashtagParser()
		if err := p.Parse(dockerfileContent); err != nil {
			return "", errors.Errorf("Failed to parse: %v, err: %v", manifestFilePath, err.Error())
		}

		if result = p.Find(string(manifestCommand)); result != "" {
			// In the future might manually substitute arguments
			result = strings.TrimSuffix(result, "\n")
		}
	}

	return result, nil
}
