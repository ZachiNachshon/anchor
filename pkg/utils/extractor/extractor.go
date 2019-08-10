package extractor

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/ZachiNachshon/anchor/pkg/utils/parser"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

var CmdExtractor Extractor

type extractor struct {
}

func New() Extractor {
	return &extractor{}
}

func (e *extractor) DockerCmd(identifier string, dockerCommand DockerCommand) (string, error) {
	if dockerfilePath, err := locator.DirLocator.Dockerfile(identifier); err != nil {
		return "", err
	} else {
		var result = ""
		if contentByte, err := ioutil.ReadFile(dockerfilePath); err != nil {
			return "", err
		} else {
			// Load .env file
			config.LoadEnvVars(identifier)

			var dockerfileContent = string(contentByte)

			p := parser.NewHashtagParser()
			if err := p.Parse(dockerfileContent); err != nil {
				return "", errors.Errorf("Failed to parse: %v, err: %v", dockerfilePath, err.Error())
			}

			if result = p.Find(string(dockerCommand)); result != "" {
				result = strings.TrimSuffix(result, "\n")
				if dockerCommand == DockerCommandBuild {
					result = replaceDockerCommandPlaceholders(result, dockerfilePath)
				}
			}

			if len(result) == 0 {
				return "", errors.Errorf(missingDockerCmdMsg(dockerCommand, identifier))
			}
		}
		return result, nil
	}
}

func replaceDockerCommandPlaceholders(content string, path string) string {
	// In case the Dockerfile is referenced by a custom path
	if strings.Contains(content, "/Dockerfile") {
		return content
	} else {
		content = strings.ReplaceAll(content, "Dockerfile", path)
		return content
	}
}

func missingDockerCmdMsg(command DockerCommand, dirname string) string {
	return fmt.Sprintf("Missing '%v' on %v Dockerfile instructions\n", command, dirname)
}

func (e *extractor) ManifestCmd(identifier string, manifestCommand ManifestCommand) (string, error) {
	if manifestFilePath, err := locator.DirLocator.Manifest(identifier); err != nil {
		return "", err
	} else {
		var result = ""
		if contentByte, err := ioutil.ReadFile(manifestFilePath); err != nil {
			return "", err
		} else {
			// Load .env file
			config.LoadEnvVars(identifier)

			var manifestContent = string(contentByte)

			p := parser.NewHashtagParser()
			if err := p.Parse(manifestContent); err != nil {
				return "", errors.Errorf("Failed to parse: %v, err: %v", manifestFilePath, err.Error())
			}

			if result = p.Find(string(manifestCommand)); result != "" {
				// In the future might manually substitute arguments
				result = strings.TrimSuffix(result, "\n")
			}

			if len(result) == 0 {
				return "", errors.Errorf(missingManifestCmdMsg(manifestCommand, identifier))
			}
		}
		return result, nil
	}
}

func missingManifestCmdMsg(manifestCommand ManifestCommand, dirname string) string {
	return fmt.Sprintf("Missing '%v' on %v K8s manifest\n", manifestCommand, dirname)
}
