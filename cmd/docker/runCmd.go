package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/spf13/cobra"
)

type runCmd struct {
	cobraCmd *cobra.Command
	opts     RunCmdOptions
}

type RunCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewRunCmd(opts *common.CmdRootOptions) *runCmd {
	var cobraCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a docker container",
		Long:  `Run a docker container`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Run: Docker Container on LOCAL Machine")

			if err := stopContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := removeContainers(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			if err := runContainer(args[0]); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var runCmd = new(runCmd)
	runCmd.cobraCmd = cobraCmd
	runCmd.opts.CmdRootOptions = opts

	if err := runCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return runCmd
}

func runContainer(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Running container [%v]", name))
	}

	if runCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandRun); err != nil {
		return err
	} else {
		if common.GlobalOptions.Verbose {
			logger.Info("\n" + runCmd + "\n")
		}

		if out, err := common.ShellExec.ExecuteWithOutput(runCmd); err != nil {
			logger.Info(out)
			return err
		} else {
			tailCmd := fmt.Sprintf("docker logs -f %v", out)
			if err := common.ShellExec.Execute(tailCmd); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cmd *runCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *runCmd) initFlags() error {
	return nil
}
