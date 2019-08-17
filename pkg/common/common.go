package common

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

var ShellExec shell.Shell

var GlobalOptions = CmdRootOptions{
	Verbose: false,

	// Docker
	DockerRegistryDns:       "registry.anchor",
	DockerRegistryDnsWithIp: "registry.anchor:32001",
	DockerImageNamespace:    "anchor",
	DockerImageTag:          "latest",
	DockerRepositoryPath:    "",

	// Kind
	KindClusterName: "anchor",
}

type CmdRootOptions struct {
	ConfigFile string

	// Log options
	Verbose bool

	// Docker
	DockerRegistryDns       string
	DockerRegistryDnsWithIp string
	DockerImageNamespace    string
	DockerRepositoryPath    string
	DockerImageTag          string

	// Kind
	KindClusterName string
}
