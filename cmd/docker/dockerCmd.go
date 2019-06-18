package docker

import (
	"fmt"
	"github.com/kit/cmd/logger"
	. "github.com/kit/cmd/types"
	"github.com/kit/cmd/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var DOCKER_IMAGE_PREFIX = "znkit"
var DOCKER_FILES_REPO_PATH string

type DockerCommand string

const (
	DockerCommandRun   DockerCommand = "run"
	DockerCommandBuild DockerCommand = "build"
)

func init() {
	if prefix := os.Getenv("DOCKER_IMAGE_PREFIX"); len(prefix) > 0 {
		DOCKER_IMAGE_PREFIX = prefix
	}
}

type DockerCmd struct {
	cobraCmd *cobra.Command
	opts     DockerCmdOptions
}

type DockerCmdOptions struct {
	*CmdRootOptions

	// Additional Build Params
}

func NewDockerCmd(opts *CmdRootOptions) *DockerCmd {
	var cobraCmd = &cobra.Command{
		Use:   "docker",
		Short: "Docker related commands",
	}

	var dockerCmd = new(DockerCmd)
	dockerCmd.cobraCmd = cobraCmd
	dockerCmd.opts.CmdRootOptions = opts

	if err := checkPrerequisites(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := dockerCmd.initFlags(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := dockerCmd.initSubCommands(); err != nil {
		logger.Fatal(err.Error())
	}

	return dockerCmd
}

func (cmd *DockerCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func checkPrerequisites() error {
	if DOCKER_FILES_REPO_PATH = os.Getenv("DOCKER_FILES"); len(DOCKER_FILES_REPO_PATH) <= 0 {
		return errors.Errorf("DOCKER_FILES environment variable is missing, must contain path to the 'dockerfiles' git repository.")
	}

	if _, err := utils.ExecShellWithOutput("which docker"); err != nil {
		return errors.Errorf("docker is missing, must be installed, cannot proceed.")
	}

	return nil
}

func (cmd *DockerCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}

func (root *DockerCmd) initSubCommands() error {

	// Docker Commands
	root.initDockerCommands()

	return nil
}

func (docker *DockerCmd) initDockerCommands() {
	opts := docker.opts.CmdRootOptions

	docker.cobraCmd.AddCommand(NewBuildCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewCleanCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewListCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewPushCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewRunCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewStopCmd(opts).GetCobraCmd())
}

func composeDockerImageIdentifier(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v:%v", DOCKER_IMAGE_PREFIX, dirname, tag)
	return imageIdentifier
}

func composeDockerImageIdentifierNoTag(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v", DOCKER_IMAGE_PREFIX, dirname)
	return imageIdentifier
}

func extractDockerCmd(dockerfilePath string, dockerCommand DockerCommand) (string, error) {
	var contentStr = ""
	if contentByte, err := ioutil.ReadFile(dockerfilePath); err != nil {
		return "", err
	} else {
		contentStr := string(contentByte)

		startIdx := strings.Index(contentStr, "docker "+string(dockerCommand))
		contentStr = contentStr[startIdx:]

		lastIdx := strings.Index(contentStr, "#\n#")
		contentStr = contentStr[:lastIdx]

		contentStr = strings.ReplaceAll(contentStr, "#", "")
	}
	return contentStr, nil
}
