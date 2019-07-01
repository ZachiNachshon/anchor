package docker

import (
	"fmt"
	"github.com/anchor/pkg/utils/installer"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/anchor/config"
	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/parser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type DockerCommand string

const (
	DockerCommandRun   DockerCommand = "docker run"
	DockerCommandBuild DockerCommand = "docker build"
	DockerCommandPush  DockerCommand = "docker push"
	DockerCommandTag   DockerCommand = "docker tag"
)

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
		Use:     "docker",
		Short:   "Docker commands",
		Aliases: []string{"d"},
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
	if err := installer.NewDockerInstaller(common.ShellExec).Check(); err != nil {
		return err
	}

	return nil
}

func (cmd *dockerCmd) initFlags() error {
	return nil
}

func (d *dockerCmd) initSubCommands() error {

	// Docker Commands
	d.initDockerCommands()

	return nil
}

func (d *dockerCmd) initDockerCommands() {
	opts := d.opts.CmdRootOptions

	d.cobraCmd.AddCommand(NewBuildCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewCleanCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewListCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewPushCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewRunCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewStopCmd(opts).GetCobraCmd())
	d.cobraCmd.AddCommand(NewPurgeCmd(opts).GetCobraCmd())
}

func composeDockerImageIdentifier(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v:%v", common.GlobalOptions.DockerImageNamespace, dirname, common.GlobalOptions.DockerImageTag)
	return imageIdentifier
}

func composeDockerImageIdentifierNoTag(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v-%v", common.GlobalOptions.DockerImageNamespace, dirname)
	return imageIdentifier
}

func extractDockerCmd(dockerfilePath string, dockerCommand DockerCommand) (string, error) {
	var result = ""
	if contentByte, err := ioutil.ReadFile(dockerfilePath); err != nil {
		return "", err
	} else {

		dirPath := filepath.Dir(dockerfilePath)

		// Load .env files from DOCKER_FILES location at root directory and then override from child
		config.LoadEnvVars(dirPath)

		var dockerfileContent = string(contentByte)

		p := parser.NewHashtagParser()
		if err := p.Parse(dockerfileContent); err != nil {
			return "", errors.Errorf("Failed to parse: %v, err: %v", dockerfilePath, err.Error())
		}

		if result = p.Find(string(dockerCommand)); result != "" {
			if dockerCommand == DockerCommandBuild {
				result = replaceDockerCommandPlaceholders(result, dockerfilePath)
			}
		}
	}
	return result, nil
}

func replaceDockerCommandPlaceholders(content string, path string) string {
	content = strings.ReplaceAll(content, "Dockerfile", path)
	return content
}

func missingDockerCmdMsg(command DockerCommand, dirname string) string {
	return fmt.Sprintf("Missing '%v' on %v Dockerfile instructions\n", command, dirname)
}
