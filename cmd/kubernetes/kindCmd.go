package kubernetes

import (
	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
	"github.com/kit/pkg/utils/shell"
	"github.com/spf13/cobra"
)

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

func (k *kindCmd) initFlags() error {
	return nil
}

func (k *kindCmd) initSubCommands() error {

	// Kind Commands
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
