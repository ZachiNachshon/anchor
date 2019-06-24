package kubernetes

import (
	"os"

	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
	"github.com/kit/pkg/utils/shell"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var DOCKER_IMAGE_NAMESPACE = "znkit"
var DOCKER_FILES_REPO_PATH string
var DOCKER_IMAGE_TAG = "latest"

var shellExec shell.Shell

func init() {
	if prefix := os.Getenv("DOCKER_IMAGE_NAMESPACE"); len(prefix) > 0 {
		DOCKER_IMAGE_NAMESPACE = prefix
	}

	shellExec = shell.NewShellExecutor(shell.BASH)
}

type kindCmd struct {
	cobraCmd *cobra.Command
	opts     KindCmdOptions
}

type KindCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewKindCmd(opts *common.CmdRootOptions) *kindCmd {
	var cobraCmd = &cobra.Command{
		Use:   "kind",
		Short: "Kind (k8s cluster) related commands",
	}

	var kindCmd = new(kindCmd)
	kindCmd.cobraCmd = cobraCmd
	kindCmd.opts.CmdRootOptions = opts

	if err := checkPrerequisites(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := kindCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := kindCmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return kindCmd
}

func (cmd *kindCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func checkPrerequisites() error {
	if DOCKER_FILES_REPO_PATH = os.Getenv("DOCKER_FILES"); len(DOCKER_FILES_REPO_PATH) <= 0 {
		return errors.Errorf("DOCKER_FILES environment variable is missing, must contain path to the 'dockerfiles' git repository.")
	}

	if err := shell.NewKindInstaller().Check(); err != nil {
		return err
	}

	return nil
}

func (k *kindCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}

func (k *kindCmd) initSubCommands() error {

	// Docker Commands
	k.initKindCommands()

	return nil
}

func (k *kindCmd) initKindCommands() {
	opts := k.opts.CmdRootOptions

	k.cobraCmd.AddCommand(NewCreateCmd(opts).GetCobraCmd())
	//k.cobraCmd.AddCommand(NewCleanCmd(opts).GetCobraCmd())
	//k.cobraCmd.AddCommand(NewListCmd(opts).GetCobraCmd())
	//k.cobraCmd.AddCommand(NewPushCmd(opts).GetCobraCmd())
	//k.cobraCmd.AddCommand(NewRunCmd(opts).GetCobraCmd())
	//k.cobraCmd.AddCommand(NewStopCmd(opts).GetCobraCmd())
}
