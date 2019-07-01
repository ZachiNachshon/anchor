![anchor-side](assets/anchor-logo-side-600px.png)

## Overview
Anchor is a utility intended for managing an ephemeral local Docker / Kubernetes development environment with ease.
<br> It is based on:
- [Docker](https://github.com/docker) for running containerized applications
- [Kind](https://github.com/kubernetes-sigs/kind) for running local K8s cluster as a docker container
- [HomeBrew](https://github.com/Homebrew/brew) for managing maxOS/Linux packages
- [Kubectl](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubectl) for running commands against Kubernetes clusters 

## Why
1. Development environment should strive to be the same as production
2. Avoid clutter by consolidate into a single folder all the docker build/run/tag/push commands and kubernetes manifests   
3. Easily manage and expand your kubernetes supported resources 
4. Utilize simple commands that encapsulate your repetitive docker and kubernetes cli commands

## How does it work?
#### Directory Structure
Anchor act as one stop shop for all Dockerfiles & Kubernetes manifests that comprise your development environment.
It relies on a `DOCKER_FILES` ENV variable that points to a local directory path containing the following structure:

    .
    ├── ...
    ├── nginx                   # Name of the docker image/container
    │   └── k8s                 # Kubernetes content
    │       └── manifest.yaml   # Kubernetes manifest
    │   ├── Dockerfile          # Docker build/run/tag/push instructions
    │   ├── .env                # Optional: Override root `NAMESPACE` & `TAG` environment vars <br>
    │   └── ...                 # Optional: files for docker build <br>
    ├── ... 
    └── .env                    # Optional: Override default `NAMESPACE` & `TAG` environment vars at root level                 
    
> It is recommended to back the `DOCKER_FILES` directory by a git repository

Anchor allow flexibility for changing the namespace and tag of all or specific container image/deployment.<br>
The following environment variables can be set on one of the `.env` files:
- `export NAMESPACE="my-namespace`
- `export TAG="v1.1.0`

> Default values are `anchor` namespace and `latest` tag

#### Dockerfile Instructions
Every Dockerfile must contain the following heading in order to integrate properly with `anchor`
```dockerfile
# OVERVIEW
# --------
# This is the Dockerfile for nginx
#
# REQUIRED BASE IMAGE TO BUILD THIS IMAGE
# ---------------------------------------
# None.
#
# REQUIRED FILES TO BUILD THIS IMAGE
# ----------------------------------
# (1) None.
#
# HOW TO BUILD THIS IMAGE
# -----------------------
# docker build -f Dockerfile \
#              -t ${NAMESPACE}/nginx:${TAG} \
#              .
#
# HOW TO RUN THIS CONTAINER
# -------------------------
# docker run -t -d \
#            -v ${HOME}/.nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
#            --name=${NAMESPACE}-nginx \
#            -p 8080:80 \
#            ${NAMESPACE}/alpine:${TAG}
#
# HOW TO TAG THIS IMAGE
# ---------------------
# docker tag ${NAMESPACE}/nginx:${TAG} \
#            ${REGISTRY}/${NAMESPACE}/nginx:${TAG}
#
# HOW TO PUSH THIS IMAGE
# ----------------------
# docker push ${REGISTRY}/${NAMESPACE}/nginx:${TAG}
# 
```

#### Kubernetes Manifest
Standard kubernetes manifests.
Anchor support ENV vars substitution within manifests using `envsubst`.

## Requirements
- Go 1.12.x

> **Important !**<br>
> Temporary add the following line to `/etc/hosts`:<br>
> `127.0.0.1	localhost registry.anchor`<br>
>
> _Issue should be automated using [go-dnsmasq](https://github.com/janeczku/go-dnsmasq) later on._


## Quick Start Guide

Clone and build Anchor repository
```bash
~$ git clone git@github.com:ZachiNachshon/anchor.git ~/anchor-example/
~$ ~/anchor-example/make build 
```

Clone an example dockerfiles git repository
```bash
~$ git clone git@github.com:ZachiNachshon/anchor-dockerfiles.git ~/anchor-example/
```

Define required environment variable (append to `$PATH` via `.bash_profile` / `.bash_rc`)
```bash
~$ export DOCKER_FILES="~/anchor-example/anchor-dockerfiles"
```

Let's create a cluster
```bash
~$ anchor cluster create
```

> Follow on-screen dashboard instructions to gain Kubernetes visibility 

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

Build an `nginx` docker image
```bash
~$ anchor docker build nginx
```

Push to private docker registry
```bash
~$ anchor docker push nginx
```

Deploy nginx kubernetes manifest
```bash
~$ anchor cluster deploy nginx
```

Interact with `nginx` service
```bash
~$ curl -X GET http://localhost:1234
```

Remove `nginx` kubernetes manifest
```bash
~$ anchor cluster remove nginx
```

Delete kubernetes cluster
```bash
~$ anchor cluster delete
```

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

Available Commands:
  create      Create a local Kubernetes cluster
  dashboard   Deploy a Kubernetes dashboard
  delete      Delete local Kubernetes cluster
  deploy      Deploy a container Kubernetes manifest
  list        List all containers with Kubernetes manifests from DOCKER_FILES repository
  registry    Create a private docker registry [registry.anchor:32001]
  remove      Removed a previously deployed container manifest
  status      Print cluster [registry.anchor:32001] status
```
---

