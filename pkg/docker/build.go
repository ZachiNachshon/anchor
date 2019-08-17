package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/extractor"
	"github.com/ZachiNachshon/anchor/pkg/utils/input"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const pullAccessDeniedRegexp = "pull access denied for .*, repository does not exist"

func extractImageNameFromPullAccessError(error string) string {
	startIdx := len("pull access denied for ")
	lastIdx := strings.Index(error, ", repository does not exist")
	result := error[startIdx:lastIdx]
	return result
}

func Build(identifier string) error {
	if name, err := locator.DirLocator.Name(identifier); err != nil {
		return err
	} else {
		logger.PrintCommandHeader(fmt.Sprintf("Building image %v", name))
	}

	if buildCmd, err := extractor.CmdExtractor.DockerCmd(identifier, extractor.DockerCommandBuild); err != nil {
		return err
	} else {
		contextPath, _ := locator.DirLocator.DockerContext(identifier)

		// Replace docker build "." with directory absolute path
		ctxIdx := strings.LastIndex(buildCmd, ".")
		buildCmd = buildCmd[:ctxIdx]
		buildCmd += contextPath

		if common.GlobalOptions.Verbose {
			logger.Info("\n" + buildCmd + "\n")
		}

		if err := common.ShellExec.Execute(buildCmd); err != nil {
			// Check is base image is missing
			if out, err := common.ShellExec.ExecuteWithOutput(buildCmd); err != nil {
				// Docker-compose behaviour, look for missing base images from DOCKER_FILES directory
				if matched, err := regexp.MatchString(pullAccessDeniedRegexp, out); err == nil && matched {

					imageName := extractImageNameFromPullAccessError(out)
					imageNoNamespace := RemoveNamespaceFromImageName(imageName)

					in := input.NewYesNoInput()
					q := fmt.Sprintf("Found missing dependant image %v, try to build?", imageNoNamespace)
					if result, err := in.WaitForInput(q); err == nil && result {

						// Build missing base image
						_ = Build(imageNoNamespace)

						// Build previous Dockerfile
						_ = Build(identifier)

					} else {
						return errors.Errorf(out)
					}
				} else {
					return err
				}
			}
		}
	}

	return nil
}
