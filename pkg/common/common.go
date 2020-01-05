package common

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"os"
	"path/filepath"
)

var ShellExec shell.Shell

var GlobalOptions = CmdRootOptions{
	Verbose: false,

	AnchorHomeDirectory: filepath.Join(os.Getenv("HOME"), ".anchor"),

	// Docker
	DockerRegistryDns:       "registry.anchor",
	DockerRegistryDnsWithIp: "registry.anchor:32001",
	DockerImageNamespace:    "anchor",
	DockerImageTag:          "latest",
	DockerRepositoryPath:    "",
	DockerRunAutoLog:        true,

	// Kind
	KindClusterName: "anchor",

	// Go
	GoPathDir: filepath.Join(os.Getenv("HOME"), "go"),
	GoRootDir: "/usr/local/opt/go", // symlink
}

type CmdRootOptions struct {
	ConfigFile string

	AnchorHomeDirectory string

	// Log options
	Verbose bool

	// Docker
	DockerRegistryDns       string
	DockerRegistryDnsWithIp string
	DockerImageNamespace    string
	DockerRepositoryPath    string
	DockerImageTag          string
	DockerRunAutoLog        bool

	// Kind
	KindClusterName string

	// Go
	GoPathDir string
	GoRootDir string
}
