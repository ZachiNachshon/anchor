package docker

import (
	"github.com/kit/cmd/logger"
	"github.com/kit/cmd/types"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type RunCmd struct {
	cobraCmd *cobra.Command
	opts     RunCmdOptions
}

type RunCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewRunCmd(opts *common.CmdRootOptions) *RunCmd {
	var cobraCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a Dockerfile",
		Long:  `Run a Dockerfile from the list of available docker image directories (ex. DIR=alpine).`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Deploying Docker Container on LOCAL Machine")
			if err := runContainer(args[0]); err != nil {
				logger.Fatal(err.Error())
			}
			logger.PrintCompletion()
		},
	}

	var runCmd = new(RunCmd)
	runCmd.cobraCmd = cobraCmd
	runCmd.opts.CmdRootOptions = opts

	if err := runCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return runCmd
}

func runContainer(dirname string) error {
	//imageIdentifier := composeDockerImageIdentifierNoSuite(dirname)

	if dockerfilePath, err := getDockerfileContextPath(dirname); err != nil {
		return err
	} else {
		_ = extractDockerRunCmd(dockerfilePath)
	}

	return nil
}

func extractDockerRunCmd(path string) error {
	if contentByte, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		contentStr := string(contentByte)
		logger.Info(contentStr)
	}
	return nil
}

func (cmd *RunCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *RunCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}
