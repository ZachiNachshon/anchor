package docker

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/utils/installer"
	"strings"
)

func RemoveNamespaceFromImageName(name string) string {
	if len(name) == 0 {
		return name
	}

	var result = name
	backslashIdx := strings.LastIndex(name, "/")
	if backslashIdx > 0 {
		result = result[backslashIdx+1:]
	}
	tagIdx := strings.LastIndex(name, ":")
	if tagIdx > 0 {
		result = result[:backslashIdx+1]
	}

	return result
}

func ComposeDockerContainerIdentifier(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v:%v", common.GlobalOptions.DockerImageNamespace, dirname, common.GlobalOptions.DockerImageTag)
	return imageIdentifier
}

func ComposeDockerContainerIdentifierNoTag(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v-%v", common.GlobalOptions.DockerImageNamespace, dirname)
	return imageIdentifier
}

func ComposeDockerImageIdentifierNoTag(dirname string) string {
	imageIdentifier := fmt.Sprintf("%v/%v", common.GlobalOptions.DockerImageNamespace, dirname)
	return imageIdentifier
}

func CheckEnvironment() error {
	if err := installer.NewDockerInstaller(common.ShellExec).Check(); err != nil {
		return err
	}

	return nil
}
