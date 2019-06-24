package common

import "github.com/kit/pkg/utils/shell"

var ShellExec shell.Shell

var GlobalOptions = CmdRootOptions{
	Verbose:              false,
	DockerImageNamespace: "znkit",
	DockerImageTag:       "latest",
	DockerRepositoryPath: "",
}

type CmdRootOptions struct {
	ConfigFile string

	// Log options
	Verbose bool

	// Additional Params
	DockerImageNamespace string
	DockerRepositoryPath string
	DockerImageTag       string
}
