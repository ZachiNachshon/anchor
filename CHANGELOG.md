## 0.4.0 (Unreleased)

NOTES:
* Updated supported Go version to 1.13.x 

FEATURES: 
* Introduced new auto-log to control docker run auto log

## 0.3.0 (August 28, 2019)

NOTES:
* Added a quick-start script

FEATURES:
* Added new cluster commands: `backup`, `connect node`, `connect pod`, `logs`, `token`
* During Kind cluster creation, ask for node workers count, defaults to 1
* K8s manifest now can specify a `hostPath` to sync a host volume to Kind container which mounts to a running pod volume
* Support stateful deployments via nodeSelector - `<app-name>: anchor-stateful`
* Apply command now identify if deployment is stateful (mounted volume) and ask for a node for deployment
* Delete command now identify if deployment is stateful (mounted volume), backup mounted volume and delete content from pod after removal
* Always backup mounted stateful volume to host path volume before deletion
* Added new flag support for cluster commands: `-namespace` / `-n` flag
* Added new library `logrusorgru/aurora` for log colors + mod vendor
* On every anchor action verify and create if needed a `${HOME}/.anchor` directory

ENHANCEMENTS:
* Changed codebase structure with cmd/pkg logic separation
* Added pod/node selector to allow SSH like behaviour with selection menu
* Cluster log command tails the running pod logs without exec into it
* WaitForInput now allow defaults by allowing new line character
* Added `ManifestContent` that extract content from manifest rather than from the comment header
* Load dashboard secret to clipboard on every dashboard command execution

BUG FIXES:
* Dashboard & registry commands verify for `kubectl` port forwarding and re-run is required
* Kill all `kubectl` port forwarding on cluster destroy
* Kill `kubectl` port forwarding if needed on k8s delete
* Setting `GOPATH` & `GOROOT` to default Homebrew locations is missing
* Prerequisites now checks if there are stopped Kind containers that should be started again after macOS restart
* Missing anchor cluster name on printConfiguration
* Tag was missing integer selection to resource name replacement

## 0.2.0 (August 10, 2019)

NOTES:
* Changed cluster delete command name to destroy
* Dashboard & registry k8s manifest are now embedded within the anchor binary instead of relying on a raw git content

FEATURES:
* Cluster creation now triggers K8s Namespace deployment under cluster name
* On Kind cluster creation, an anchor dedicated namespace is being created, can be overridden by `${NAMESPACE}` env var
* Added cluster connect command for accessing a running pod by easy selection mod
* Added new kubernetes command that loads `KUBECONFIG` export script into clipboard
* Added Go lang installer as mandatory requirement
* Introduce `hostess` library for managing /etc/hosts on private docker registry creation/removal
* Created `bash` / `zsh` auto completion script generator command
* Changed `cluster` command name to `kubernetes` with alias `k`
* Change docker `clean` command name to `remove`

ENHANCEMENTS:
* Docker registry now prints the entire private registry catalog content including image name and tags list
* Docker build command allows the usage of custom Dockerfile path
* Added a new numeric input implementation

BUG FIXES:
* Docker clean didn't remove images properly
* Clean docker images now properly handle multiple image ids
* `kubectl` should be installed via brew package and not cask
* Missing organization module name in `go.mod`

## 0.1.0 (July 3, 2019)

NOTES:
* This is the 1st pre-release version of the `Ancor` utility

FEATURES:
* Kubernetes local ephemeral cluster based on [Kind](https://github.com/kubernetes-sigs/kind)
* Docker private registry creation with `reigstry.anchor` DNS record support
* Utilize a central repository containing all docker/kubernetes instructions exposed by `DOCKER_FILES` ENV var 

ENHANCEMENTS:
* Anchor docker commands: `build`, `clean`, `list`, `purge`, `push`, `run`, `stop`
* Anchor cluster commands: `create`, `dashboard`, `delete`, `deploy`, `expose`, `list`, `registry`, `remove`, `status`
