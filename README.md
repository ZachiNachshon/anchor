![anchor-side](assets/anchor-logo-side-600px.png)

## Overview
Anchor is a utility intended for managing an ephemeral local Docker / Kubernetes development environment with ease.
<br> It is based on:
- [Docker](https://github.com/docker) for running containerized applications
- [Kind](https://github.com/kubernetes-sigs/kind) for running local K8s cluster as a docker container
- [Kubectl](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubectl) for running commands against Kubernetes clusters 
- [Hostess](https://github.com/cbednarski/hostess) for managing your `/etc/hosts` file
- [Homebrew](https://github.com/Homebrew/brew) for managing macOS/Linux packages

[![Build Status](https://travis-ci.com/ZachiNachshon/anchor.svg "Travis CI status")](https://travis-ci.com/ZachiNachshon/anchor)

> Note:<br>
> Anchor is utilizing the following components: `docker`, `kind`, `kubectl`, `envsubst`, `hostess`.<br>
> If they can't be found on your machine, `Homebrew` is being installed and fetches them for you.

## Why
1. Allow a repository to become an anchor for all docker / kubernetes scripts that you manage on local / CI environment
2. Avoid clutter by consolidate into a single repository all docker `build`/`run`/`tag`/`push` commands and kubernetes manifests   
3. Encapsulate commonly used, repetitive docker / kubernetes actions as simple cli commands
4. Development environment should strive to be the same as production, deploy locally the same as you deploy to production

## What's in the box?
- Private docker registry deployed on control-plane with multi node support
- Allow stateful volume `mountPath` / `hostPath` between local and a running pod within the Kind docker via the node selector: `app-name: anchor-stateful`
- Selection menu for quick connect to any node/pod 

## How does it work?
#### Directory Structure & Instructions
`Anchor` relies on a `DOCKER_FILES` environment variable that points to a local directory path containing the dockerfiles and k8s manifests.<br/> 
Please refer to the sample [anchor-dockerfiles](https://github.com/ZachiNachshon/anchor-dockerfiles) repository for additional details.

#### Quick Start Setup 
```bash
~$ sh -c "$(curl -fsSL https://raw.githubusercontent.com/ZachiNachshon/anchor/master/scripts/quick-start.sh)"
```

##### What's included?
1. Clone `anchor-dockerfiles` to `${HOME}/.anchor` 
2. Set `DOCKER_FILES` environment variable to `${HOME}/.anchor/anchor-dockerfiles` 

> Note:<br/>
> Consider setting `DOCKER_FILES` as permanent environment variable (append to `$PATH`)

## Requirements
- Go 1.13.x

## Download

#### I don't have GO environment 
Download your OS and ARCH relevant binary from [releases](https://github.com/ZachiNachshon/anchor/releases), unzip and place in `/usr/bin` or `usr/local/bin`.

#### I do have GO environment

##### without source
```bash
~$ go get github.com/ZachiNachshon/anchor@v0.3.0
```

##### with source
```bash
~$ git clone https://github.com/ZachiNachshon/anchor.git ${GOPATH}/src/github.com/anchor
~$ cd {GOPATH}/src/github.com/anchor
~$ make build
```

## Usage Example 
Create an `anchor` Kubernetes cluster:
[![https://youtu.be/7PtbKPpiJIA](assets/thumbnails/anchor-cluster-create-tn.png)](https://youtu.be/4XCf3M424Gk)

Build and run `nginx` as a docker container:
[![https://youtu.be/7PtbKPpiJIA](assets/thumbnails/anchor-docker-nginx-tn.png)](https://youtu.be/7PtbKPpiJIA)

Auto deploy `nginx` to Kubernetes:
[![https://youtu.be/urmfVmYi5BE](assets/thumbnails/anchor-deploy-auto-nginx-tn.png)](https://youtu.be/7Tdx1GHaQ50)

Manual deploy `nginx` to Kubernetes:
[![https://youtu.be/urmfVmYi5BE](assets/thumbnails/anchor-deploy-manual-nginx-tn.png)](https://youtu.be/urmfVmYi5BE)

Connect to a running Kubernetes pod/node:
[![https://youtu.be/O25weLHGC-M](assets/thumbnails/anchor-cluster-connect-tn.png)](https://youtu.be/O25weLHGC-M)


## Still in progress
- After macOS restart `kubectl` is losing cluster context
- Use `stty sane` to avoid terminal input errors such as `^M` on Enter
- Add the ability to delete an image from private docker registry 
- Add `-y` flag to skip all prompts using default values
- Tests coverage 

## Available Anchor commands

List of available `anchor docker` commands:
```bash
Usage:
  anchor docker [command]

Aliases:
  docker, d

Available Commands:
  build       Builds a docker image
  purge       Purge all docker images and containers
  push        Push a docker image to repository [registry.anchor]
  remove      Remove docker containers and images
  run         Run a docker container
  stop        Stop a docker container

Flags:
  -h, --help   help for docker

Global Flags:
  -v, --verbose   anchor <command> -v
```

List of available `anchor cluster` commands:
```bash
Usage:
  anchor cluster [command]

Aliases:
  cluster, c

Available Commands:
  apply       Apply a Kubernetes manifest resource
  backup      Backup a stateful mounted volume
  connect     Connect to a kubernetes [node, pod] by name
  create      Create a local Kubernetes cluster
  dashboard   Deploy a Kubernetes dashboard
  delete      Delete a previously deployed Kubernetes resource
  deploy      Deploy a fully managed Kubernetes resource
  destroy     Destroy local Kubernetes cluster
  expose      Expose to the host instance a container port of a deployed Kubernetes resource
  logs        Log a running kubernetes pod by name
  registry    Create a private docker registry [registry.anchor]
  status      Print cluster [anchor] status
  token       Generate export KUBECONFIG command and load to clipboard

Flags:
  -h, --help   help for cluster

Global Flags:
  -v, --verbose   anchor <command> -v
```