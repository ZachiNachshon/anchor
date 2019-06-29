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

var MANIFESTS_PATH = "k8s/manifest.yaml"

type listCmd struct {
	cobraCmd *cobra.Command
	opts     ListCmdOptions
}

type ListCmdOptions struct {
	*common.CmdRootOptions

	// Additional Params
}

func NewListCmd(opts *common.CmdRootOptions) *listCmd {
	var cobraCmd = &cobra.Command{
		Use:   "list",
		Short: "List all containers with kubernetes manifests",
		Long:  `List all containers with kubernetes manifests`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.PrintHeadline("Listing Containers With K8S Manifests")
			if _, err := listContainersManifestsDirs(true); err != nil {
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

func getContainerManifestsDir(dirname string) (string, error) {
	expected := fmt.Sprintf("%v/%v/%v", common.GlobalOptions.DockerRepositoryPath, dirname, MANIFESTS_PATH)
	dirNames, _ := listContainersManifestsDirs(false)

	for _, e := range dirNames {
		if strings.EqualFold(expected, e) {
			return e, nil
		}
	}

	return "", errors.Errorf("Cannot find container manifests(s) for %v", dirname)
}

func listContainersManifestsDirs(verbose bool) ([]string, error) {
	var dirNames = make([]string, 0)
	err := filepath.Walk(common.GlobalOptions.DockerRepositoryPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Continue to the next path
			if !strings.Contains(path, MANIFESTS_PATH) {
				return nil
			}

			if manifestPath, err := filepath.Abs(path); err != nil {
				return err
			} else {
				dirName := extractManifestDirName(manifestPath)

				if verbose {
					logger.Info("  " + dirName)
				}

				dirNames = append(dirNames, manifestPath)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return dirNames, nil
}

func extractManifestDirName(path string) string {
	// Dirname to k8s
	dirName := filepath.Dir(path)
	// Dirname to container dir name
	dirName = filepath.Dir(dirName)
	dirName = strings.TrimPrefix(dirName, common.GlobalOptions.DockerRepositoryPath+"/")
	return dirName
}
