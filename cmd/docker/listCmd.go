package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var DOCKER_FILE_NAME = "/Dockerfile"

type listCmd struct {
	cobraCmd *cobra.Command
	opts     ListCmdOptions
}

type ListCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewListCmd(opts *common.CmdRootOptions) *listCmd {
	var cobraCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available docker supported images from DOCKER_FILES repository",
		Long:  `List all available docker supported images from DOCKER_FILES repository`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Listing all Docker images")

			if _, err := listDockerfilesDirs(true); err != nil {
				logger.Fatal(err.Error())
			}

			logger.PrintCompletion()
		},
	}

	var listCmd = new(listCmd)
	listCmd.cobraCmd = cobraCmd
	listCmd.opts.CmdRootOptions = opts

	if err := listCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	return listCmd
}

func (cmd *listCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *listCmd) initFlags() error {
	return nil
}

func getDockerfileContextPath(dirname string) (string, error) {
	expected := fmt.Sprintf("%v/%v/Dockerfile", common.GlobalOptions.DockerRepositoryPath, dirname)
	dirNames, _ := listDockerfilesDirs(false)

	for _, e := range dirNames {
		if strings.EqualFold(expected, e) {
			return e, nil
		}
	}

	return "", errors.Errorf("Cannot find Dockerfile for %v", dirname)
}

func listDockerfilesDirs(verbose bool) ([]string, error) {
	var dirNames = make([]string, 0)
	err := filepath.Walk(common.GlobalOptions.DockerRepositoryPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Continue to the next path
			if !strings.Contains(path, DOCKER_FILE_NAME) {
				return nil
			}

			if dockerfilePath, err := filepath.Abs(path); err != nil {
				return err
			} else {
				dirName := extractDockerfileDirName(dockerfilePath)

				if verbose {
					logger.Info("  " + dirName)
				}

				dirNames = append(dirNames, dockerfilePath)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return dirNames, nil
}

func extractDockerfileDirName(path string) string {
	dirName := strings.TrimPrefix(path, common.GlobalOptions.DockerRepositoryPath+"/")
	dirName = strings.TrimSuffix(dirName, DOCKER_FILE_NAME)
	return dirName
}
