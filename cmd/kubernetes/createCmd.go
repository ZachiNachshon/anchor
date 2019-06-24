package kubernetes

import (
	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
	"github.com/spf13/cobra"
)

type createCmd struct {
	cobraCmd *cobra.Command
	opts     BuildCmdOptions
}

type BuildCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewCreateCmd(opts *common.CmdRootOptions) *createCmd {
	var cobraCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a local Kubernetes cluster",
		Long:  `Create a local Kubernetes cluster based on Kind.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Creating Kubernetes Cluster")
			//if err := createKubernetesCluster(); err != nil {
			//	logger.Fatal(err.Error())
			//}
			logger.PrintCompletion()
		},
	}

	var createCmd = new(createCmd)
	createCmd.cobraCmd = cobraCmd
	createCmd.opts.CmdRootOptions = opts

	if err := createCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return createCmd
}

func (cmd *createCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *createCmd) initFlags() error {
	cmd.cobraCmd.Flags().StringVarP(&DOCKER_IMAGE_TAG, "DOCKER_IMAGE_TAG", "s", "latest", "docker image DOCKER_IMAGE_TAG")
	return nil
}

//func createKubernetesCluster() error {
//	if exists, err := checkForActiveCluster(); err != nil {
//		return err
//	} else if !exists {
//		if createCmd, err := extractDockerCmd(dockerfilePath, DockerCommandBuild); err != nil {
//			return err
//		} else {
//			dirPath := filepath.Dir(dockerfilePath)
//			ctxIdx := strings.LastIndex(createCmd, ".")
//			createCmd = createCmd[:ctxIdx]
//			createCmd += dirPath
//			if common.GlobalOptions.Verbose {
//				logger.Info("\n" + createCmd + "\n")
//			}
//			if err = shellExec.Execute(createCmd); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
//
//func checkForActiveCluster() (bool, error) {
//
//}
