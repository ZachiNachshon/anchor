package config

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/installer"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"os"

	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
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

	if err := installer.NewGoInstaller(common.ShellExec).Check(); err != nil {
		return err
	}

	return nil
}

func setDefaultEnvVar() {
	// Docker
	_ = os.Setenv("REGISTRY", common.GlobalOptions.DockerRegistryDnsWithIp)
	_ = os.Setenv("NAMESPACE", common.GlobalOptions.DockerImageNamespace)
	_ = os.Setenv("TAG", common.GlobalOptions.DockerImageTag)
}

func LoadEnvVars(identifier string) {
	var envFilePath = ""

	if identifier == common.GlobalOptions.DockerRepositoryPath {
		// dirname should be the repository root folder only at config CheckPrerequisites stage
		envFilePath = common.GlobalOptions.DockerRepositoryPath + "/.env"
	} else {
		ctxDir, _ := locator.DirLocator.DockerContext(identifier)
		envFilePath = ctxDir + "/.env"
	}
	loadEnvVarsInner(envFilePath)
}

func loadEnvVarsInner(envFilePath string) {
	if err := godotenv.Overload(envFilePath); err != nil {
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

const KindClusterManifestFormat = `# Kind cluster creation config
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
nodes:
%v
%v`

const KubernetesNamespaceManifest = `
apiVersion: v1
kind: Namespace
metadata:
  name: NAMESPACE-TO-REPLACE
`

const RegistryContainerdConfigTemplate = `disabled_plugins = ["aufs", "btrfs", "zfs"]
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

const KubernetesDashboardManifest = `# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ------------------- Dashboard Secret ------------------- #

apiVersion: v1
kind: Secret
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-certs
  namespace: kube-system
type: Opaque

---
# ------------------- Dashboard Service Account ------------------- #

apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system

---
# ------------------- Dashboard Role & Role Binding ------------------- #

kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kubernetes-dashboard-minimal
  namespace: kube-system
rules:
  # Allow Dashboard to create 'kubernetes-dashboard-key-holder' secret.
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["create"]
    # Allow Dashboard to create 'kubernetes-dashboard-settings' config map.
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["create"]
    # Allow Dashboard to get, update and delete Dashboard exclusive secrets.
  - apiGroups: [""]
    resources: ["secrets"]
    resourceNames: ["kubernetes-dashboard-key-holder", "kubernetes-dashboard-certs"]
    verbs: ["get", "update", "delete"]
    # Allow Dashboard to get and update 'kubernetes-dashboard-settings' config map.
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["kubernetes-dashboard-settings"]
    verbs: ["get", "update"]
    # Allow Dashboard to get metrics from heapster.
  - apiGroups: [""]
    resources: ["services"]
    resourceNames: ["heapster"]
    verbs: ["proxy"]
  - apiGroups: [""]
    resources: ["services/proxy"]
    resourceNames: ["heapster", "http:heapster:", "https:heapster:"]
    verbs: ["get"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kubernetes-dashboard-minimal
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kubernetes-dashboard-minimal
subjects:
  - kind: ServiceAccount
    name: kubernetes-dashboard
    namespace: kube-system

---
# ------------------- Dashboard Deployment ------------------- #

kind: Deployment
apiVersion: apps/v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s-app: kubernetes-dashboard
  template:
    metadata:
      labels:
        k8s-app: kubernetes-dashboard
    spec:
      containers:
        - name: kubernetes-dashboard
          image: k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1
          ports:
            - containerPort: 8443
              protocol: TCP
          args:
            - --auto-generate-certificates
            # Uncomment the following line to manually specify Kubernetes API server Host
            # If not specified, Dashboard will attempt to auto discover the API server and connect
            # to it. Uncomment only if the default does not work.
            # - --apiserver-host=http://my-address:port
          volumeMounts:
            - name: kubernetes-dashboard-certs
              mountPath: /certs
              # Create on-disk volume to store exec logs
            - mountPath: /tmp
              name: tmp-volume
          livenessProbe:
            httpGet:
              scheme: HTTPS
              path: /
              port: 8443
            initialDelaySeconds: 30
            timeoutSeconds: 30
      volumes:
        - name: kubernetes-dashboard-certs
          secret:
            secretName: kubernetes-dashboard-certs
        - name: tmp-volume
          emptyDir: {}
      serviceAccountName: kubernetes-dashboard
      # Comment the following tolerations if Dashboard must not be deployed on master
      tolerations:
        - key: node-role.kubernetes.io/master
          effect: NoSchedule

---
# ------------------- Dashboard Service ------------------- #

kind: Service
apiVersion: v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  ports:
    - port: 443
      targetPort: 8443
  selector:
    k8s-app: kubernetes-dashboard`

const KubernetesRegistryManifest = `---
apiVersion: v1
kind: Namespace
metadata:
  name: container-registry
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: registry-claim
  namespace: container-registry
spec:
  accessModes:
    - ReadWriteMany
  volumeMode: Filesystem
  resources:
    requests:
      storage: 5Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: registry
  name: registry
  namespace: container-registry
spec:
  replicas: 1
  selector:
    matchLabels:
      app: registry
  template:
    metadata:
      labels:
        app: registry
    spec:
      containers:
        - name: registry
          image: cdkbot/registry-amd64:2.6
          env:
            - name: REGISTRY_HTTP_ADDR
              value: :5000
            - name: REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY
              value: /var/lib/registry
            - name: REGISTRY_STORAGE_DELETE_ENABLED
              value: "yes"
          ports:
            - containerPort: 5000
              name: registry
              protocol: TCP
          volumeMounts:
            - mountPath: /var/lib/registry
              name: registry-data
      volumes:
        - name: registry-data
          persistentVolumeClaim:
            claimName: registry-claim
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: registry
  name: registry
  namespace: container-registry
spec:
  type: NodePort
  selector:
    app: registry
  ports:
    - name: "registry"
      port: 5000
      targetPort: 5000
      nodePort: 32001`
