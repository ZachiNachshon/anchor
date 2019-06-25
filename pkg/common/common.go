package common

import "github.com/anchor/pkg/utils/shell"

var ShellExec shell.Shell

var GlobalOptions = CmdRootOptions{
	Verbose: false,

	// Docker
	DockerImageNamespace: "anchor",
	DockerImageTag:       "latest",
	DockerRepositoryPath: "",

	// Kind
	KindClusterName: "anchor",
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
