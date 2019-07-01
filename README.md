![anchor-side](assets/anchor-logo-side-600px.png)

## Overview
Anchor is a utility intended for managing a volatile local Docker / Kubernetes development environment with ease.
<br> It is based on:
- [Docker](https://github.com/docker) for running containerized applications.
- [Kind](https://github.com/kubernetes-sigs/kind) for running local K8s cluster as a docker container.
- [HomeBrew](https://github.com/Homebrew/brew) for managing maxOS/Linux packages.
- [Kubectl](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubectl) for running commands against Kubernetes clusters. 

## Why
1. Development environment should strive to be the same as production
2. Avoid clutter by consolidate into a single folder all the docker build/run/tag/push commands and kubernetes manifests   
3. Easily manage and expand your kubernetes supported resources 
4. Utilize simple commands that encapsulate your repetitive docker and kubernetes cli commands

## How does it work?
Anchor act as one stop shop for all Dockerfiles & Kubernetes manifests that comprise your development environment.
It relies on a `DOCKER_FILES` ENV variable that points to a local directory path containing the following structure:
. <br>
├── ... <br>
├── nginx                   # Name of the docker image/container <br> 
│   ├── k8s                 # Kubernetes content <br>
│   │   ├── manifest.yaml   # Kubernetes manifest <br>
│   ├── Dockerfile          # Docker build/run/tag/push instructions <br>
│   ├── .env                # Optional: Override root `NAMESPACE` & `TAG` environment vars <br> 
│   └── ...                 # Optional: files for docker build <br>
└── ... <br>
├── .env                    # Optional: Override default `NAMESPACE` & `TAG` environment vars at root level 

> It is recommended to back the `DOCKER_FILES` directory by a git repository

## Requirements
- Go 1.12.x

## Getting Started

List of available `anchor docker` commands:

1. `list`  - List all available images
2. `build` - Builds an image by name
3. `clean` - Cleaning unknown and previous container images by name
4. `push`  - Push dockerfile image to repository by name
5. `run`   - Run a Dockerfile by name
6. `stop`  - Stop container by name

List of available `anchor kind` commands:

1. `create`    - Create a Kubernetes cluster 
2. `dashboard` - Deploy a Kubernetes dashboard pod to allow Web UI 

#### Examples:

List and run a container for mysql:

```bash
~$ anchor docker list
----------------------- Listing all Docker images ------------------------
  alpine
  automation
  centos
  grafana
  jenkins
  jfrog-artifactory
  kafka
  mysql-percona-client
  mysql
  percona-server
  postgresql-percona-client
  postgresql-pgwatch2
  postgresql
  zookeeper

    Done.

~$ anchor docker build mysql && anchor docker run mysql
```
---

Retrieve built container images:

```bash
~$ docker images

REPOSITORY                                TAG                 IMAGE ID            CREATED             SIZE
anchor/centos                             latest              1aaabbbcccdd        4 weeks ago         202MB
anchor/mysql                              latest              2aaabbbcccdd        4 weeks ago         477MB
anchor/mysql-percona-client               latest              3aaabbbcccdd        4 weeks ago         396MB
anchor/postgresql-percona-client          latest              4aaabbbcccdd        4 weeks ago         396MB
anchor/postgresql                         latest              5aaabbbcccdd        4 weeks ago         312MB
anchor/percona-server                     latest              6aaabbbcccdd        4 weeks ago         1.08GB
anchor/pgwatch2                           latest              7aaabbbcccdd        4 weeks ago         1.06GB
anchor/grafana                            latest              8aaabbbcccdd        5 weeks ago         244MB
anchor/automation                         latest              9aaabbbcccdd        6 weeks ago         1.97GB
```

