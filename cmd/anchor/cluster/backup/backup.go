package backup

import (
	"github.com/ZachiNachshon/anchor/pkg/cluster"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/spf13/cobra"
)

var namespace = common.GlobalOptions.DockerImageNamespace

type backupCmd struct {
	cobraCmd *cobra.Command
	opts     BackupOptions
}

type BackupOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

const longDescription = `
Backup a mounted volume from <node>:/opt/stateful/<app-name> to ${HOME}/.anchor/<app-name>

Example usage:

  - anchor cluster backup <app-name>

Note:
<app-name> must exist as a DOCKER_FILES repository resource
`

func NewCommand(opts *common.CmdRootOptions) *backupCmd {
	var cobraCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup a stateful mounted volume",
		Long:  longDescription,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cluster.CheckEnvironment(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintHeadline(logger.ClusterHeadline, "Backup")
			if valid := cluster.Prerequisites(); !valid {
				return
			}
			if _, err := cluster.Backup(args[0], namespace); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var backupCmd = new(backupCmd)
	backupCmd.cobraCmd = cobraCmd
	backupCmd.opts.CmdRootOptions = opts

	if err := backupCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return backupCmd
}

func (cmd *backupCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *backupCmd) initFlags() error {
	cmd.cobraCmd.PersistentFlags().StringVarP(
		&namespace,
		"namespace",
		"n",
		namespace,
		"anchor cluster backup <app-name> -n <namespace>")
	return nil
}
