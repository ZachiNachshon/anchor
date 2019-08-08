package docker

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/extractor"
	"github.com/anchor/pkg/utils/locator"
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

			if err := tagDockerImage(args[0]); err != nil {
				logger.Fatal(err.Error())
			} else {
				if err := pushDockerImage(args[0]); err != nil {
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

func tagDockerImage(identifier string) error {
	logger.PrintCommandHeader(fmt.Sprintf("Tagging image %v", identifier))
	if tagCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandTag); err != nil {
		return err
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + tagCmd + "\n")
		}
		if err = common.ShellExec.Execute(tagCmd); err != nil {
			return err
		}
		logger.Info(" Successfully tagged.")
		return nil
	}
}

func pushDockerImage(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Pushing image [%v]", name))
	}

	if pushCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandPush); err != nil {
		return err
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + pushCmd + "\n")
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
