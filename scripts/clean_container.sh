#!/bin/bash
#
# This script allows you to stop running containers
# from previous docker run.
#
set -e
set -o pipefail

DOCKER_REPO_PREFIX=${DOCKER_REPO_PREFIX}
DOCKER_REGISTRY=${DOCKER_REGISTRY}
DOT_FILES=${DOT_FILES}

echo ----------------------- Cleanup: Previous Containers -------------------------

#-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
#       CLEANUP: CONTAINER
#-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

stop_container() {
  img_name=$1

  containersToStop=$(docker ps | grep "${DOCKER_REPO_PREFIX}/${img_name}" | awk {'print $1'})
  if [[ ! -z "${containersToStop}" ]]; then
      echo "Stopping all running container(s) of "${DOCKER_REPO_PREFIX}/${img_name}"..."
      docker stop ${containersToStop}
  fi
}

remove_running_containers() {
  img_name=$1

  runningContainerToRemove=$(docker ps | grep "${DOCKER_REPO_PREFIX}/${img_name}" | awk {'print $1'})
  if [[ ! -z "${runningContainerToRemove}" ]]; then
      echo "Removing running container(s) of "${DOCKER_REPO_PREFIX}/${img_name}"..."
      docker rm -f ${runningContainerToRemove}
  fi
}

remove_stopped_containers() {
  img_name=$1

  stoppedContainerToRemove=$(docker ps -a | grep "${DOCKER_REPO_PREFIX}/${img_name}" | awk {'print $1'})
  if [[ ! -z "${stoppedContainerToRemove}" ]]; then
      echo "Removing stopped container(s) of "${DOCKER_REPO_PREFIX}/${img_name}"..."
      docker rm -f ${stoppedContainerToRemove}
  fi
}

main() {

  if [[ $# -eq 0 ]]; then
    echo "Usage: $0 image1 image2 ..."
    exit 1
  fi

  for name in "$@"; do
    if [[ ! -d "$name" ]]; then
      echo "Unable to find container configuration directory with name: $name"
      exit 1
    else
      echo "Cleaning: ${name}"
      stop_container ${name} || true
      remove_running_containers ${name} || true
      remove_stopped_containers ${name} || true
    fi

    shift
  done
}

main "$@"

echo -e "\n    Done.
  "