package config

import (
	"os"

	"github.com/anchor/pkg/common"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/shell"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func CheckPrerequisites() error {
	var repoPath = ""
	if repoPath = os.Getenv("DOCKER_FILES"); len(repoPath) <= 0 {
		return errors.Errorf("DOCKER_FILES environment variable is missing, must contain path to 'dockerfiles' git repository.")
	}
	common.GlobalOptions.DockerRepositoryPath = repoPath

	// TODO: resolve shell type from configuration (https://github.com/spf13/viper ?)
	common.ShellExec = shell.NewShellExecutor(shell.BASH)

	setDefaultEnvVar()
	LoadEnvVars(common.GlobalOptions.DockerRepositoryPath)

	return nil
}

func setDefaultEnvVar() {
	// Docker
	_ = os.Setenv("REGISTRY", common.GlobalOptions.DockerRegistryDns)
	_ = os.Setenv("NAMESPACE", common.GlobalOptions.DockerImageNamespace)
	_ = os.Setenv("TAG", common.GlobalOptions.DockerImageTag)
}

func LoadEnvVars(path string) {
	envFile := path + "/.env"
	if err := godotenv.Overload(envFile); err != nil {
		if common.GlobalOptions.Verbose {
			// TODO: Change to warn once implemented
			logger.Info(err.Error())
		}
	}

	if v := os.Getenv("NAMESPACE"); len(v) > 0 {
		common.GlobalOptions.DockerImageNamespace = v
	}

	if v := os.Getenv("TAG"); len(v) > 0 {
		common.GlobalOptions.DockerImageTag = v
	}
}

var KubernetesNamespaceManifest = `
apiVersion: v1
kind: Namespace
metadata:
  name: NAMESPACE-TO-REPLACE
`

var RegistryContainerdConfigTemplate = `disabled_plugins = ["aufs", "btrfs", "zfs"]
root = "/var/lib/containerd"
state = "/run/containerd"
oom_score = 0

[grpc]
  address = "/run/containerd/containerd.sock"
  uid = 0
  gid = 0
  max_recv_message_size = 16777216
  max_send_message_size = 16777216

[debug]
  address = ""
  uid = 0
  gid = 0
  level = ""

[metrics]
  address = ""
  grpc_histogram = false

[cgroup]
  path = ""

[plugins]
  [plugins.cgroups]
    no_prometheus = false
  [plugins.cri]
    stream_server_address = "127.0.0.1"
    stream_server_port = "0"
    enable_selinux = false
    sandbox_image = "k8s.gcr.io/pause:3.1"
    stats_collect_period = 10
    systemd_cgroup = false
    enable_tls_streaming = false
    max_container_log_line_size = 16384
    [plugins.cri.containerd]
      snapshotter = "overlayfs"
      no_pivot = false
      [plugins.cri.containerd.default_runtime]
        runtime_type = "io.containerd.runtime.v1.linux"
        runtime_engine = ""
        runtime_root = ""
      [plugins.cri.containerd.untrusted_workload_runtime]
        runtime_type = ""
        runtime_engine = ""
        runtime_root = ""
    [plugins.cri.cni]
      bin_dir = "/opt/cni/bin"
      conf_dir = "/etc/cni/net.d"
      conf_template = ""
    [plugins.cri.registry]
      [plugins.cri.registry.mirrors]
        [plugins.cri.registry.mirrors."local.insecure-registry.io"]
          endpoint = ["http://127.0.0.1:32001"]
        [plugins.cri.registry.mirrors."registry.anchor:32001"]
          endpoint = ["http://127.0.0.1:32001"]
        [plugins.cri.registry.mirrors."docker.io"]
          endpoint = ["https://registry-1.docker.io"]
    [plugins.cri.x509_key_pair_streaming]
      tls_cert_file = ""
      tls_key_file = ""
  [plugins.diff-service]
    default = ["walking"]
  [plugins.linux]
    shim = "containerd-shim"
    runtime = "runc"
    runtime_root = ""
    no_shim = false
    shim_debug = false
  [plugins.opt]
    path = "/opt/containerd"
  [plugins.restart]
    interval = "10s"
  [plugins.scheduler]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = "0s"
    startup_delay = "100ms"
`
