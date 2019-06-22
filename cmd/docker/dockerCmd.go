package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
	"github.com/kit/pkg/utils/parser"
	"github.com/kit/pkg/utils/shell"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var DOCKER_IMAGE_NAMESPACE = "znkit"
var DOCKER_FILES_REPO_PATH string
var DOCKER_IMAGE_TAG = "latest"

type DockerCommand string

var shellExec shell.Shell

const (
	DockerCommandRun   DockerCommand = "docker run"
	DockerCommandBuild DockerCommand = "docker build"
)

func init() {
	if prefix := os.Getenv("DOCKER_IMAGE_NAMESPACE"); len(prefix) > 0 {
		DOCKER_IMAGE_NAMESPACE = prefix
	}

	shellExec = shell.NewShellExecutor(shell.BASH)
}

type dockerCmd struct {
	cobraCmd *cobra.Command
	opts     DockerCmdOptions
}

type DockerCmdOptions struct {
	*common.CmdRootOptions

	// Additional Build Params
}

func NewDockerCmd(opts *common.CmdRootOptions) *dockerCmd {
	var cobraCmd = &cobra.Command{
		Use:   "docker",
		Short: "Docker related commands",
	}

	var dockerCmd = new(dockerCmd)
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

func (cmd *dockerCmd) GetCobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func checkPrerequisites() error {
	if DOCKER_FILES_REPO_PATH = os.Getenv("DOCKER_FILES"); len(DOCKER_FILES_REPO_PATH) <= 0 {
		return errors.Errorf("DOCKER_FILES environment variable is missing, must contain path to the 'dockerfiles' git repository.")
	}

	if err := shell.NewDockerInstaller().Check(); err != nil {
		return err
	}

	return nil
}

func (cmd *dockerCmd) initFlags() error {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	return nil
}

func (root *dockerCmd) initSubCommands() error {

	// Docker Commands
	root.initDockerCommands()

	return nil
}

func (docker *dockerCmd) initDockerCommands() {
	opts := docker.opts.CmdRootOptions

	docker.cobraCmd.AddCommand(NewBuildCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewCleanCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewListCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewPushCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewRunCmd(opts).GetCobraCmd())
	docker.cobraCmd.AddCommand(NewStopCmd(opts).GetCobraCmd())
}

func composeDockerImageIdentifier(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v:%v", DOCKER_IMAGE_NAMESPACE, dirname, DOCKER_IMAGE_TAG)
	return imageIdentifier
}

func composeDockerImageIdentifierNoTag(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v", DOCKER_IMAGE_NAMESPACE, dirname)
	return imageIdentifier
}

func extractDockerCmd(dockerfilePath string, dockerCommand DockerCommand) (string, error) {
	var dockerfileContent = ""
	if contentByte, err := ioutil.ReadFile(dockerfilePath); err != nil {
		return "", err
	} else {
		dockerfileContent = string(contentByte)

		p := parser.NewHashtagParser()
		if err := p.Parse(dockerfileContent); err != nil {
			logger.Fatalf("Failed to parse: %v, err: %v", dockerfilePath, err.Error())
		}

		if cmd := p.Find(string(dockerCommand)); cmd != "" {
			dockerfileContent = replaceDockerCommandPlaceholders(cmd, dockerfilePath)
		}
	}
	return dockerfileContent, nil
}

func replaceDockerCommandPlaceholders(content string, path string) string {
	content = strings.ReplaceAll(content, "NAMESPACE", DOCKER_IMAGE_NAMESPACE)
	content = strings.ReplaceAll(content, "TAG", DOCKER_IMAGE_TAG)
	content = strings.ReplaceAll(content, "DOCKERFILE", path)
	return content
}
