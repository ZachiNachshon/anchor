## 0.3.0 (Unreleased)
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
