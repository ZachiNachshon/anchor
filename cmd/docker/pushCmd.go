package docker

import (
	"fmt"
	"github.com/anchor/config"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type pushCmd struct {
	cobraCmd *cobra.Command
	opts     PushCmdOptions
}

type PushCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewPushCmd(opts *common.CmdRootOptions) *pushCmd {
	var cobraCmd = &cobra.Command{
		Use:   "push",
		Short: "Push a docker image to remote/local repository",
		Long:  `Push a docker image to remote/local repository`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Push Docker Image")
			if err := pushDockerImage(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var pushCmd = new(pushCmd)
	pushCmd.cobraCmd = cobraCmd
	pushCmd.opts.CmdRootOptions = opts

	if err := pushCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return pushCmd
}

func pushDockerImage(dirname string) error {
	if dockerfilePath, err := getDockerfileContextPath(dirname); err != nil {
		return err
	} else {
		dirPath := filepath.Dir(dockerfilePath)
		config.LoadEnvVars(dirPath)

		registry := os.Getenv("REGISTRY")
		namespace := os.Getenv("NAMESPACE")
		tag := os.Getenv("TAG")

		pushCmd := fmt.Sprintf("docker push %v/%v/%v:%v", registry, namespace, dirname, tag)
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + pushCmd + "\n")
		}
		if err = common.ShellExec.Execute(pushCmd); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *pushCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *pushCmd) initFlags() error {
	return nil
}
