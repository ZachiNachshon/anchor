## 0.2.0 (Unreleased)
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
