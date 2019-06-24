package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kit/pkg/common"
	"github.com/kit/pkg/logger"
	"github.com/kit/pkg/utils/shell"
	"github.com/pkg/errors"
)

func init() {
	if prefix := os.Getenv("DOCKER_IMAGE_NAMESPACE"); len(prefix) > 0 {
		common.GlobalOptions.DockerImageTag = prefix
	}
}

func CheckPrerequisites() error {
	var repoPath = ""
	if repoPath = os.Getenv("DOCKER_FILES"); len(repoPath) <= 0 {
		return errors.Errorf("DOCKER_FILES environment variable is missing, must contain path to 'dockerfiles' git repository.")
	}
	common.GlobalOptions.DockerRepositoryPath = repoPath

	// TODO: resolve shell type from configuration
	common.ShellExec = shell.NewShellExecutor(shell.BASH)

	setDefaultEnvVar()
	LoadEnvVars(common.GlobalOptions.DockerRepositoryPath)

	return nil
}

func setDefaultEnvVar() {
	_ = os.Setenv("NAMESPACE", common.GlobalOptions.DockerImageNamespace)
	_ = os.Setenv("TAG", common.GlobalOptions.DockerImageTag)
}

func LoadEnvVars(path string) {
	envFile := path + "/.env"
	if err := godotenv.Overload(envFile); err != nil {
		// TODO: Change to warn once implemented
		logger.Info(err.Error())
	}
}
