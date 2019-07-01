package docker

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
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
		Short: fmt.Sprintf("Push a docker image to repository [%v]", common.GlobalOptions.DockerRegistryDns),
		Long:  fmt.Sprintf("Push a docker image to repository [%v]", common.GlobalOptions.DockerRegistryDns),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Push: Docker Image")

			if dockerfilePath, err := tagDockerImage(args[0]); err != nil {
				logger.Fatal(err.Error())
			} else {
				if err := pushDockerImage(args[0], dockerfilePath); err != nil {
					logger.Fatal(err.Error())
				}
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

func tagDockerImage(dirname string) (string, error) {
	if dockerfilePath, err := getDockerfileContextPath(dirname); err != nil {
		return "", err
	} else {
		logger.Info("\n==> Tagging image...\n")
		if tagCmd, err := extractDockerCmd(dockerfilePath, DockerCommandTag); err != nil {
			return "", err
		} else if len(tagCmd) == 0 {
			return "", errors.Errorf(missingDockerCmdMsg(DockerCommandTag, dirname))
		} else {
			if common.GlobalOptions.Verbose {
				logger.Info("\n" + tagCmd)
			}
			if err = common.ShellExec.Execute(tagCmd); err != nil {
				return "", err
			}
			logger.Info(" Successfully tagged.")
			return dockerfilePath, nil
		}
	}
}

func pushDockerImage(dirname string, dockerfilePath string) error {
	logger.Info("\n==> Pushing image...\n")
	if pushCmd, err := extractDockerCmd(dockerfilePath, DockerCommandPush); err != nil {
		return err
	} else if len(pushCmd) == 0 {
		return errors.Errorf(missingDockerCmdMsg(DockerCommandPush, dirname))
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + pushCmd)
		}
		if err = common.ShellExec.Execute(pushCmd); err != nil {
			return err
		}

		_ = untagDockerImage(pushCmd)

		logger.Info(" Successfully pushed to registry.")
	}

	return nil
}

func untagDockerImage(pushCommand string) error {
	removeImage := strings.Replace(pushCommand, "push", "rmi -f", 1)
	if err := common.ShellExec.Execute(removeImage); err != nil {
		return err
	}
	return nil
}

func (cmd *pushCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *pushCmd) initFlags() error {
	return nil
}
