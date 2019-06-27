package kubernetes

import (
	"fmt"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var RESOURCES_DIR_NAME = "k8s"

type listCmd struct {
	cobraCmd *cobra.Command
	opts     CreateCmdOptions
}

type ListCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewListCmd(opts *common.CmdRootOptions) *listCmd {
	var cobraCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available container kubernetes resources",
		Long:  `List all available container kubernetes resources`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Listing Containers Kubernetes Resources")
			if _, err := listContainersResourceDirs(true); err != nil {
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

func getContainerResourceDir(dirname string) (string, error) {
	expected := fmt.Sprintf("%v/%v/%v", common.GlobalOptions.DockerRepositoryPath, dirname, RESOURCES_DIR_NAME)
	dirNames, _ := listContainersResourceDirs(false)

	for _, e := range dirNames {
		if strings.EqualFold(expected, e) {
			return e, nil
		}
	}

	return "", errors.Errorf("Cannot find container resource(s) for %v", dirname)
}

func listContainersResourceDirs(verbose bool) ([]string, error) {
	var dirNames = make([]string, 0)
	err := filepath.Walk(common.GlobalOptions.DockerRepositoryPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Continue to the next path
			if !info.IsDir() || !strings.Contains(path, RESOURCES_DIR_NAME) {
				return nil
			}

			if resourcePath, err := filepath.Abs(path); err != nil {
				return err
			} else {
				dirName := extractResourceDirName(resourcePath)

				if verbose {
					logger.Info("  " + dirName)
				}

				dirNames = append(dirNames, resourcePath)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return dirNames, nil
}

func extractResourceDirName(path string) string {
	dirName := filepath.Dir(path)
	dirName = strings.TrimPrefix(dirName, common.GlobalOptions.DockerRepositoryPath+"/")
	return dirName
}
