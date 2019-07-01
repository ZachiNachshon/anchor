package common

import "github.com/anchor/pkg/utils/shell"

var ShellExec shell.Shell

var GlobalOptions = CmdRootOptions{
	Verbose: false,

	// Docker
	DockerRegistryDns:    "registry.anchor:32001",
	DockerImageNamespace: "anchor",
	DockerImageTag:       "latest",
	DockerRepositoryPath: "",

	// Kind
	KindClusterName: "anchor",

	// Remote Files
	DashboardManifest:        "https://raw.githubusercontent.com/ZachiNachshon/anchor-files/master/v1/dashboard/dashboard.yaml",
	RegistryManifest:         "https://raw.githubusercontent.com/ZachiNachshon/anchor-files/master/v1/docker-registry/registry.yaml",
	RegistryContainerdConfig: "https://raw.githubusercontent.com/ZachiNachshon/anchor-files/master/v1/docker-registry/config_template.toml",
}

type CmdRootOptions struct {
	ConfigFile string

	// Log options
	Verbose bool

	// Docker
	DockerRegistryDns    string
	DockerImageNamespace string
	DockerRepositoryPath string
	DockerImageTag       string

	// Kind
	KindClusterName string

	// Remote Files
	DashboardManifest        string
	RegistryManifest         string
	RegistryContainerdConfig string
}
