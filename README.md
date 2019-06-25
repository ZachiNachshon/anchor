![anchor-side](assets/anchor-logo-side-700px.png)

Anchor is a utility intended for managing a local Docker / Kubernetes development environment.

## Overview
Anchor rely on `DOCKER_FILES` ENV variable to contain local path to a `dockerfiles` repository.<br/>
Every directory in the `dockerfiles` repository is identified as a container name and must contain a single `Dockerfile` inside.<br/>
It is optional to include a `.env` file in the `dockerfiles` root directory or in each of the sub directories, it'll get exported before image/container creation.

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

