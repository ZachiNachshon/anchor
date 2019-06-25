package common

import "github.com/kit/pkg/utils/shell"

var ShellExec shell.Shell

var GlobalOptions = CmdRootOptions{
	Verbose: false,

	// Docker
	DockerImageNamespace: "znkit",
	DockerImageTag:       "latest",
	DockerRepositoryPath: "",

	// Kind
	KindClusterName: "znkit",
}

type CmdRootOptions struct {
	ConfigFile string

	// Log options
	Verbose bool

	// Docker
	DockerImageNamespace string
	DockerRepositoryPath string
	DockerImageTag       string

	// Kind
	KindClusterName string
}
