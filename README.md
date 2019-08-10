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
1. Development environment should strive to be the same as production
2. Avoid clutter by consolidate into a single folder all the docker `build`/`run`/`tag`/`push` commands and kubernetes manifests   
3. Easily manage and expand your kubernetes supported resources 
4. Utilize simple commands that encapsulate your repetitive docker and kubernetes cli commands

## How does it work?
#### Directory Structure & Instructions
Please refer to [anchor-dockerfiles](https://github.com/ZachiNachshon/anchor-dockerfiles) repository for additional details.

## Requirements
- Go 1.12.x

## Download

#### I don't have GO environment 
Download your OS and ARCH relevant binary from [releases](https://github.com/ZachiNachshon/anchor/releases), unzip and place in `/usr/bin` / `ust/local/bin`.

#### I do have GO environment
Clone anchor repository and build as follows:
```bash
~$ git clone https://github.com/ZachiNachshon/anchor.git ~/anchor-example/anchor
~$ cd ~/anchor-example/anchor
~$ make build
```

## Quick Start Guide

Clone an example dockerfiles git repository
```bash
~$ git clone https://github.com/ZachiNachshon/anchor-dockerfiles.git ~/anchor-example/anchor-dockerfiles
```

Define required environment variable (append to `$PATH` via `~/.bash_profile` / `~/.bashrc`)
```bash
~$ cd ~/anchor-example/anchor-dockerfiles
~$ export DOCKER_FILES="$(PWD)"
``` 

List all available docker supported images/manifests
```bash
~$ anchor docker list

----------------------- Listing all Docker images ------------------------
  alpine
  nginx

    Done.
```

```bash
~$ anchor cluster list

----------------------- Listing Containers With K8S Manifests -----------------------
  alpine
  nginx

    Done.
```

Let's create a cluster (patience needed since the first creation pulls image `kindest/node:v1.15.0`)
```bash
~$ anchor cluster create
```

> Follow on-screen dashboard instructions to enter the Kubernetes Web-UI

Make sure all pods are running as expected
```bash
~$ anchor cluster status
```

Build an `nginx` docker image
```bash
~$ anchor docker build nginx
```

Push to private docker registry
```bash
~$ anchor docker push nginx
```

Verify `nginx` image exists on registry catalog
```bash
~$ anchor cluster registry
```

Deploy `nginx` kubernetes manifest and check status on K8s dashboard
```bash
~$ anchor cluster deploy nginx
```

Expose `nginx` port forwarding to the host instance
```bash
~$ anchor cluster expose nginx
```

Interact with `nginx` service
```bash
~$ curl -X GET http://localhost:1234
```

Connect to an `nginx` pod
```bash
~$ anchor cluster connect nginx
```

Remove `nginx` kubernetes manifest
```bash
~$ anchor cluster remove nginx
```

Delete kubernetes cluster
```bash
~$ anchor cluster delete
```
---
List of available `anchor docker` commands:
```bash
Usage:
  anchor docker [command]

Aliases:
  docker, d

Available Commands:
  build       Builds a docker image
  clean       Clean docker containers and images
  list        List all available docker supported images from DOCKER_FILES repository
  purge       Purge all docker images and containers
  push        Push a docker image to repository [registry.anchor:32001]
  run         Run a docker container
  stop        Stop a docker container
```

List of available `anchor cluster` commands:
```bash
Usage:
  anchor cluster [command]

Aliases:
  cluster, c

Available Commands:
  connect     Connect to a kubernetes pod by name
  create      Create a local Kubernetes cluster
  dashboard   Deploy a Kubernetes dashboard
  deploy      Deploy a container Kubernetes manifest
  destroy     Destroy local Kubernetes cluster
  expose      Expose a container port to the host instance
  list        List all containers with Kubernetes manifests from DOCKER_FILES repository
  registry    Create a private docker registry [registry.anchor:32001]
  remove      Removed a previously deployed container manifest
  status      Print cluster [registry.anchor:32001] status
```
---

